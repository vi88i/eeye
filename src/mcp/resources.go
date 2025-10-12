package mcp

import (
	"context"
	"eeye/src/db"
	"eeye/src/models"
	"eeye/src/utils"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func handleResource(
	_ context.Context,
	req *mcp.ReadResourceRequest,
) (*mcp.ReadResourceResult, error) {
	u, err := url.Parse(req.Params.URI)
	if err != nil {
		return nil, fmt.Errorf("invalid URI: %w", err)
	}

	scheme := u.Scheme
	log.Printf("HandleResource scheme: %v\n", scheme)

	switch scheme {
	case "db":
		resource := u.Opaque
		log.Printf("HandleResource resource: %v\n", resource)

		switch resource {
		case "stocks":
			return handleStocksResource(req)
		default:
			return nil, fmt.Errorf("invalid resource: %v", resource)
		}
	}

	return nil, fmt.Errorf("invalid scheme: %v", scheme)
}

func handleStocksResource(req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	stocks, err := db.FetchAllStocks()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stocks: %w", err)
	}

	text := strings.Join(
		utils.Map(
			stocks,
			func(e models.Stock) string {
				return e.Symbol
			},
		),
		",",
	)
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "text/plain",
				Text:     text,
			},
		},
	}, nil
}

func addResources(server *mcp.Server) {
	server.AddResource(
		&mcp.Resource{
			MIMEType:    "text/plain",
			Name:        "nseStocks",
			Title:       "NSE stocks",
			Description: "List of NSE stocks",
			URI:         "db:stocks",
		},
		handleResource,
	)
}

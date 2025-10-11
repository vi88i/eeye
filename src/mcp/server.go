// Package mcp creates MCP server
package mcp

import (
	"context"
	"eeye/src/config"
	"fmt"
	"log"
	"net/http"
	"net/url"

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

// Init is a facade for MCP server functionality
func Init() {
	serverImpl := &mcp.Implementation{
		Name:    "eeye-mcp",
		Version: "v0.0.1",
		Title:   "eeye stock screener",
	}

	serverOpts := &mcp.ServerOptions{
		Instructions: "Use this server for NSE stock analysis queries!",
		HasResources: true,
	}

	server := mcp.NewServer(serverImpl, serverOpts)

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

	reqHandler := mcp.NewStreamableHTTPHandler(
		func(_ *http.Request) *mcp.Server {
			return server
		},
		nil,
	)

	url := fmt.Sprintf("%v:%v", config.MCP.Host, config.MCP.Port)
	log.Printf("Starting MCP HTTP streamable transport server on %v\n", url)
	if err := http.ListenAndServe(url, reqHandler); err != nil {
		log.Fatal(err)
	}
}

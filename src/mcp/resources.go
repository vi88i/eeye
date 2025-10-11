package mcp

import (
	"eeye/src/db"
	"eeye/src/models"
	"eeye/src/utils"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

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

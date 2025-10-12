// Package mcp creates MCP server
package mcp

import (
	"eeye/src/config"
	"fmt"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

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
		HasTools:     true,
	}

	server := mcp.NewServer(serverImpl, serverOpts)
	addResources(server)
	addTools(server)

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

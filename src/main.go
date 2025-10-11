// Package main is the entry point for the eeye trading system.
// It initializes the trading environment, loads configurations, and starts
// the trading strategies execution pipeline.
package main

import (
	"eeye/src/api"
	"eeye/src/config"
	"eeye/src/db"
	"eeye/src/handlers"
	"eeye/src/mcp"
	"eeye/src/strategy"
	"flag"
	"log"
)

func main() {
	mcpMode := flag.Bool("mcp", false, "Enable to start MCP server")
	flag.Parse()

	applog := handlers.GetAppLog()
	config.Load()
	api.InitGrowwTradingClient()
	api.InitNSEClient()
	db.Connect()

	if *mcpMode {
		mcp.Init()
	} else {
		quit := handlers.GetInterruptHandlerChannel()
		done := strategy.Analyze()

		select {
		case sig := <-quit:
			log.Println("Shutting down gracefully, signal caught:", sig.String())
		case <-done:
			// Do the deletion only after analysis
			db.DeleteDelistedStocks()
		}
	}

	db.Disconnect()
	applog.Close()
}

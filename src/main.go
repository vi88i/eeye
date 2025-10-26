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
	cleanUp := flag.Bool("cleanup", false, "Clean up de-listed stocks")
	verbose := flag.Bool("verbose", false, "Print logs in stdout/stderr")
	flag.Parse()

	applog := handlers.GetAppLog(*verbose)
	config.Load()
	api.InitGrowwTradingClient()
	api.InitNseClient()
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
			if *cleanUp {
				db.DeleteDelistedStocks()
			}
		}
	}

	db.Disconnect()
	applog.Close()
}

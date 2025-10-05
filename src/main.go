// Package main is the entry point for the eeye trading system.
// It initializes the trading environment, loads configurations, and starts
// the trading strategies execution pipeline.
package main

import (
	"eeye/src/api"
	"eeye/src/config"
	"eeye/src/db"
	"eeye/src/handlers"
	"eeye/src/strategy"
)

func main() {
	applog := handlers.GetAppLog()
	config.Load()
	api.InitGrowwTradingClient()
	api.InitNSEClient()
	db.Connect()

	done := strategy.Analyze()
	quit := handlers.GetInterruptHandlerChannel()

	select {
	case <-quit:
	case <-done:
		db.DeleteDelistedStocks()
	}

	db.Disconnect()
	applog.Close()
}

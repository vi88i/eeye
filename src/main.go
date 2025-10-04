// Package main is the entry point for the eeye trading system.
// It initializes the trading environment, loads configurations, and starts
// the trading strategies execution pipeline.
package main

import (
	"eeye/src/api"
	"eeye/src/applog"
	"eeye/src/config"
	"eeye/src/db"
	"eeye/src/strategy"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logFile := applog.Init()
	config.Load()
	api.InitGrowwTradingClient()
	api.InitNSEClient()
	db.Connect()

	done := strategy.Analyze()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-quit:
	case <-done:
		db.DeleteDelistedStocks()
	}

	db.Disconnect()
	var _ = logFile.Close()
}

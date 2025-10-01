// Package main is the entry point for the eeye trading system.
// It initializes the trading environment, loads configurations, and starts
// the trading strategies execution pipeline.
package main

import (
	"eeye/src/api"
	"eeye/src/config"
	"eeye/src/db"
	"eeye/src/steps"
	"eeye/src/strategy"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.Load()
	api.InitGrowwTradingClient()
	api.InitNSEClient()
	db.Connect()

	go func() {
		stocks := steps.GetStocks()
		strategy.Executor(stocks)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	db.Disconnect()
}

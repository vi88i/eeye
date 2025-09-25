// Package main is the entry point for the eeye trading system.
// It initializes the trading environment, loads configurations, and starts
// the trading strategies execution pipeline.
package main

import (
	"eeye/src/api"
	"eeye/src/config"
	"eeye/src/constants"
	"eeye/src/db"
	"eeye/src/strategy"
	"eeye/src/utils"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	var (
		stocksYamlPathPtr = flag.String(
			"stocks",
			constants.StocksYamlPath,
			"yaml file with list of stocks. Eg: examples/stocks.yaml",
		)
	)

	flag.Parse()

	stocks := utils.GetStocksFromYaml(filepath.Clean(*stocksYamlPathPtr))
	config.Load()
	api.InitTradingClient()
	db.Connect()

	go func() {
		strategy.Executor(stocks)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	db.Disconnect()
}

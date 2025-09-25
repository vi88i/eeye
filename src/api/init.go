// Package api provides functionality for interacting with trading APIs.
// It handles API client initialization, data fetching, and communication
// with external trading services.
package api

import (
	"eeye/src/config"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// Client is the shared HTTP client for making API requests.
// It is configured with the necessary headers and base URL for the trading API.
var Client *resty.Client

// InitTradingClient initializes the global HTTP client with proper configuration
// for making requests to the trading API. This includes setting up authentication,
// API version headers, and base URL.
func InitTradingClient() {
	Client = resty.New()
	Client.SetHeader("Authorization", "Bearer "+config.TradingAPIConfig.AccessToken)
	Client.SetHeader("X-API-VERSION", config.TradingAPIConfig.XAPIVersion)
	Client.SetHeader("Accept", "application/json")
	Client.SetBaseURL(fmt.Sprintf("%v/%v", config.TradingAPIConfig.BaseURL, config.TradingAPIConfig.APIVersion))
}

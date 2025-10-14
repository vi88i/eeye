// Package api provides functionality for interacting with trading APIs.
// It handles API client initialization, data fetching, and communication
// with external trading services.
package api

import (
	"eeye/src/config"
	"eeye/src/constants"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// GrowwClient is the shared HTTP client for making API requests.
// It is configured with the necessary headers and base URL for the trading API.
var GrowwClient *resty.Client

// NseClient is the shared HTTP client for making requests to the NSE API.
var NseClient *resty.Client

// InitNSEClient initializes the global HTTP client with proper configuration
// for making requests to the NSE API. This includes setting up headers
// and base URL.
func InitNSEClient() {
	NseClient = resty.New()
	NseClient.SetHeader("User-Agent", constants.ReqNSEUserAgent)
	NseClient.SetHeader("Accept", "application/json")
	NseClient.SetHeader("Cache-Control", "no-cache, no-store, must-revalidate")
	NseClient.SetHeader("Pragma", "no-cache")
	NseClient.SetHeader("Expires", "0")
	NseClient.SetBaseURL(config.NSE.BaseURL)
}

// InitGrowwTradingClient initializes the global HTTP client with proper configuration
// for making requests to the trading API. This includes setting up authentication,
// API version headers, and base URL.
func InitGrowwTradingClient() {
	GrowwClient = resty.New()
	GrowwClient.SetHeader("Authorization", "Bearer "+config.Groww.AccessToken)
	GrowwClient.SetHeader("X-API-VERSION", config.Groww.XAPIVersion)
	GrowwClient.SetHeader("Accept", "application/json")
	GrowwClient.SetBaseURL(fmt.Sprintf("%v/%v", config.Groww.BaseURL, config.Groww.APIVersion))
}

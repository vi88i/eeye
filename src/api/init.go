package api

import (
	"eeye/src/config"
	"fmt"

	"github.com/go-resty/resty/v2"
)

var Client *resty.Client

func InitTradingClient() {
	Client = resty.New()
	Client.SetHeader("Authorization", "Bearer "+ config.TradingAPIConfig.AccessToken)
	Client.SetHeader("X-API-VERSION", config.TradingAPIConfig.XAPIVersion)
	Client.SetHeader("Accept", "application/json")
	Client.SetBaseURL(fmt.Sprintf("%v/%v", config.TradingAPIConfig.BaseURL, config.TradingAPIConfig.APIVersion))
}

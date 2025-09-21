package config

import (
	"eeye/src/constants"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var TradingAPIConfig = struct {
	AccessToken string
	BaseURL     string
	APIVersion  string
	XAPIVersion string
	RateLimit   int
}{}

var DBConfig = struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Tz       string
}{}

var StepsConfig = struct {
	Concurrency int
}{constants.MinConcurrency}

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file")
	}

	TradingAPIConfig.AccessToken = os.Getenv("GROWW_ACCESS_TOKEN")
	TradingAPIConfig.BaseURL = os.Getenv("GROWW_BASE_URL")
	TradingAPIConfig.APIVersion = os.Getenv("GROWW_API_VERSION")
	TradingAPIConfig.XAPIVersion = os.Getenv("GROWW_X_API_VERSION")

	DBConfig.Host = os.Getenv("EEYE_DB_HOST")
	DBConfig.Port = os.Getenv("EEYE_DB_PORT")
	DBConfig.User = os.Getenv("EEYE_DB_USER")
	DBConfig.Password = os.Getenv("EEYE_DB_PASSWORD")
	DBConfig.Name = os.Getenv("EEYE_DB_NAME")
	DBConfig.Tz = os.Getenv("EEYE_TZ")

	concurrency, err := strconv.Atoi(os.Getenv("EEYE_CONCURRENCY"))
	if err == nil {
		StepsConfig.Concurrency = concurrency
	} else {
		log.Println("invalid EEYE_CONCURRENCY")
	}

	rateLimit, err := strconv.Atoi(os.Getenv("GROWW_RATE_LIMIT_PER_SECOND"))
	if err == nil {
		TradingAPIConfig.RateLimit = 1
	} else {
		TradingAPIConfig.RateLimit = rateLimit
	}
}

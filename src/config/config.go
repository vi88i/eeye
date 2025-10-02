// Package config manages the application configuration.
// It provides structures and functions for loading and accessing
// configuration parameters for database, API, and trading steps.
package config

import (
	"eeye/src/constants"
	"eeye/src/utils"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// TradingAPIConfig holds configuration for the trading API connection.
// It includes authentication, endpoint, and rate limiting settings.
var TradingAPIConfig = struct {
	// AccessToken is the API authentication token
	AccessToken string

	// BaseURL is the root URL for API requests
	BaseURL string

	// APIVersion specifies the API version to use
	APIVersion string

	// XAPIVersion is the custom API version header value
	XAPIVersion string

	// RequestPerSecond defines the maximum API requests per second
	RequestPerSecond int
}{RequestPerSecond: constants.MinRequestPerSecond}

// DBConfig holds the PostgreSQL database connection configuration.
var DBConfig = struct {
	// Host is the database server hostname
	Host string

	// Port is the database server port
	Port string

	// User is the database username
	User string

	// Password is the database user password
	Password string

	// Name is the database name to connect to
	Name string

	// Tz is the timezone for database connections
	Tz string
}{}

// NSEConfig holds configuration for NSE Bhavcopy downloads
var NSEConfig = struct {
	// BaseURL is the base URL for NSE Bhavcopy downloads
	BaseURL string
}{}

// Load reads configuration from environment variables and initializes
// the application's configuration structures. It will panic if required
// environment variables are missing or invalid.
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

	requestPerSecond, err := strconv.Atoi(os.Getenv("GROWW_RPS"))
	if err == nil {
		TradingAPIConfig.RequestPerSecond = utils.Clamp[int, int](
			requestPerSecond,
			constants.MinRequestPerSecond,
			constants.MaxRequestPerSecond,
		)
	} else {
		log.Println("invalid GROWW_RPS defaulting to", constants.MaxRequestPerSecond)
	}

	NSEConfig.BaseURL = os.Getenv("NSE_BHAVCOPY_BASE_URL")
}

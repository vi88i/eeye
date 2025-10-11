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

// Groww holds configuration for the trading API connection.
// It includes authentication, endpoint, and rate limiting settings.
var Groww = struct {
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

// DB holds the PostgreSQL database connection configuration.
var DB = struct {
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

// NSE holds configuration for NSE Bhavcopy downloads
var NSE = struct {
	// BaseURL is the base URL for NSE Bhavcopy downloads
	BaseURL string
}{}

// MCP holds the MCP server configuration
var MCP = struct {
	// Host is the MCP server host
	Host string

	// Port is the MCP server port
	Port string
}{}

// Load reads configuration from environment variables and initializes
// the application's configuration structures. It will panic if required
// environment variables are missing or invalid.
func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file")
	}

	Groww.AccessToken = os.Getenv("GROWW_ACCESS_TOKEN")
	Groww.BaseURL = os.Getenv("GROWW_BASE_URL")
	Groww.APIVersion = os.Getenv("GROWW_API_VERSION")
	Groww.XAPIVersion = os.Getenv("GROWW_X_API_VERSION")

	DB.Host = os.Getenv("EEYE_DB_HOST")
	DB.Port = os.Getenv("EEYE_DB_PORT")
	DB.User = os.Getenv("EEYE_DB_USER")
	DB.Password = os.Getenv("EEYE_DB_PASSWORD")
	DB.Name = os.Getenv("EEYE_DB_NAME")
	DB.Tz = os.Getenv("EEYE_TZ")

	requestPerSecond, err := strconv.Atoi(os.Getenv("GROWW_RPS"))
	if err == nil {
		Groww.RequestPerSecond = utils.Clamp[int, int](
			requestPerSecond,
			constants.MinRequestPerSecond,
			constants.MaxRequestPerSecond,
		)
	} else {
		log.Println("invalid GROWW_RPS defaulting to", constants.MaxRequestPerSecond)
	}

	NSE.BaseURL = os.Getenv("NSE_BHAVCOPY_BASE_URL")

	MCP.Host = os.Getenv("MCP_HOST")
	MCP.Port = os.Getenv("MCP_PORT")
}

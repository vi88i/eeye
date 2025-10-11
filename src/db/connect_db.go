// Package db handles database operations for the trading system.
// It provides functionality for connecting to PostgreSQL, managing connections,
// and executing queries for storing and retrieving trading data.
package db

import (
	"context"
	"eeye/src/config"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Pool is the global connection pool for PostgreSQL database access.
// It provides managed, concurrent access to database connections.
var Pool *pgxpool.Pool

// Connect initializes the global database connection pool using configuration
// from DB. It will panic if the connection cannot be established.
func Connect() {
	var databaseURL = fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v",
		config.DB.User,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
	)

	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("Unable to parse config: %v\n", err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	Pool = pool
	log.Println("Connected to database")
}

// Disconnect closes the database connection pool if it exists.
// This should be called when shutting down the application to ensure
// all database connections are properly closed.
func Disconnect() {
	if Pool != nil {
		Pool.Close()
		log.Println("Disconnected from database")
	}
}

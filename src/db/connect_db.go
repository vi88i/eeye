package db

import (
	"context"
	"eeye/src/config"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() {
	var databaseUrl = fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v",
		config.DBConfig.User,
		config.DBConfig.Password,
		config.DBConfig.Host,
		config.DBConfig.Port,
		config.DBConfig.Name,
	)

	cfg, err := pgxpool.ParseConfig(databaseUrl)
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

func Disconnect() {
	if Pool != nil {
		Pool.Close()
		log.Println("Disconnected from database")
	}
}

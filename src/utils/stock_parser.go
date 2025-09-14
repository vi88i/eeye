package utils

import (
	"eeye/src/models"
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

func GetStocksFromYaml(path string) []models.Stock {
	data, err := os.ReadFile(path)

	if err != nil {
		log.Fatalf("failed to read stocks yaml: %v", err)
	}

	var cfg = models.StocksConfig{Stocks: make([]models.Stock, 0, 100)}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("failed to parse yaml: %v", err)
	}

	return cfg.Stocks
}

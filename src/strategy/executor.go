package strategy

import (
	"eeye/src/models"
	"log"
)

func Executor(stocks []models.Stock) {
	results := []string{
		lowerBollingerBandBullish(stocks),
		bullishSwing(stocks),
		emaFakeBreakdown(stocks, 50),
	}

	log.Println("================= Strategy Results =================")
	for _, result := range results {
		log.Println(result)
	}
}

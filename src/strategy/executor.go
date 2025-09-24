package strategy

import (
	"eeye/src/models"
	"fmt"
)

func Executor(stocks []models.Stock) {
	results := []string{
		lowerBollingerBandBullish(stocks),
		bullishSwing(stocks),
	}

	fmt.Println("================= Strategy Results =================")
	for _, result := range results {
		fmt.Println(result)
	}
}

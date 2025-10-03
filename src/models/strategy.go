package models

import "eeye/src/constants"

// Strategy defines the interface that all trading strategies must implement.
type Strategy interface {
	// Execute runs the strategy logic on the given stock and sends results to the output channel.
	Execute(stock *Stock)

	// Name returns the name of the strategy.
	Name() string

	// GetSink returns the output channel for the strategy.
	GetSink() chan *Stock

	// mustEmbedStrategyBaseImpl is a marker function to ensure that
	// StrategyBaseImpl is embedded in all strategies which helps in code re-using.
	mustEmbedStrategyBaseImpl()
}

// StrategyBaseImpl provides a base implementation for the Strategy interface.
type StrategyBaseImpl struct {
	// sink is the output channel for the strategy.
	sink chan *Stock
}

// mustEmbedStrategyBaseImpl is intentionally left blank to enforce embedding.
//
//nolint:unused
func (s *StrategyBaseImpl) mustEmbedStrategyBaseImpl() {}

// GetSink returns the output channel for the strategy, initializing it if necessary.
func (s *StrategyBaseImpl) GetSink() chan *Stock {
	if s.sink == nil {
		s.sink = make(chan *Stock, constants.AggregatorBufferSize)
	}

	return s.sink
}

// StrategyResult combines the strategy and the the result satisfying the strategy
type StrategyResult struct {
	// Strategy config
	Strategy Strategy

	// Stocks is list of stocks satisfying the strategy
	Stocks []*Stock
}

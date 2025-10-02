// Package constants defines system-wide constant values.
// It includes configuration defaults, API endpoints, time formats,
// and other immutable values used throughout the application.
package constants

const (
	// NumOfStrategyWorkers defines the maximum number of concurrent strategy workers allowed
	NumOfStrategyWorkers = 12

	// StrategyWorkerInputBufferSize defines the buffer size for each strategy worker's input channel
	StrategyWorkerInputBufferSize = 100

	// StrategyWorkerOutputBufferSize defines the buffer size for each strategy worker's output channel
	StrategyWorkerOutputBufferSize = 500

	// NumOfIngestionWorkers defines the number of concurrent ingestion workers
	NumOfIngestionWorkers = 4

	// IngestionBufferSize defines the buffer size for the ingestion worker's input channel
	IngestionBufferSize = 20

	// AggregatorBufferSize defines the buffer size for the aggregator's input channel
	AggregatorBufferSize = 20
)

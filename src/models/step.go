package models

import "log"

// Step defines the interface that each screen step should implement.
// All step implementations must embed StepBaseImpl to satisfy this interface.
type Step interface {
	// Name returns the name of step
	Name() string

	// Screen returns whether the stock passed the screening test for the strategy
	Screen(strategy string, stock *Stock) bool

	// mustEmbedStepBaseImpl is a marker method to force consumers to compose StepBaseImpl.
	// This ensures all Step implementations have access to shared helper methods.
	mustEmbedStepBaseImpl()
}

// StepBaseImpl provides base implementation and helper methods for all Step implementations.
// All step types must embed this struct to satisfy the Step interface.
type StepBaseImpl struct{}

// mustEmbedStepBaseImpl is a marker method that forces all Step implementations
// to embed StepBaseImpl, ensuring consistent behavior across all steps.
//
//nolint:unused
func (s *StepBaseImpl) mustEmbedStepBaseImpl() {}

// TruthyCheck is a helper method that executes the provided assertion function,
// logs a failure message if the test fails, and returns the test result.
// This reduces code duplication across step implementations by centralizing
// the common pattern of test execution and conditional logging.
func (s *StepBaseImpl) TruthyCheck(
	strategy string,
	step string,
	stock *Stock,
	assert func() bool,
) bool {
	test := assert()
	if !test {
		log.Printf("[%v - %v] test failed: %v\n", strategy, step, stock.Symbol)
	}
	return test
}

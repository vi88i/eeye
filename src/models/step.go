package models

// Step defines the interface that each screen step should implement
type Step interface {
	// Name returns the name of step
	Name() string

	// Screen returns whether the stock passed the screening test for the strategy
	Screen(strategy string, stock *Stock) bool
}

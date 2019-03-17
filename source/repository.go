package source

import "context"

// Repository it's the interface for dealing with the
// Source entity
type Repository interface {
	// Filter returns a list of Source without
	// any filter
	Filter(ctx context.Context) ([]Source, error)
}

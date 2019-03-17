package notification

import "context"

// Repository it's the interface for dealing with the
// Notification entity
type Repository interface {
	// Filter returns a list of Notification without
	// any filter
	Filter(ctx context.Context) ([]*Notification, error)
}

package notification

import "context"

type Repository interface {
	Filter(ctx context.Context) ([]*Notification, error)
}

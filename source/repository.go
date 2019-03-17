package source

import "context"

type Repository interface {
	Filter(ctx context.Context) ([]Source, error)
}

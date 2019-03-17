package immem

import (
	"context"

	"github.com/xescugc/notigator/source"
)

type sourceRepository struct {
	sources []source.Source
}

// NewSourceRepository returns the implementation of a source.Repository
// from an "in memory" sotrage
func NewSourceRepository(srcs []source.Source) source.Repository {
	return &sourceRepository{
		sources: srcs,
	}
}

func (r *sourceRepository) Filter(ctx context.Context) ([]source.Source, error) {
	return r.sources, nil
}

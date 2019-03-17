package service

import (
	"context"
	"fmt"

	"github.com/xescugc/notigator/notification"
	"github.com/xescugc/notigator/source"
)

// Service it's the public interface for the Notigator Service
type Service interface {
	// GetSources returns a list of all the configured source.Source
	GetSources(ctx context.Context) ([]source.Source, error)

	// GetSourceNotifications returns all the notifications of the srcID
	GetSourceNotifications(ctx context.Context, srcID string) ([]*notification.Notification, error)
}

type service struct {
	sources       source.Repository
	notifications map[string]notification.Repository
}

// New returns an implementation of the Service interface
func New(srcs source.Repository, nots map[string]notification.Repository) Service {
	return &service{
		sources:       srcs,
		notifications: nots,
	}
}

func (s *service) GetSources(ctx context.Context) ([]source.Source, error) {
	return s.sources.Filter(ctx)
}

func (s *service) GetSourceNotifications(ctx context.Context, srcID string) ([]*notification.Notification, error) {
	nr, ok := s.notifications[srcID]
	if !ok {
		return nil, fmt.Errorf("source not defined %q", srcID)
	}

	nts, err := nr.Filter(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while Filtering notifications: %s", err)
	}

	return nts, nil
}

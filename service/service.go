package service

import (
	"context"
	"fmt"

	"github.com/xescugc/notigator/notification"
	"github.com/xescugc/notigator/source"
)

type Service interface {
	GetSources(ctx context.Context) ([]source.Source, error)
	GetSourceNotifications(ctx context.Context, srcCan source.Canonical) ([]*notification.Notification, error)
}

type service struct {
	gh notification.Repository
	gl notification.Repository
	tr notification.Repository
}

func New(gh, gl, tr notification.Repository) Service {
	return &service{
		gh: gh,
		gl: gl,
		tr: tr,
	}
}

func (s *service) GetSources(ctx context.Context) ([]source.Source, error) {
	return source.Sources, nil
}

func (s *service) GetSourceNotifications(ctx context.Context, srcCan source.Canonical) ([]*notification.Notification, error) {
	if !srcCan.IsACanonical() {
		return nil, fmt.Errorf("could not found source %q", srcCan)
	}

	nr, err := s.selectSource(srcCan)
	if err != nil {
		return nil, fmt.Errorf("could not select source: %s", err)
	}

	nts, err := nr.Filter(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while Filtering notifications: %s", err)
	}

	return nts, nil
}

func (s *service) selectSource(srcCan source.Canonical) (notification.Repository, error) {
	switch srcCan {
	case source.Github:
		return s.gh, nil
	case source.Gitlab:
		return s.gl, nil
	case source.Trello:
		return s.tr, nil
	default:
		return nil, fmt.Errorf("could not found source %q", srcCan)
	}
}

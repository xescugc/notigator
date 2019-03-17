package service

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/xescugc/notigator/service/models"
)

type response struct {
	Data interface{} `json:"data,omitempty"`
}

func makeGetSources(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		sources, err := s.GetSources(ctx)
		if err != nil {
			return nil, err
		}
		return response{Data: models.NewSources(sources)}, nil
	}
}

type getSourceNotificationsRequest struct {
	SourceID string
}

func makeGetSourceNotifications(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getSourceNotificationsRequest)
		notifications, err := s.GetSourceNotifications(ctx, req.SourceID)
		if err != nil {
			return nil, err
		}
		return response{Data: models.NewNotifications(notifications)}, nil
	}
}

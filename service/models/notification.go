package models

import (
	"time"

	"github.com/xescugc/notigator/notification"
)

type Notification struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Scope     string    `json:"scope"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewNotifications(ns []*notification.Notification) []*Notification {
	rns := make([]*Notification, 0, len(ns))

	for _, v := range ns {
		rns = append(rns, &Notification{
			ID:        v.ID,
			Title:     v.Title,
			URL:       v.URL.String(),
			Scope:     v.Scope,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return rns
}

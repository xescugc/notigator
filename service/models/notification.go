package models

import (
	"time"

	"github.com/xescugc/notigator/notification"
)

// Notification it's the entity used to
// transform to JSON the notification.Notification
type Notification struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Scope     string    `json:"scope"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewNotifications returns a conversion of the ns to Notification
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

package notification

import (
	"net/url"
	"time"
)

// Notification is the basic entity
// that represents a notification of
// any Source
type Notification struct {
	ID        string
	Title     string
	URL       url.URL
	UpdatedAt time.Time
	Scope     string
}

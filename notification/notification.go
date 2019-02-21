package notification

import (
	"net/url"
	"time"
)

type Notification struct {
	ID        string
	Title     string
	URL       url.URL
	UpdatedAt time.Time
	Scope     string
}

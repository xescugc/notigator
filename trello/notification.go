package trello

import (
	"context"
	"fmt"
	"net/url"

	"github.com/adlio/trello"
	"github.com/xescugc/notigator/notification"
)

type notificationRepository struct {
	client *trello.Client
}

// NewNotificationRepository returns the implementation of
// a notification.Repository for a trello Source
func NewNotificationRepository(apiKey, token string) notification.Repository {
	return &notificationRepository{
		client: trello.NewClient(apiKey, token),
	}
}

func (n *notificationRepository) Filter(ctx context.Context) ([]*notification.Notification, error) {
	nots, err := n.client.GetMyNotifications(trello.Arguments(map[string]string{
		"limit":       "1000",
		"read_filter": "unread",
	}))
	if err != nil {
		return nil, err
	}

	notifications := make([]*notification.Notification, 0)
	titles := make(map[string]struct{})
	for _, n := range nots {
		// If it's not a Card action we do not want it
		if n.Data.Card == nil {
			continue
		}
		// We do not want to repeat the notification for the same card again
		// The order of the notifications ensure that the new ones are first
		// so we do not have to replace them for old ones
		if _, ok := titles[n.Data.Card.Name]; ok {
			continue
		}

		surl := fmt.Sprintf("https://trello.com/c/%s", n.Data.Card.ShortLink)
		u, err := url.Parse(surl)
		if err != nil {
			return nil, fmt.Errorf("could not parse url %q: %s", surl, err)
		}
		notifications = append(notifications, &notification.Notification{
			ID:        n.ID,
			Title:     n.Data.Card.Name,
			URL:       *u,
			Scope:     n.Data.Board.Name,
			UpdatedAt: n.Date,
		})
		titles[n.Data.Card.Name] = struct{}{}
	}

	return notifications, nil
}

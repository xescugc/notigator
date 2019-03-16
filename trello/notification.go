package trello

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/adlio/trello"
	"github.com/xescugc/notigator/notification"
)

type notificationRepository struct {
	client *trello.Client
}

func NewNotificationRepository(ctx context.Context, apiKey, token string) notification.Repository {
	return &notificationRepository{
		client: trello.NewClient(apiKey, token),
	}
}

type Notification struct {
	ID              string           `json:"id"`
	IDAction        string           `json:"idAction"`
	Unread          bool             `json:"unread"`
	Type            string           `json:"type"`
	IDMemberCreator string           `json:"idMemberCreator"`
	Date            time.Time        `json:"date"`
	DateRead        time.Time        `json:"dataRead"`
	Data            NotificationData `json:"data,omitempty"`
	MemberCreator   *trello.Member   `json:"memberCreator,omitempty"`
}

type NotificationData struct {
	Text  string                 `json:"text"`
	Card  *NotificationDataCard  `json:"card,omitempty"`
	Board *NotificationDataBoard `json:"board,omitempty"`
}

type NotificationDataBoard struct {
	ID        string `json:"id"`
	ShortLink string `json:"shortLink"`
	Name      string `json:"name"`
}

type NotificationDataCard struct {
	ID        string `json:"id"`
	IDShort   int    `json:"idShort"`
	Name      string `json:"name"`
	ShortLink string `json:"shortLink"`
}

func (n *notificationRepository) Filter(ctx context.Context) ([]*notification.Notification, error) {
	nots := make([]Notification, 0)
	err := n.client.Get("members/me/notifications", trello.Arguments(map[string]string{
		"limit":       "1000",
		"read_filter": "unread",
	}), &nots)
	if err != nil {
		return nil, err
	}

	notifications := make([]*notification.Notification, 0)
	titles := make(map[string]struct{}, 0)
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

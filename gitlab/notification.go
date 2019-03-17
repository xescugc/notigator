package gitlab

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	gitlab "github.com/xanzy/go-gitlab"
	"github.com/xescugc/notigator/notification"
)

type notificationRepository struct {
	client *gitlab.Client
}

// NewNotificationRepository returns the implementation of
// a notification.Repository for a gitlab Source
func NewNotificationRepository(token string) notification.Repository {
	return &notificationRepository{
		client: gitlab.NewClient(nil, token),
	}
}

func (n *notificationRepository) Filter(ctx context.Context) ([]*notification.Notification, error) {
	todos, _, err := n.client.Todos.ListTodos(&gitlab.ListTodosOptions{
		State: gitlab.String("pending"),
	}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("could not fetch the list of notifications: %s", err)
	}

	notifications := make([]*notification.Notification, 0, len(todos))
	for _, t := range todos {
		u, err := url.Parse(t.TargetURL)
		if err != nil {
			return nil, fmt.Errorf("could not parse url %q: %s", t.TargetURL, err)
		}
		notifications = append(notifications, &notification.Notification{
			ID:        strconv.Itoa(t.ID),
			Title:     t.Target.Title,
			URL:       *u,
			Scope:     t.Project.PathWithNamespace,
			UpdatedAt: *t.CreatedAt,
		})
	}

	return notifications, nil
}

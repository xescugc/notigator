package zeplin

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/xescugc/notigator/notification"
)

type notificationRepository struct {
	token  string
	client *http.Client
}

func NewNotificationRepository(ctx context.Context, token string) notification.Repository {
	return &notificationRepository{
		token:  token,
		client: &http.Client{},
	}
}

type Notification struct {
	Params struct {
		Source struct {
			ID       string      `json:"_id"`
			Email    string      `json:"email"`
			Username string      `json:"username"`
			Emotar   interface{} `json:"emotar"`
			Avatar   interface{} `json:"avatar"`
		} `json:"source"`
		Project struct {
			ID   string `json:"_id"`
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"project"`
	} `json:"params"`
	Updated time.Time `json:"updated"`
	ID      string    `json:"_id"`
	Events  []struct {
		Screen struct {
			ID   string `json:"_id"`
			Name string `json:"name"`
		} `json:"screen"`
	} `json:"events"`
	ActionName string `json:"actionName"`
}

func (n *Notification) getTitle() string {
	return fmt.Sprintf("%s %s", n.Params.Source.Username, n.ActionName)
}

func (n *notificationRepository) Filter(ctx context.Context) ([]*notification.Notification, error) {
	nots, err := n.getNotifications(ctx)
	if err != nil {
		return nil, err
	}

	notifications := make([]*notification.Notification, 0, len(nots))
	for _, n := range nots {

		// For now ignore those
		if n.Params.Project.Name == "" {
			continue
		}

		surl := fmt.Sprintf("https://app.zeplin.io/project/%s", n.Params.Project.ID)
		u, err := url.Parse(surl)
		if err != nil {
			return nil, fmt.Errorf("could not parse url %q: %s", surl, err)
		}
		notifications = append(notifications, &notification.Notification{
			ID:        n.ID,
			Title:     n.getTitle(),
			URL:       *u,
			Scope:     n.Params.Project.Name,
			UpdatedAt: n.Updated,
		})
	}

	return notifications, nil
}

func (n *notificationRepository) getNotifications(ctx context.Context) ([]*Notification, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.zeplin.io/notifications?count=50", nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %s", err)
	}

	req.Header.Set("zeplin-token", n.token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not make request: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %s", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("http status %d with contentt %s", resp.StatusCode, body)
	}

	var nots struct {
		Notifications []*Notification `json:"notifications"`
	}

	err = json.Unmarshal(body, &nots)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal body: %s", err)
	}

	return nots.Notifications, nil
}

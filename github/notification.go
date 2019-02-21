package github

import (
	"context"
	"fmt"
	"strings"

	"net/url"

	"github.com/xescugc/go-github/github"
	"github.com/xescugc/notigator/notification"
	"golang.org/x/oauth2"
)

type notificationRepository struct {
	client *github.Client
}

func NewNotificationRepository(ctx context.Context, token string) notification.Repository {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &notificationRepository{
		client: client,
	}
}

func (n *notificationRepository) Filter(ctx context.Context) ([]*notification.Notification, error) {
	nots, _, err := n.client.Activity.ListNotifications(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the list of notifications: %s", err)
	}

	notifications := make([]*notification.Notification, 0, len(nots))
	for _, not := range nots {
		u, err := buildURL(not)
		if err != nil {
			return nil, fmt.Errorf("could not parse url: %s", err)
		}
		notifications = append(notifications, &notification.Notification{
			ID:        *not.ID,
			Title:     *not.Subject.Title,
			URL:       *u,
			Scope:     *not.Repository.FullName,
			UpdatedAt: *not.UpdatedAt,
		})
	}

	return notifications, nil
}

func buildURL(n *github.Notification) (*url.URL, error) {
	base := *n.Repository.HTMLURL
	issueURLSplit := strings.Split(*n.Subject.URL, "/")
	issueID := issueURLSplit[len(issueURLSplit)-1]
	switch *n.Subject.Type {
	case "Issue":
		base = fmt.Sprintf("%s/issues/%s", base, issueID)
	case "PullRequest":
		base = fmt.Sprintf("%s/pull/%s", base, issueID)
	}

	if n.Subject.LatestCommentURL != nil {
		commentURLSplit := strings.Split(*n.Subject.LatestCommentURL, "/")
		commentID := commentURLSplit[len(commentURLSplit)-1]
		switch *n.Reason {
		case "review_requested":
			base = fmt.Sprintf("%s#issuecomment-%s", base, commentID)
			// TODO: Find a way to point to the review comment
			//base = fmt.Sprintf("%s#discussion_r%s", base, commentID)
		case "subscribed":
			base = fmt.Sprintf("%s#issuecomment-%s", base, commentID)
		}
	}

	return url.Parse(base)
}

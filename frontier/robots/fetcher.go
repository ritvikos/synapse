package robots

import (
	"context"
	"net/http"
	"time"

	"github.com/temoto/robotstxt"
)

type DefaultRobotsFetcher struct {
	client http.Client
	data   *robotstxt.RobotsData
}

func NewDefaultRobotsTxtFetcher(client http.Client) *DefaultRobotsFetcher {
	return &DefaultRobotsFetcher{
		client: http.Client{
			Timeout: 12 * time.Second,
		},
	}
}

func (r *DefaultRobotsFetcher) Fetch(ctx context.Context, host string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", host+"/robots.txt", nil)
	if err != nil {
		return nil, err
	}

	return r.client.Do(req)
}

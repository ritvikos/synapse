// Copyright 2025-2026 Ritvik Gupta
// SPDX-License-Identifier: Apache-2.0

package robots

import (
	"context"
	"net/http"
	"time"

	"github.com/temoto/robotstxt"
)

type RobotsFetcher interface {
	Fetch(ctx context.Context, host string) (*http.Response, error)
}

// TODO: Add TTL
type RobotsEntry struct {
	Group       *robotstxt.Group
	LastFetched time.Time
}

func (e *RobotsEntry) Test(path string) bool {
	if e.Group == nil {
		return true
	}
	return e.Group.Test(path)
}

func (e *RobotsEntry) CrawlDelay() time.Duration {
	if e.Group == nil {
		return time.Second * 0
	}
	return e.Group.CrawlDelay
}

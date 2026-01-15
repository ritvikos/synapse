// Copyright 2025-2026 Ritvik Gupta
// SPDX-License-Identifier: Apache-2.0

package robots

import (
	"context"
	"time"

	"github.com/ritvikos/synapse/frontier/backend"
	"github.com/temoto/robotstxt"
	"golang.org/x/sync/singleflight"
)

type RobotsBackend = backend.Store[*RobotsEntry]

type RobotsResolver struct {
	fetcher   RobotsFetcher
	backend   RobotsBackend
	sf        singleflight.Group
	userAgent string
}

func NewRobotsResolver(
	userAgent string,
	fetcher RobotsFetcher,
	backend RobotsBackend,
) *RobotsResolver {
	return &RobotsResolver{
		userAgent: userAgent,
		fetcher:   fetcher,
		backend:   backend,
		sf:        singleflight.Group{},
	}
}

// origin: [Scheme + Host]
func (r *RobotsResolver) Resolve(ctx context.Context, origin string) (*RobotsEntry, error) {
	entry, err := r.backend.Get(ctx, origin)
	if err == nil {
		return entry, nil
	}

	result, err, _ := r.sf.Do(origin, func() (any, error) {
		if e, err := r.backend.Get(ctx, origin); err == nil {
			return e, nil
		}

		data, err := r.resolve(ctx, origin)
		if err != nil {
			return nil, err
		}

		entry := &RobotsEntry{
			Group:       data.FindGroup(r.userAgent),
			LastFetched: time.Now(),
		}

		return entry, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*RobotsEntry), nil
}

func (r *RobotsResolver) resolve(ctx context.Context, origin string) (*robotstxt.RobotsData, error) {
	resp, err := r.fetcher.Fetch(ctx, origin)
	if err != nil {
		return nil, err
	}

	data, err := robotstxt.FromResponse(resp)
	if err != nil {
		return nil, err
	}

	entry := &RobotsEntry{
		Group:       data.FindGroup(r.userAgent),
		LastFetched: time.Now(),
	}

	if err := r.backend.Put(ctx, origin, entry); err != nil {
		return nil, err
	}

	return data, nil
}

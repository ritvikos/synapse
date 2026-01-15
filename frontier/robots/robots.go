// Copyright 2025-2026 Ritvik Gupta
// SPDX-License-Identifier: Apache-2.0

package robots

import (
	"context"
	"errors"
	"time"

	"github.com/ritvikos/synapse/frontier/backend"
	"github.com/temoto/robotstxt"
	"golang.org/x/sync/singleflight"
)

// ------ NOTE ------
// Currently, the [RobotsResolver] uses origin (url scheme + host) as the cache key.
// As soon as the fingerprint usage is clarified in the overall architecture,
// the internal usage will be changed. The cache will store entries based on
// fingerprint (string) instead of origin (string).

var ErrRobotsInvalidTTL = errors.New("robots resolver: invalid TTL duration")

type RobotsCache = backend.Cache[*RobotsEntry]

type RobotsResolver struct {
	sf        singleflight.Group
	fetcher   RobotsFetcher
	cache     RobotsCache
	userAgent string
	ttl       time.Duration
}

func NewRobotsResolver(
	userAgent string,
	fetcher RobotsFetcher,
	cache RobotsCache,
	ttl time.Duration,
) (*RobotsResolver, error) {
	if ttl <= 0 {
		return nil, ErrRobotsInvalidTTL
	}

	return &RobotsResolver{
		userAgent: userAgent,
		fetcher:   fetcher,
		cache:     cache,
		sf:        singleflight.Group{},
		ttl:       ttl,
	}, nil
}

func (r *RobotsResolver) Resolve(ctx context.Context, origin string) (*RobotsEntry, error) {
	return r.resolve(ctx, origin, r.ttl)
}

func (r *RobotsResolver) ResolveWithTTL(ctx context.Context, origin string, ttl time.Duration) (*RobotsEntry, error) {
	if ttl <= 0 {
		return nil, ErrRobotsInvalidTTL
	}
	return r.resolve(ctx, origin, ttl)
}

// origin: [Scheme + Host]
func (r *RobotsResolver) resolve(ctx context.Context, origin string, ttl time.Duration) (*RobotsEntry, error) {
	entry, err := r.cache.Get(ctx, origin)
	if err == nil {
		return entry, nil
	}

	result, err, _ := r.sf.Do(origin, func() (any, error) {
		if e, err := r.cache.Get(ctx, origin); err == nil {
			return e, nil
		}

		entry, err := r.get(ctx, origin)
		if err != nil {
			return nil, err
		}

		if err := r.cache.Set(ctx, origin, entry, ttl); err != nil {
			return nil, err
		}

		return entry, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*RobotsEntry), nil
}

// Fetch and parse 'robots.txt' from origin
func (r *RobotsResolver) get(ctx context.Context, origin string) (*RobotsEntry, error) {
	resp, err := r.fetcher.Fetch(ctx, origin)
	if err != nil {
		return nil, err
	}

	data, err := robotstxt.FromResponse(resp)
	if err != nil {
		return nil, err
	}

	return &RobotsEntry{
		Group:       data.FindGroup(r.userAgent),
		LastFetched: time.Now(),
	}, nil
}

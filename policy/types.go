// Copyright 2025-2026 Ritvik Gupta
// SPDX-License-Identifier: Apache-2.0

package policy

import (
	"context"
	"net/http"
	"time"
)

type BackoffPolicy interface {
	NextRetry(attempt int) time.Duration
}

type RetryPolicy interface {
	Next() (time.Duration, bool)
}

type RateLimitPolicy interface {
	Wait(ctx context.Context, req *http.Request) error
}

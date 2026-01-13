// Copyright 2025-2026 Ritvik Gupta
// SPDX-License-Identifier: Apache-2.0

package lifecycle

import "context"

type Lifecycle interface {
	Start(ctx context.Context) error

	Stop(ctx context.Context) error
}

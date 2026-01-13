// Copyright 2025-2026 Ritvik Gupta
// SPDX-License-Identifier: Apache-2.0

package http

import "io"

type readCloser struct {
	io.Reader
	closer io.Closer
}

func (rc *readCloser) Close() error {
	if rc.closer != nil {
		return rc.closer.Close()
	}
	return nil
}

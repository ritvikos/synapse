// Copyright 2025-2026 Ritvik Gupta
// SPDX-License-Identifier: Apache-2.0

package pipeline

import "io"

type Handler interface {
	Execute(r io.Reader) error
	ContentType() string
}

type Parser[T any] interface {
	Parse(r io.Reader) (T, error)
}

type Processor[T any] interface {
	Process(data T) (T, error)
}

type Sink[T any] interface {
	Write(data T) error
	Close() error
}

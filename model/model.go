package model

import (
	"time"
)

type Task[T any] struct {
	Url       string
	ExecuteAt time.Time
	Score     float64
	Metadata  T
	// Fingerprint string
}

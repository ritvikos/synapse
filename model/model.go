package model

import (
	"time"
)

type Task[T any] struct {
	Url       string
	ExecuteAt time.Time
	Metadata  T
	// Fingerprint string
}

type ScoredTask[T any] struct {
	Task  *Task[T]
	Score float64
}

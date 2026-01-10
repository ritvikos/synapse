package backend

import (
	"context"
)

// Generic Queue interface
type Queue[T any] interface {
	// Insert items into the queue.
	Enqueue(ctx context.Context, items []T) error

	// Retrieve items from the queue into the provided buffer.
	// Returns the number of items retrieved.
	// If the queue is empty, it returns (0, nil)
	Dequeue(ctx context.Context, n int, buf chan<- T) (int, error)

	// Number of pending items in the queue.
	Len(ctx context.Context) (int, error)
}

type Store[T any] interface {
	Put(ctx context.Context, key string, value T) error
	Get(ctx context.Context, key string) (T, error)
	Delete(ctx context.Context, key string) error
}

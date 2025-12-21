package backend

import (
	"context"
)

// Generic Queue interface
type Queue[T any] interface {
	// Insert items into the queue.
	Enqueue(ctx context.Context, items ...T) error

	// Retrieve upto 'n' items from the queue.
	Dequeue(ctx context.Context, n int) ([]T, error)

	// Number of pending items in the queue.
	Len(ctx context.Context) (int, error)
}

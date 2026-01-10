package sched

import (
	"context"

	"github.com/ritvikos/synapse/frontier/backend"
	"github.com/ritvikos/synapse/model"
)

type ScoredTask[T any] = *model.ScoredTask[T]
type Queue[T any] = backend.Queue[ScoredTask[T]]

type Scheduler[T any] interface {
	// TODO: Replace Start() and Stop() with lifecycle interface, once fixed :p
	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	// Enqueues a task into the scheduler,
	// eventally to be flushed to the underlying queue.
	Enqueue(ctx context.Context, task ScoredTask[T]) error

	// Dequeues 'n' tasks from the underlying queue into the scheduler's buffer
	// (if any, based on underlying implementation) and returns the number of tasks dequeued.
	Dequeue(ctx context.Context) ScoredTask[T]
}

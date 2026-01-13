// Copyright 2025-2026 Ritvik Gupta
// SPDX-License-Identifier: Apache-2.0

package sched

import (
	"context"
	"fmt"
	"sync"
)

var _ Scheduler[any] = (*UnbufferedScheduler[any])(nil)

// Directly enqueue/dequeue to/from the backend queue without any buffering.
type UnbufferedScheduler[T any] struct {
	queue     Queue[T]
	dequeueCh chan ScoredTask[T]

	// Internal
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
}

func NewUnbufferedScheduler[T any](queue Queue[T]) *UnbufferedScheduler[T] {
	return &UnbufferedScheduler[T]{
		queue:     queue,
		dequeueCh: make(chan ScoredTask[T], 1),
	}
}

func (s *UnbufferedScheduler[T]) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancel != nil {
		return fmt.Errorf("[blocking scheduler]: already started")
	}
	s.ctx, s.cancel = context.WithCancel(ctx)
	return nil
}

func (s *UnbufferedScheduler[T]) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancel == nil {
		return fmt.Errorf("[blocking scheduler]: not started")
	}
	s.cancel()
	s.cancel = nil
	return nil
}

func (s *UnbufferedScheduler[T]) Dequeue(ctx context.Context) ScoredTask[T] {
	n, _ := s.queue.Dequeue(ctx, 1, s.dequeueCh)
	if n == 0 {
		return nil
	}

	select {
	case task := <-s.dequeueCh:
		return task
	case <-ctx.Done():
		return nil
	default:
	}

	return nil
}

func (s *UnbufferedScheduler[T]) Enqueue(ctx context.Context, task ScoredTask[T]) error {
	return s.queue.Enqueue(s.ctx, []ScoredTask[T]{task})
}

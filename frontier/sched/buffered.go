// Copyright 2025-2026 Ritvik Gupta
// SPDX-License-Identifier: Apache-2.0

package sched

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// TODO: Support on-disk persistence (configurable) for prefetch/flush buffers
// for recovery and prevent task loss.

type BufferedScheduler[T any] struct {
	queue  Queue[T]
	policy BufferPolicy

	prefetchChan       chan ScoredTask[T]
	prefetchSignalChan chan struct{}

	flushChan       chan ScoredTask[T]
	flushSignalChan chan struct{}
	flushTimer      *time.Timer

	// Internal
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
	wg     sync.WaitGroup
}

func NewBufferedScheduler[T any](
	queue Queue[T],
	policy BufferPolicy,
	prefetchBufSize uint,
	flushBufSize uint,
) *BufferedScheduler[T] {
	return &BufferedScheduler[T]{
		queue:              queue,
		policy:             policy,
		prefetchChan:       make(chan ScoredTask[T], prefetchBufSize),
		prefetchSignalChan: make(chan struct{}, 1),
		flushChan:          make(chan ScoredTask[T], flushBufSize),
		flushSignalChan:    make(chan struct{}, 1),
	}
}

func (s *BufferedScheduler[T]) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.cancel != nil {
		s.mu.Unlock()
		return fmt.Errorf("[buffered scheduler]: already started")
	}

	s.ctx, s.cancel = context.WithCancel(ctx)
	s.mu.Unlock()

	s.wg.Add(2)
	go s.prefetchWorker()

	s.triggerPrefetch()

	return nil
}

func (s *BufferedScheduler[T]) Stop(ctx context.Context) error {
	s.mu.Lock()
	if s.cancel == nil {
		s.mu.Unlock()
		return fmt.Errorf("[buffered scheduler]: not started")
	}

	s.cancel()
	s.cancel = nil
	s.mu.Unlock()

	s.wg.Wait()
	return nil
}

func (s *BufferedScheduler[T]) Dequeue(ctx context.Context) ScoredTask[T] {
	// fast path
	select {
	case task, ok := <-s.prefetchChan:
		if !ok {
			return nil
		}
		return task
	default:
	}

	// slow path
	s.triggerPrefetch()

	select {
	case <-ctx.Done():
		return nil
	case <-s.ctx.Done():
		return nil
	case task, ok := <-s.prefetchChan:
		if !ok {
			return nil
		}
		return task
	}
}

func (s *BufferedScheduler[T]) Enqueue(ctx context.Context, task ScoredTask[T]) error {
	// fast path
	select {
	case s.flushChan <- task:
		return nil
	default:
	}

	// slow path
	select {
	case s.flushSignalChan <- struct{}{}:
	default:
	}

	select {
	case <-s.ctx.Done():
		return s.ctx.Err()
	case s.flushChan <- task:
		return nil
	}
}

func (s *BufferedScheduler[T]) prefetchWorker() {
	defer s.wg.Done()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.prefetchSignalChan:
			for {
				state := State{
					BufLen: len(s.prefetchChan),
					BufCap: cap(s.prefetchChan),
				}

				count := s.policy.Prefetch(state)
				if count <= 0 {
					break
				}

				// TODO: backoff/retry on error
				n, err := s.queue.Dequeue(s.ctx, count, s.prefetchChan)
				if err != nil {
					log.Printf("[buffered scheduler]: prefetch dequeue error: %v", err)
					break
				}
				if n == 0 {
					break
				}
			}
		}
	}
}

func (s *BufferedScheduler[T]) triggerPrefetch() {
	select {
	case s.prefetchSignalChan <- struct{}{}:
	default:
	}
}

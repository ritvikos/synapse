// Copyright 2025-2026 Ritvik Gupta
// SPDX-License-Identifier: Apache-2.0

package pipeline

import "io"

type LocalPipelineOption[T any] func(*LocalPipeline[T])

func WithProcessors[T any](procs ...Processor[T]) LocalPipelineOption[T] {
	return func(p *LocalPipeline[T]) {
		p.processors = append(p.processors, procs...)
	}
}

func WithSink[T any](sink Sink[T]) LocalPipelineOption[T] {
	return func(p *LocalPipeline[T]) {
		p.sink = sink
	}
}

func WithExecutor[T any](executor func(*ExecutionContext[T], io.Reader) error) LocalPipelineOption[T] {
	return func(p *LocalPipeline[T]) {
		p.executor = executor
	}
}

func (p LocalPipeline[T]) Execute(r io.Reader) error {
	return p.executor(&ExecutionContext[T]{
		parser:     p.parser,
		processors: p.processors,
		sink:       p.sink,
	}, r)
}

func (p LocalPipeline[T]) ContentType() string {
	return p.contentType
}

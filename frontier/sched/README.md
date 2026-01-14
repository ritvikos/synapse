# Scheduler

## Purpose

Instead of pinging the underlying queue backend for every enqueue/dequeue op, which is expensive and inefficient especially when it isn't co-located with the frontier. To address this, the [`Scheduler`](./sched.go) is an abstraction that uses pluggable [`BufferPolicy`](./types.go) to determine when URLs should be fetched/flushed from/to the underlying `Queue`.

Internally, it provides two pluggable implementations (based on different scaling requirements):

1. [**Buffered Scheduler**](./buffered.go) batches urls locally before enqueue/dequeue is performed. It maintains internal buffers for prefetching/flushing urls from/to the underlying queue. Once the buffer limits are reached (based on the configured [`BufferPolicy`](./types.go)), it performs bulk enqueue/dequeue operations to/from the backend queue.

2. [**Unbuffered Scheduler**](./unbuffered.go) provides a direct pass-through to the backend. Basically, it synchronously performs enqueue/dequeue operations to/from the underlying queue without any intermediate buffering. This is useful when the queue is co-located with the frontier.

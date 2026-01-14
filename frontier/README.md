# Frontier

## Purpose

Internally, it orchestrates the following sub-components that handle robots.txt compliance, URL prioritization, scheduling and exposes a facade for the end-user to enqueue/dequeue urls for crawling:

1. [**Robots Resolver**](./robots/) fetches `robots.txt` when needed and enforces compliance with the [Robots Exclusion Protocol](https://en.wikipedia.org/wiki/Robots.txt). It resolves `robots.txt` for target hosts to apply crawl-delay directives and path-based exclusions, persist in storage (configured by the end-user with [`Store`](./backend/types.go) interface).

2. [**Scheduler**](./sched/) There're two types of pluggable schedulers:

   1. [**Buffered Scheduler**](./sched/buffered.go) which handles the buffering of the scored tasks to be enqueued/dequeued to/from the underlying [`Queue`](./backend/types.go) backend.

   2. [**Unbuffered Scheduler**](./sched/unbuffered.go) which directly interacts with the underlying [`Queue`](./backend/types.go) backend without any intermediate buffering.

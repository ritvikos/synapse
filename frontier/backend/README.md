# Backend

## Purpose

It provides the fundamental abstraction layer for data persistence and queuing via [`Store`](./types.go) and [`Queue`](./types.go) interfaces, respectively. This design allows the [`Frontier`](../frontier.go) to operate independently of the underlying storage infra, allowing pluggable backends.

1. [**Queue**](./types.go) is a generic interface for FIFO operations (`Enqueue`, `Dequeue`, `Len`).

2. [**Store**](./types.go) is a generic key-value interface (`Put`, `Get`, `Delete`).

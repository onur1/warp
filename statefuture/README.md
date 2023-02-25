# statefuture

A StateFuture represents a value or an error that exists in the future and depends on itself through some computation.

```go
type StateFuture[S, A any] func(s S) Future[These[A, S]]
```

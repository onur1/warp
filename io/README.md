# io

An IO represents the result of a non-deterministic computation that may cause side-effects, but never fails and yields a value of type A.

```go
type IO[A any] func() A
```

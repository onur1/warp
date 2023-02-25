# io
An IO is a computation which, when performed, does some I/O before returning a value of type A.

```go
type IO[A any] func() A
```

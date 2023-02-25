# state

A State represents a computation that depend on and modify some internal state where parameter `S` is a state type to carry and `A` is the type of a return value.

```go
type State[S, A any] func(s S) (A, S)
```

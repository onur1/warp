# state

A State represents a value which depends on itself through some computation, where parameter `S` is a state type to carry and `A` is the type of a return value.

```go
type State[S, A any] func(s S) (A, S)
```

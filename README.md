# data

Package data provides experimental monadic data types in Go.

## Types

An IO represents a computation that never fails and yields a value of type A.

```go
type IO[A any] func() A
```

A Result represents a result of a computation which is either a value of type A,
or an error.

```go
type Result[A any] func() (A, error)
```

A Nilable represents an optional value which is either some value or nil.

```go
type Nilable[A any] *A
```

An Event represents a collection of discrete occurrences of events with associated
values.

```go
type Event[A any] func(context.Context, chan<- A)
```

A Future represents a collection of discrete occurrences of events with associated
values or errors, in that, a Future is actually an Event that may fail and emits
a value which is encapsulated in a Result.

```go
type Future[A any] Event[Result[A]]
```

A State represents a value which depends on itself through some computation, where
parameter S is the state type to carry and A is the type of a return value.

```go
type State[S, A any] func(s S) (A, S)
```

// Package fpgo provides experimental monads.
package fpgo

import (
	"context"
)

// An IO represents a computation that never fails and yields a value of type A.
type IO[A any] func() A

// A Result represents a result of a computation which is either a value of type A,
// or an error.
type Result[A any] func(context.Context) (A, error)

// An Event represents a collection of discrete occurrences of events with associated
// values.
type Event[A any] func(context.Context, chan<- A)

// A Future represents a collection of discrete occurrences of events with associated
// values or errors, in that, a Future is actually an Event that may fail and emits
// a value which is encapsulated in a Result.
type Future[A any] Event[Result[A]]

// A Predicate represents a predicate (boolean-valued function) of one argument.
type Predicate[A any] func(A) bool

// A Nilable represents an optional value which is either some value or nil.
type Nilable[A any] *A

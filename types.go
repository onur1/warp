package data

import (
	"context"
)

// An Either represents a value pair which contains two values that can never
// co-exist, either the left one the right one will have zero value.
type Either[E, A comparable] func() (E, A)

// A These represents a value pair which contains two values that may have zero
// values or not. This is "inclusive-or" (as opposed to "exclusive-or" provided
// by Either), both values can have zero values, or only one of them, or both of
// them can have non-zero values.
type These[E, A comparable] func() (E, A)

// A Result represents a result of a computation that either yields a value
// of type A, or fails with an error.
type Result[A any] func() (A, error)

// A Nilable represents an optional value which is either some value or nil.
type Nilable[A any] *A

// An Event represents a collection of discrete occurrences with associated values.
type Event[A any] func(context.Context, chan<- A)

// A Future represents an Event that may fail, in that, it returns a value which is
// encapsulated in a Result.
type Future[A any] Event[Result[A]]

// A State represents a value which depends on itself through some computation.
type State[S, A any] func(s S) (A, S)

// A StateFuture represents a value or an error that exists in the future and
// depends on itself through some computation.
type StateFuture[S, A comparable] func(s S) Future[These[A, S]]

// An IO represents the result of a non-deterministic computation that may cause
// side-effects, but never fails and yields a value of type A.
type IO[A any] func() A

type Predicate[A any] func(A) bool

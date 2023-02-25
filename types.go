package data

import (
	"context"
)

// A Result represents a result of a computation that either yields a value
// of type A, or fails with an error.
type Result[A any] func() (A, error)

// A Nilable represents an optional value which is either some value or nil.
type Nilable[A any] *A

// An Event represents a collection of discrete occurrences with associated values.
type Event[A any] func(context.Context, chan<- A)

// A Future represents a collection of asynchronous computations with
// their associated results.
type Future[A any] Event[Result[A]]

type Predicate[A any] func(A) bool

type Lazy[A any] func() A

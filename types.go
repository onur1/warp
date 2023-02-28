// Package data provides a set of useful data types for working with different
// types of data, such as data that contains an error, data that doesn't exist,
// data that doesn't exist yet, data that keeps changing, and etc.
//
// The practical part of [functional programming] is all about higher level
// abstractions (like these types that describe some nature of data) and combining
// their behaviors to gain properties and functionality to deal with data in more
// detail.
//
// Data types must obey certain laws to implement a [typeclass] instance (which
// contains implementations of functions defined in a typeclass, such as [functors, applicatives and monads]),
// specialized to a particular type. These laws are rooted in [category theory]
// and together they form a type system that ensures safety and composability.
//
// [functional programming]: https://github.com/enricopolanski/functional-programming
// [typeclass]: https://wiki.haskell.org/Typeclassopedia
// [functors, applicatives and monads]: https://www.adit.io/posts/2013-04-17-functors,_applicatives,_and_monads_in_pictures.html
// [category theory]: https://www.youtube.com/watch?v=gui_SE8rJUM
package data

import (
	"context"
)

// An IO represents a computation that never fails and yields a value of type A.
type IO[A any] func() A

// A Result represents a result of a computation which is either a value of type A,
// or an error.
type Result[A any] func() (A, error)

// A Nilable represents an optional value which is either some value or nil.
type Nilable[A any] *A

// An Event represents a collection of discrete occurrences of events with associated
// values.
type Event[A any] func(context.Context, chan<- A)

// A Future represents a collection of discrete occurrences of events with associated
// values or errors, in that, a Future is actually an Event that may fail and emits
// a value which is encapsulated in a Result.
type Future[A any] Event[Result[A]]

// A State represents a value which depends on itself through some computation, where
// parameter S is the state type to carry and A is the type of a return value.
type State[S, A any] func(s S) (A, S)

// A Predicate represents a predicate (boolean-valued function) of one argument.
type Predicate[A any] func(A) bool

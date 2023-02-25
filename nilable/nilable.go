// Package nilable implements the Nilable type.
package nilable

import (
	"github.com/onur1/data"
)

// IsNil returns true if the value is nil.
func IsNil[A any](a data.Nilable[A]) bool {
	return a == nil
}

// IsSome returns true if the value is not nil.
func IsSome[A any](a data.Nilable[A]) bool {
	return a != nil
}

// Nil creates a nilable with nil value.
func Nil[A any]() data.Nilable[A] {
	return nil
}

// Some creates a nilable with some value.
func Some[A any](a A) data.Nilable[A] {
	return &a
}

// Map creates a nilable by applying a function on an existing value.
func Map[A, B any](fa data.Nilable[A], f func(A) B) data.Nilable[B] {
	if fa == nil {
		return nil
	}
	return Some(f(*fa))
}

// Ap creates a nilable by applying a function contained in the first nilable on
// the value contained in the second nilable if they both exist.
func Ap[A, B any](fab data.Nilable[func(A) B], fa data.Nilable[A]) data.Nilable[B] {
	if fab == nil {
		return nil
	}
	if fa == nil {
		return nil
	}
	return Some((*fab)(*fa))
}

// Chain creates a nilable which combines two nilables in sequence, using the
// return value of one nilable to determine the next one.
func Chain[A, B any](ma data.Nilable[A], f func(A) data.Nilable[B]) data.Nilable[B] {
	if ma == nil {
		return nil
	}
	return f(*ma)
}

// ApFirst creates a nilable by combining two nilables, keeping only the result
// of the first.
func ApFirst[A, B any](fa data.Nilable[A], fb data.Nilable[B]) data.Nilable[A] {
	return Ap(Map(fa, func(a A) func(B) A {
		return func(_ B) A {
			return a
		}
	}), fb)
}

// ApSecond creates a nilable by combining two nilables, keeping only the result
// of the second.
func ApSecond[A, B any](fa data.Nilable[A], fb data.Nilable[B]) data.Nilable[B] {
	return Ap(Map(fa, func(_ A) func(B) B {
		return func(b B) B {
			return b
		}
	}), fb)
}

// FromResult creates a nilable from a result, returning nil for errors.
func FromResult[A any](ma data.Result[A]) data.Nilable[A] {
	a, err := ma()
	if err != nil {
		return nil
	}
	return Some(a)
}

// FromPredicate creates a nilable by testing a value against a predicate first.
func FromPredicate[A any](a A, predicate data.Predicate[A]) data.Nilable[A] {
	if predicate(a) {
		return Some(a)
	} else {
		return Nil[A]()
	}
}

// Attempt creates a nilable by running a function which returns a value,
// recovering with nil if a panic is thrown.
func Attempt[A any](f data.IO[A]) data.Nilable[A] {
	defer func() {
		recover()
	}()
	return Some(f())
}

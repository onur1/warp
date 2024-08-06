// Package nilable implements the Nilable type.
package nilable

import (
	"context"

	"github.com/onur1/warp"
)

// IsNil returns true if the value is nil.
func IsNil[A any](a warp.Nilable[A]) bool {
	return a == nil
}

// IsSome returns true if the value is not nil.
func IsSome[A any](a warp.Nilable[A]) bool {
	return a != nil
}

// Nil creates a nilable with nil value.
func Nil[A any]() warp.Nilable[A] {
	return nil
}

// Some creates a nilable with some value.
func Some[A any](a A) warp.Nilable[A] {
	return &a
}

// Map creates a nilable by applying a function on an existing value.
func Map[A, B any](fa warp.Nilable[A], f func(A) B) warp.Nilable[B] {
	if fa == nil {
		return nil
	}
	return Some(f(*fa))
}

// Ap creates a nilable by applying a function contained in the first nilable on
// the value contained in the second nilable if they both exist.
func Ap[A, B any](fab warp.Nilable[func(A) B], fa warp.Nilable[A]) warp.Nilable[B] {
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
func Chain[A, B any](ma warp.Nilable[A], f func(A) warp.Nilable[B]) warp.Nilable[B] {
	if ma == nil {
		return nil
	}
	return f(*ma)
}

// ApFirst creates a nilable by combining two nilables, keeping only the result
// of the first.
func ApFirst[A, B any](fa warp.Nilable[A], fb warp.Nilable[B]) warp.Nilable[A] {
	return Ap(Map(fa, func(a A) func(B) A {
		return func(_ B) A {
			return a
		}
	}), fb)
}

// ApSecond creates a nilable by combining two nilables, keeping only the result
// of the second.
func ApSecond[A, B any](fa warp.Nilable[A], fb warp.Nilable[B]) warp.Nilable[B] {
	return Ap(Map(fa, func(_ A) func(B) B {
		return func(b B) B {
			return b
		}
	}), fb)
}

// FromResult creates a nilable from a result, returning nil for errors.
func FromResult[A any](ctx context.Context, ma warp.Result[A]) warp.Nilable[A] {
	a, err := ma(ctx)
	if err != nil {
		return nil
	}
	return Some(a)
}

// FromPredicate creates a nilable by testing a value against a predicate first.
func FromPredicate[A any](a A, predicate warp.Predicate[A]) warp.Nilable[A] {
	if predicate(a) {
		return Some(a)
	} else {
		return Nil[A]()
	}
}

// Attempt creates a nilable by running a function which returns a value,
// recovering with nil if a panic is thrown.
func Attempt[A any](f warp.IO[A]) warp.Nilable[A] {
	defer func() {
		recover()
	}()
	return Some(f())
}

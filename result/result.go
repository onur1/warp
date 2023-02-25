// Package result implements the Result type.
package result

import (
	"github.com/onur1/data"
)

// Succeed creates a result which always returns a value.
func Succeed[A any](a A) data.Result[A] {
	return func() (A, error) {
		return a, nil
	}
}

// Fail creates a result which always fails with an error.
func Fail[A any](err error) data.Result[A] {
	return func() (a A, _ error) {
		return a, err
	}
}

// Map creates a result by applying a function on a succeeding result.
func Map[A, B any](fa data.Result[A], f func(A) B) data.Result[B] {
	a, err := fa()
	if err != nil {
		return Fail[B](err)
	}
	return Succeed(f(a))
}

// MapError creates a result by applying a function on a failing result.
func MapError[A any](fa data.Result[A], f func(error) error) data.Result[A] {
	_, err := fa()
	if err != nil {
		return Fail[A](f(err))
	}
	return fa
}

// Ap creates a result by applying a function contained in the first result
// on the value contained in the second result.
func Ap[A, B any](fab data.Result[func(A) B], fa data.Result[A]) data.Result[B] {
	var (
		err error
		ab  func(A) B
		a   A
	)

	ab, err = fab()
	if err != nil {
		return Fail[B](err)
	}

	a, err = fa()
	if err != nil {
		return Fail[B](err)
	}

	return Succeed(ab(a))
}

// Chain creates a result which combines two results in sequence, using the
// return value of one result to determine the next one.
func Chain[A, B any](ma data.Result[A], f func(A) data.Result[B]) data.Result[B] {
	a, err := ma()
	if err != nil {
		return Fail[B](err)
	}
	return f(a)
}

// Bimap creates a result by mapping a pair of functions over an error or a value
// contained in a result.
func Bimap[A, B any](fa data.Result[A], f func(error) error, g func(A) B) data.Result[B] {
	a, err := fa()
	if err != nil {
		return Fail[B](f(err))
	}
	return Succeed(g(a))
}

// ApFirst creates a result by combining two effectful computations, keeping
// only the result of the first.
func ApFirst[A, B any](fa data.Result[A], fb data.Result[B]) data.Result[A] {
	return Ap(Map(fa, func(a A) func(B) A {
		return func(_ B) A {
			return a
		}
	}), fb)
}

// ApSecond creates a result by combining two effectful computations, keeping
// only the result of the second.
func ApSecond[A, B any](fa data.Result[A], fb data.Result[B]) data.Result[B] {
	return Ap(Map(fa, func(_ A) func(B) B {
		return func(b B) B {
			return b
		}
	}), fb)
}

// Fold takes two functions and a result and returns a value by applying
// one of the supplied functions to the inner value.
func Fold[A, B any](ma data.Result[A], onError func(error) B, onSuccess func(A) B) B {
	a, err := ma()
	if err != nil {
		return onError(err)
	}
	return onSuccess(a)
}

// FromNilable creates a result from a nilable, returning the supplied error
// for nil values.
func FromNilable[A any](ma data.Nilable[A], onNil func() error) data.Result[A] {
	if ma == nil {
		return Fail[A](onNil())
	}
	return Succeed(*ma)
}

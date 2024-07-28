// Package result implements the Result type.
package result

import (
	"context"
	"time"

	"github.com/onur1/gofp"
)

// Ok creates a result which never fails and returns a value of type A.
func Ok[A any](a A) gofp.Result[A] {
	return func(_ context.Context) (A, error) {
		return a, nil
	}
}

func After[A any](dur time.Duration, a A) gofp.Result[A] {
	return func(_ context.Context) (A, error) {
		time.Sleep(dur)
		return a, nil
	}
}

func ErrorAfter[A any](dur time.Duration, err error) gofp.Result[A] {
	return func(_ context.Context) (a A, _ error) {
		time.Sleep(dur)
		return a, err
	}
}

// Error creates a result which always fails with an error.
func Error[A any](err error) gofp.Result[A] {
	return func(_ context.Context) (a A, _ error) {
		return a, err
	}
}

// Zero creates a result which never fails and returns a zero value
// of the type that it is initialized with.
func Zero[A any]() gofp.Result[A] {
	return func(_ context.Context) (a A, _ error) {
		return
	}
}

// Map creates a result by applying a function on a succeeding result.
func Map[A, B any](fa gofp.Result[A], f func(A) B) gofp.Result[B] {
	return func(ctx context.Context) (b B, err error) {
		var a A
		if a, err = fa(ctx); err != nil {
			return
		}
		b = f(a)
		return
	}
}

// MapError creates a result by applying a function on a failing result.
func MapError[A any](fa gofp.Result[A], f func(error) error) gofp.Result[A] {
	return func(ctx context.Context) (a A, err error) {
		if a, err = fa(ctx); err != nil {
			err = f(err)
			return
		}
		return
	}
}

// Ap creates a result by applying a function contained in the first result
// on the value contained in the second result.
func Ap[A, B any](fab gofp.Result[func(A) B], fa gofp.Result[A]) gofp.Result[B] {
	return func(ctx context.Context) (b B, err error) {
		var ab func(A) B

		if ab, err = fab(ctx); err != nil {
			return
		}

		var a A

		if a, err = fa(ctx); err != nil {
			return
		}

		b = ab(a)

		return
	}
}

// Chain creates a result which combines two results in sequence, using the
// return value of one result to determine the next one.
func Chain[A, B any](ma gofp.Result[A], f func(A) gofp.Result[B]) gofp.Result[B] {
	return func(ctx context.Context) (_ B, err error) {
		var a A
		if a, err = ma(ctx); err != nil {
			return
		}
		return f(a)(ctx)
	}
}

// ChainFirst composes two results in sequence, using the return value of one result
// to determine the next one, keeping only the first result.
func ChainFirst[A, B any](ma gofp.Result[A], f func(A) gofp.Result[B]) gofp.Result[A] {
	return Chain(ma, func(a A) gofp.Result[A] {
		return Map(f(a), fst[A, B](a))
	})
}

// Bimap creates a result by mapping a pair of functions over an error or a value
// contained in a result.
func Bimap[A, B any](fa gofp.Result[A], f func(error) error, g func(A) B) gofp.Result[B] {
	return func(ctx context.Context) (b B, err error) {
		var a A
		if a, err = fa(ctx); err != nil {
			err = f(err)
		} else {
			b = g(a)
		}
		return
	}
}

// ApFirst creates a result by combining two effectful computations, keeping
// only the result of the first.
func ApFirst[A, B any](fa gofp.Result[A], fb gofp.Result[B]) gofp.Result[A] {
	return Ap(Map(fa, fst[A, B]), fb)
}

// ApSecond creates a result by combining two effectful computations, keeping
// only the result of the second.
func ApSecond[A, B any](fa gofp.Result[A], fb gofp.Result[B]) gofp.Result[B] {
	return Ap(Map(fa, snd[A, B]), fb)
}

// Reduce takes two functions and a result and returns a value by applying
// one of the supplied functions to the inner value.
func Reduce[A, B any](ctx context.Context, ma gofp.Result[A], onError func(error) B, onSuccess func(A) B) B {
	if a, err := ma(ctx); err != nil {
		return onError(err)
	} else {
		return onSuccess(a)
	}
}

// GetOrElse creates a result which can be used to recover from a failing result
// with a new value.
func GetOrElse[A any](ctx context.Context, ma gofp.Result[A], onError func(error) A) A {
	if a, err := ma(ctx); err != nil {
		return onError(err)
	} else {
		return a
	}
}

// OrElse creates a result which can be used to recover from a failing result
// by switching to a new result.
func OrElse[A any](ma gofp.Result[A], onError func(error) gofp.Result[A]) gofp.Result[A] {
	return func(ctx context.Context) (a A, err error) {
		if a, err = ma(ctx); err != nil {
			return onError(err)(ctx)
		}
		return
	}
}

// FilterOrElse creates a result which can be used to fail with an error unless
// a predicate holds on a succeeding result.
func FilterOrElse[A any](ma gofp.Result[A], predicate gofp.Predicate[A], onFalse func(A) error) gofp.Result[A] {
	return Chain(ma, func(a A) gofp.Result[A] {
		if predicate(a) {
			return Ok(a)
		} else {
			return Error[A](onFalse(a))
		}
	})
}

// Fork is like Reduce but it doesn't have a return value.
func Fork[A any](ctx context.Context, ma gofp.Result[A], onError func(error), onSuccess func(A)) {
	if a, err := ma(ctx); err != nil {
		onError(err)
	} else {
		onSuccess(a)
	}
}

// FromNilable creates a result from a nilable, returning the supplied error
// for nil values.
func FromNilable[A any](ma gofp.Nilable[A], onNil func() error) gofp.Result[A] {
	if ma == nil {
		return Error[A](onNil())
	}
	return Ok(*ma)
}

func FromResult1[A, I any](f func(context.Context, I) (A, error), i I) gofp.Result[A] {
	return func(ctx context.Context) (A, error) {
		return f(ctx, i)
	}
}

func FromResult[A any](result gofp.Result[A]) (r gofp.Result[A]) {
	r = result
	return
}

func fst[A, B any](a A) func(B) A {
	return func(B) A {
		return a
	}
}

func snd[A, B any](A) func(B) B {
	return func(b B) B {
		return b
	}
}

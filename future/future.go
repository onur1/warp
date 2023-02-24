// Package future implements the Future type and associated operations.
package future

import (
	"context"
	"time"

	"github.com/onur1/fp/event"
	"github.com/onur1/fp/result"
)

// A Future represents a collection of asynchronous computations with
// their associated results.
type Future[A any] event.Event[result.Result[A]]

// Succeed creates a future that succeeds with a value.
func Succeed[A any](a A) Future[A] {
	return Future[A](event.Map(event.Of(a), result.Succeed[A]))
}

// Fail creates a future that fails with an error.
func Fail[A any](err error) Future[A] {
	return Future[A](event.Map(event.Of(err), result.Fail[A]))
}

// Success creates a future which always succeeds with values received from
// a source event.
func Success[A any](ea event.Event[A]) Future[A] {
	return Future[A](event.Map(ea, result.Succeed[A]))
}

// Failure creates a future which always fails with errors received from
// a source event.
func Failure[A any](ea event.Event[error]) Future[A] {
	return Future[A](event.Map(ea, result.Fail[A]))
}

// After creates a future that succeeds after a timeout.
func After[A any](dur time.Duration, a A) Future[A] {
	return Success(event.After(dur, a))
}

// FailAfter creates a future that fails with an error after a timeout.
func FailAfter[A any](dur time.Duration, err error) Future[A] {
	return Failure[A](event.After(dur, err))
}

// Attempt converts a function which returns a (value, error) pair into a future
// that either succeeds with a value or fails with an error.
func Attempt[A any](ra result.Result[A], onPanic func(any) error) Future[A] {
	return func(ctx context.Context, sub chan<- result.Result[A]) {
		defer close(sub)

		var (
			a    A
			err  error
			done = ctx.Done()
		)

		defer func() {
			if r := recover(); r != nil {
				select {
				case <-done:
				default:
					select {
					case <-done:
					case sub <- result.Fail[A](onPanic(r)):
					}
				}
			}
		}()

		a, err = ra()
		if err != nil {
			select {
			case <-done:
			default:
				select {
				case <-done:
				case sub <- result.Fail[A](err):
				}
			}
			return
		}

		select {
		case <-done:
		default:
			select {
			case <-done:
			case sub <- result.Succeed(a):
			}
		}
	}
}

func Map[A, B any](fa Future[A], f func(A) B) Future[B] {
	return Future[B](
		event.Map(
			event.Event[result.Result[A]](fa),
			func(ra result.Result[A]) result.Result[B] {
				return result.Map(ra, f)
			},
		),
	)
}

func Ap[A, B any](fab Future[func(A) B], fa Future[A]) Future[B] {
	return Future[B](event.Ap(
		event.Map(
			event.Event[result.Result[func(A) B]](fab),
			func(gab result.Result[func(A) B]) func(result.Result[A]) result.Result[B] {
				return func(ga result.Result[A]) result.Result[B] {
					return result.Ap(gab, ga)
				}
			},
		),
		event.Event[result.Result[A]](fa),
	))
}

func Chain[A, B any](ma Future[A], f func(A) Future[B]) Future[B] {
	return Future[B](
		event.Chain(event.Event[result.Result[A]](ma), func(ra result.Result[A]) event.Event[result.Result[B]] {
			return event.Event[result.Result[B]](result.Fold(ra, Fail[B], f))
		}),
	)
}

// Package future implements the Future type.
package future

import (
	"context"
	"time"

	"github.com/onur1/data"
	"github.com/onur1/data/event"
	"github.com/onur1/data/result"
)

// Succeed creates a future that succeeds with a value.
func Succeed[A any](a A) data.Future[A] {
	return data.Future[A](event.Map(event.Of(a), result.Ok[A]))
}

// Fail creates a future that fails with an error.
func Fail[A any](err error) data.Future[A] {
	return data.Future[A](event.Map(event.Of(err), result.Error[A]))
}

// Success creates a future which always succeeds with values received from
// a source event.
func Success[A any](ea data.Event[A]) data.Future[A] {
	return data.Future[A](event.Map(ea, result.Ok[A]))
}

// Failure creates a future which always fails with errors received from
// a source event.
func Failure[A any](ea data.Event[error]) data.Future[A] {
	return data.Future[A](event.Map(ea, result.Error[A]))
}

// After creates a future that succeeds after a timeout.
func After[A any](dur time.Duration, a A) data.Future[A] {
	return Success(event.After(dur, a))
}

// FailAfter creates a future that fails with an error after a timeout.
func FailAfter[A any](dur time.Duration, err error) data.Future[A] {
	return Failure[A](event.After(dur, err))
}

// Attempt converts a function which returns a (value, error) pair into a future
// that either succeeds with a value or fails with an error.
func Attempt[A any](ra data.Result[A], onPanic func(any) error) data.Future[A] {
	return func(ctx context.Context, sub chan<- data.Result[A]) {
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
					case sub <- result.Error[A](onPanic(r)):
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
				case sub <- result.Error[A](err):
				}
			}
			return
		}

		select {
		case <-done:
		default:
			select {
			case <-done:
			case sub <- result.Ok(a):
			}
		}
	}
}

func Map[A, B any](fa data.Future[A], f func(A) B) data.Future[B] {
	return data.Future[B](
		event.Map(
			data.Event[data.Result[A]](fa),
			func(ra data.Result[A]) data.Result[B] {
				return result.Map(ra, f)
			},
		),
	)
}

func Ap[A, B any](fab data.Future[func(A) B], fa data.Future[A]) data.Future[B] {
	return data.Future[B](event.Ap(
		event.Map(
			data.Event[data.Result[func(A) B]](fab),
			func(gab data.Result[func(A) B]) func(data.Result[A]) data.Result[B] {
				return func(ga data.Result[A]) data.Result[B] {
					return result.Ap(gab, ga)
				}
			},
		),
		data.Event[data.Result[A]](fa),
	))
}

func Chain[A, B any](ma data.Future[A], f func(A) data.Future[B]) data.Future[B] {
	return data.Future[B](
		event.Chain(data.Event[data.Result[A]](ma), func(ra data.Result[A]) data.Event[data.Result[B]] {
			return data.Event[data.Result[B]](result.Fold(ra, Fail[B], f))
		}),
	)
}

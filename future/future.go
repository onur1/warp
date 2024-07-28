// Package future implements the Future type.
package future

import (
	"context"
	"time"

	"github.com/onur1/fpgo"
	"github.com/onur1/fpgo/event"
	"github.com/onur1/fpgo/result"
)

// Succeed creates a future that succeeds with a value.
func Succeed[A any](a A) fpgo.Future[A] {
	return fpgo.Future[A](event.Map(event.Of(a), result.Ok[A]))
}

// Fail creates a future that fails with an error.
func Fail[A any](err error) fpgo.Future[A] {
	return fpgo.Future[A](event.Map(event.Of(err), result.Error[A]))
}

// Success creates a future which always succeeds with values received from
// a source event.
func Success[A any](ea fpgo.Event[A]) fpgo.Future[A] {
	return fpgo.Future[A](event.Map(ea, result.Ok[A]))
}

// Failure creates a future which always fails with errors received from
// a source event.
func Failure[A any](ea fpgo.Event[error]) fpgo.Future[A] {
	return fpgo.Future[A](event.Map(ea, result.Error[A]))
}

// After creates a future that succeeds after a timeout.
func After[A any](dur time.Duration, a A) fpgo.Future[A] {
	return Success(event.After(dur, a))
}

// FailAfter creates a future that fails with an error after a timeout.
func FailAfter[A any](dur time.Duration, err error) fpgo.Future[A] {
	return Failure[A](event.After(dur, err))
}

// Attempt converts a function which returns a (value, error) pair into a future
// that either succeeds with a value or fails with an error.
func Attempt[A any](ra fpgo.Result[A], onPanic func(any) error) fpgo.Future[A] {
	return func(ctx context.Context, sub chan<- fpgo.Result[A]) {
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

		a, err = ra(ctx)
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

func AttemptForever[A any](ra fpgo.Result[A], onPanic func(any) error) fpgo.Future[A] {
	return func(ctx context.Context, sub chan<- fpgo.Result[A]) {
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

		for {
			a, err = ra(ctx)
			if err != nil {
				select {
				case <-done:
					return
				default:
					select {
					case <-done:
						return
					case sub <- result.Error[A](err):
					}
				}
			}

			select {
			case <-done:
				return
			default:
				select {
				case <-done:
					return
				case sub <- result.Ok(a):
				}
			}
		}
	}
}

func Map[A, B any](fa fpgo.Future[A], f func(A) B) fpgo.Future[B] {
	return fpgo.Future[B](
		event.Map(
			fpgo.Event[fpgo.Result[A]](fa),
			func(ra fpgo.Result[A]) fpgo.Result[B] {
				return result.Map(ra, f)
			},
		),
	)
}

func Ap[A, B any](fab fpgo.Future[func(A) B], fa fpgo.Future[A]) fpgo.Future[B] {
	return fpgo.Future[B](event.Ap(
		event.Map(
			fpgo.Event[fpgo.Result[func(A) B]](fab),
			func(gab fpgo.Result[func(A) B]) func(fpgo.Result[A]) fpgo.Result[B] {
				return func(ga fpgo.Result[A]) fpgo.Result[B] {
					return result.Ap(gab, ga)
				}
			},
		),
		fpgo.Event[fpgo.Result[A]](fa),
	))
}

func Alt[A any](x fpgo.Future[A], y fpgo.Future[A]) fpgo.Future[A] {
	return fpgo.Future[A](event.Alt(fpgo.Event[fpgo.Result[A]](x), fpgo.Event[fpgo.Result[A]](y)))
}

func Chain[A, B any](ma fpgo.Future[A], f func(A) fpgo.Future[B]) fpgo.Future[B] {
	return fpgo.Future[B](
		event.Chain(fpgo.Event[fpgo.Result[A]](ma), func(ra fpgo.Result[A]) fpgo.Event[fpgo.Result[B]] {
			return func(ctx context.Context, c chan<- fpgo.Result[B]) {
				result.Reduce(ctx, ra, Fail[B], f)(ctx, c)
			}
		}),
	)
}

func ChainEvent[A, B any](ma fpgo.Event[A], f func(A) fpgo.Future[B]) fpgo.Future[B] {
	return fpgo.Future[B](
		event.Chain(ma, func(ra A) fpgo.Event[fpgo.Result[B]] {
			return func(ctx context.Context, c chan<- fpgo.Result[B]) {
				result.Reduce(ctx, result.Ok(ra), Fail[B], f)(ctx, c)
			}
		}),
	)
}

func FromEvent[A any](fa fpgo.Event[A]) fpgo.Future[A] {
	return fpgo.Future[A](event.Map(fa, result.Ok[A]))
}

func FromResult[A any](fa fpgo.Result[A]) fpgo.Future[A] {
	return fpgo.Future[A](event.Of(fa))
}

func FromResults[A any](fas []fpgo.Result[A]) fpgo.Future[A] {
	return fpgo.Future[A](event.From(fas))
}

func From[A any](as []A) fpgo.Future[A] {
	ras := make([]fpgo.Result[A], len(as))
	for i, v := range as {
		ras[i] = result.Ok(v)
	}
	return fpgo.Future[A](event.From(ras))
}

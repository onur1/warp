// Package future implements the Future type.
package future

import (
	"context"
	"time"

	"github.com/onur1/gofp"
	"github.com/onur1/gofp/event"
	"github.com/onur1/gofp/result"
)

// Succeed creates a future that succeeds with a value.
func Succeed[A any](a A) gofp.Future[A] {
	return gofp.Future[A](event.Map(event.Of(a), result.Ok[A]))
}

// Fail creates a future that fails with an error.
func Fail[A any](err error) gofp.Future[A] {
	return gofp.Future[A](event.Map(event.Of(err), result.Error[A]))
}

// Success creates a future which always succeeds with values received from
// a source event.
func Success[A any](ea gofp.Event[A]) gofp.Future[A] {
	return gofp.Future[A](event.Map(ea, result.Ok[A]))
}

// Failure creates a future which always fails with errors received from
// a source event.
func Failure[A any](ea gofp.Event[error]) gofp.Future[A] {
	return gofp.Future[A](event.Map(ea, result.Error[A]))
}

// After creates a future that succeeds after a timeout.
func After[A any](dur time.Duration, a A) gofp.Future[A] {
	return Success(event.After(dur, a))
}

// FailAfter creates a future that fails with an error after a timeout.
func FailAfter[A any](dur time.Duration, err error) gofp.Future[A] {
	return Failure[A](event.After(dur, err))
}

// Attempt converts a function which returns a (value, error) pair into a future
// that either succeeds with a value or fails with an error.
func Attempt[A any](ra gofp.Result[A], onPanic func(any) error) gofp.Future[A] {
	return func(ctx context.Context, sub chan<- gofp.Result[A]) {
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

func AttemptForever[A any](ra gofp.Result[A], onPanic func(any) error) gofp.Future[A] {
	return func(ctx context.Context, sub chan<- gofp.Result[A]) {
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

func Map[A, B any](fa gofp.Future[A], f func(A) B) gofp.Future[B] {
	return gofp.Future[B](
		event.Map(
			gofp.Event[gofp.Result[A]](fa),
			func(ra gofp.Result[A]) gofp.Result[B] {
				return result.Map(ra, f)
			},
		),
	)
}

func Ap[A, B any](fab gofp.Future[func(A) B], fa gofp.Future[A]) gofp.Future[B] {
	return gofp.Future[B](event.Ap(
		event.Map(
			gofp.Event[gofp.Result[func(A) B]](fab),
			func(gab gofp.Result[func(A) B]) func(gofp.Result[A]) gofp.Result[B] {
				return func(ga gofp.Result[A]) gofp.Result[B] {
					return result.Ap(gab, ga)
				}
			},
		),
		gofp.Event[gofp.Result[A]](fa),
	))
}

func Alt[A any](x gofp.Future[A], y gofp.Future[A]) gofp.Future[A] {
	return gofp.Future[A](event.Alt(gofp.Event[gofp.Result[A]](x), gofp.Event[gofp.Result[A]](y)))
}

func Chain[A, B any](ma gofp.Future[A], f func(A) gofp.Future[B]) gofp.Future[B] {
	return gofp.Future[B](
		event.Chain(gofp.Event[gofp.Result[A]](ma), func(ra gofp.Result[A]) gofp.Event[gofp.Result[B]] {
			return func(ctx context.Context, c chan<- gofp.Result[B]) {
				result.Reduce(ctx, ra, Fail[B], f)(ctx, c)
			}
		}),
	)
}

func ChainEvent[A, B any](ma gofp.Event[A], f func(A) gofp.Future[B]) gofp.Future[B] {
	return gofp.Future[B](
		event.Chain(ma, func(ra A) gofp.Event[gofp.Result[B]] {
			return func(ctx context.Context, c chan<- gofp.Result[B]) {
				result.Reduce(ctx, result.Ok(ra), Fail[B], f)(ctx, c)
			}
		}),
	)
}

func FromEvent[A any](fa gofp.Event[A]) gofp.Future[A] {
	return gofp.Future[A](event.Map(fa, result.Ok[A]))
}

func FromResult[A any](fa gofp.Result[A]) gofp.Future[A] {
	return gofp.Future[A](event.Of(fa))
}

func FromResults[A any](fas []gofp.Result[A]) gofp.Future[A] {
	return gofp.Future[A](event.From(fas))
}

func From[A any](as []A) gofp.Future[A] {
	ras := make([]gofp.Result[A], len(as))
	for i, v := range as {
		ras[i] = result.Ok(v)
	}
	return gofp.Future[A](event.From(ras))
}

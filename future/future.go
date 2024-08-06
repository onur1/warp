// Package future implements the Future type.
package future

import (
	"context"
	"time"

	"github.com/onur1/warp"
	"github.com/onur1/warp/event"
	"github.com/onur1/warp/result"
)

// Succeed creates a future that succeeds with a value.
func Succeed[A any](a A) warp.Future[A] {
	return warp.Future[A](event.Map(event.Of(a), result.Ok[A]))
}

// Fail creates a future that fails with an error.
func Fail[A any](err error) warp.Future[A] {
	return warp.Future[A](event.Map(event.Of(err), result.Error[A]))
}

// Success creates a future which always succeeds with values received from
// a source event.
func Success[A any](ea warp.Event[A]) warp.Future[A] {
	return warp.Future[A](event.Map(ea, result.Ok[A]))
}

// Failure creates a future which always fails with errors received from
// a source event.
func Failure[A any](ea warp.Event[error]) warp.Future[A] {
	return warp.Future[A](event.Map(ea, result.Error[A]))
}

// After creates a future that succeeds after a timeout.
func After[A any](dur time.Duration, a A) warp.Future[A] {
	return Success(event.After(dur, a))
}

// FailAfter creates a future that fails with an error after a timeout.
func FailAfter[A any](dur time.Duration, err error) warp.Future[A] {
	return Failure[A](event.After(dur, err))
}

// Attempt converts a function which returns a (value, error) pair into a future
// that either succeeds with a value or fails with an error.
func Attempt[A any](ra warp.Result[A], onPanic func(any) error) warp.Future[A] {
	return func(ctx context.Context, sub chan<- warp.Result[A]) {
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

func AttemptForever[A any](ra warp.Result[A], onPanic func(any) error) warp.Future[A] {
	return func(ctx context.Context, sub chan<- warp.Result[A]) {
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

func Map[A, B any](fa warp.Future[A], f func(A) B) warp.Future[B] {
	return warp.Future[B](
		event.Map(
			warp.Event[warp.Result[A]](fa),
			func(ra warp.Result[A]) warp.Result[B] {
				return result.Map(ra, f)
			},
		),
	)
}

func Ap[A, B any](fab warp.Future[func(A) B], fa warp.Future[A]) warp.Future[B] {
	return warp.Future[B](event.Ap(
		event.Map(
			warp.Event[warp.Result[func(A) B]](fab),
			func(gab warp.Result[func(A) B]) func(warp.Result[A]) warp.Result[B] {
				return func(ga warp.Result[A]) warp.Result[B] {
					return result.Ap(gab, ga)
				}
			},
		),
		warp.Event[warp.Result[A]](fa),
	))
}

func Alt[A any](x warp.Future[A], y warp.Future[A]) warp.Future[A] {
	return warp.Future[A](event.Alt(warp.Event[warp.Result[A]](x), warp.Event[warp.Result[A]](y)))
}

func Chain[A, B any](ma warp.Future[A], f func(A) warp.Future[B]) warp.Future[B] {
	return warp.Future[B](
		event.Chain(warp.Event[warp.Result[A]](ma), func(ra warp.Result[A]) warp.Event[warp.Result[B]] {
			return func(ctx context.Context, c chan<- warp.Result[B]) {
				result.Reduce(ctx, ra, Fail[B], f)(ctx, c)
			}
		}),
	)
}

func ChainEvent[A, B any](ma warp.Event[A], f func(A) warp.Future[B]) warp.Future[B] {
	return warp.Future[B](
		event.Chain(ma, func(ra A) warp.Event[warp.Result[B]] {
			return func(ctx context.Context, c chan<- warp.Result[B]) {
				result.Reduce(ctx, result.Ok(ra), Fail[B], f)(ctx, c)
			}
		}),
	)
}

func FromEvent[A any](fa warp.Event[A]) warp.Future[A] {
	return warp.Future[A](event.Map(fa, result.Ok[A]))
}

func FromResult[A any](fa warp.Result[A]) warp.Future[A] {
	return warp.Future[A](event.Of(fa))
}

func FromResults[A any](fas []warp.Result[A]) warp.Future[A] {
	return warp.Future[A](event.From(fas))
}

func From[A any](as []A) warp.Future[A] {
	ras := make([]warp.Result[A], len(as))
	for i, v := range as {
		ras[i] = result.Ok(v)
	}
	return warp.Future[A](event.From(ras))
}

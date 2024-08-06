// Package event implements the Event type.
package event

import (
	"context"
	"time"

	"github.com/onur1/warp"
	"github.com/onur1/warp/nilable"
)

// Map creates an event by applying a function on each value received from a source
// event.
func Map[A, B any](fa warp.Event[A], f func(A) B) warp.Event[B] {
	return func(ctx context.Context, sub chan<- B) {
		defer close(sub)

		var (
			as = make(chan A)
			a  A
		)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		go fa(ctx, as)

		for a = range as {
			select {
			case <-done:
				return
			default:
				select {
				case <-done:
					return
				case sub <- f(a):
				}
			}
		}
	}
}

// Ap creates an event by applying the latest observed function from the first event on
// each value received from the second event.
func Ap[A, B any](fab warp.Event[func(A) B], fa warp.Event[A]) warp.Event[B] {
	return func(ctx context.Context, sub chan<- B) {
		defer close(sub)

		var (
			abs      = make(chan func(A) B)
			as       = make(chan A)
			abLatest func(A) B
			ab       func(A) B
			a        A
			ok       bool
		)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		go fab(ctx, abs)
		go fa(ctx, as)

		select {
		case <-done:
			return
		default:
			select {
			case <-done:
				return
			case abLatest, ok = <-abs:
				if !ok {
					return
				}
			}
		}

		for {
			select {
			case ab, ok = <-abs:
				if !ok {
					abs = nil
					if as == nil {
						return
					}
					continue
				}
				abLatest = ab
			case a, ok = <-as:
				if !ok {
					as = nil
					if abs == nil {
						return
					}
					continue
				}
				select {
				case <-done:
					return
				default:
					select {
					case <-done:
						return
					case sub <- abLatest(a):
					}
				}
			}
		}
	}
}

// Chain creates an event which composes two events in sequence, using the return value
// of the first event to determine the next one.
func Chain[A, B any](fa warp.Event[A], f func(A) warp.Event[B]) warp.Event[B] {
	return func(ctx context.Context, sub chan<- B) {
		defer close(sub)

		var (
			as = make(chan A)
			bs chan B
			a  A
			b  B
			fb warp.Event[B]
		)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		go fa(ctx, as)

		for a = range as {
			fb = f(a)
			bs = make(chan B)

			go fb(ctx, bs)

			for b = range bs {
				select {
				case <-done:
					return
				default:
					select {
					case <-done:
						return
					case sub <- b:
					}
				}
			}
		}
	}
}

// Reduce returns a value by applying a function on each value received from an event,
// in order, passing in the value and the return value from the calculation on the
// preceding element.
func Reduce[A, B any](ctx context.Context, fa warp.Event[A], b B, f func(B, A) B) B {
	as := make(chan A)

	go fa(ctx, as)

	r := b

	for a := range as {
		r = f(r, a)
	}

	return r
}

// ReduceRight applies a function against an accumulator and each observed value of
// the event (from right-to-left) to reduce it to a single value.
// Same as Reduce but applied from end to start.
func ReduceRight[A, B any](ctx context.Context, fa warp.Event[A], b B, f func(A, B) B) B {
	asc := make(chan A)

	go fa(ctx, asc)

	var as []A

	for a := range asc {
		as = append(as, a)
	}

	l := len(as)
	i := l - 1

	r := b

	for ; i >= 0; i-- {
		r = f(as[i], r)
	}

	return r
}

// SampleOn creates an event which samples the latest values from the first event at
// the times when the second event fires.
func SampleOn[A, B any](fa warp.Event[A], fab warp.Event[func(A) B]) warp.Event[B] {
	return func(ctx context.Context, sub chan<- B) {
		defer close(sub)

		var (
			abs     = make(chan func(A) B)
			as      = make(chan A)
			a       A
			aLatest A
			ab      func(A) B
			ok      bool
		)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		go fa(ctx, as)
		go fab(ctx, abs)

		aLatest, ok = <-as
		if !ok {
			return
		}

		for {
			select {
			case <-done:
				return
			default:
				select {
				case <-done:
					return
				case a, ok = <-as:
					if !ok {
						as = nil
						if abs == nil {
							return
						}
						continue
					}
					aLatest = a
				case ab, ok = <-abs:
					if !ok {
						abs = nil
						if as == nil {
							return
						}
						continue
					}
					select {
					case <-done:
						return
					default:
						select {
						case <-done:
							return
						case sub <- ab(aLatest):
						}
					}
				}
			}
		}
	}
}

func identity[A any](a A) A {
	return a
}

// SampleOn_ creates an event which samples the latest values from the first event at the
// times when the second event fires, ignoring the values produced by the second event.
func SampleOn_[A, B any](fa warp.Event[A], fb warp.Event[B]) warp.Event[A] {
	return SampleOn(fa, Map(fb, func(_ B) func(A) A {
		return identity[A]
	}))
}

// Alt creates an event which emits values simultaneously from two source events.
func Alt[A any](x warp.Event[A], y warp.Event[A]) warp.Event[A] {
	return func(ctx context.Context, sub chan<- A) {
		defer close(sub)

		var (
			xs = make(chan A)
			ys = make(chan A)
			a  A
			ok bool
		)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		go x(ctx, xs)
		go y(ctx, ys)

		for {
			select {
			case a, ok = <-xs:
				if !ok {
					xs = nil
					if ys == nil {
						return
					}
					break
				}
				select {
				case <-done:
					return
				default:
					select {
					case <-done:
						return
					case sub <- a:
					}
				}
			case a, ok = <-ys:
				if !ok {
					ys = nil
					if xs == nil {
						return
					}
					break
				}
				select {
				case <-done:
					return
				default:
					select {
					case <-done:
						return
					case sub <- a:
					}
				}
			}
		}
	}
}

// Filter creates an event which emits values from a source event when a predicate holds.
func Filter[A any](fa warp.Event[A], predicate warp.Predicate[A]) warp.Event[A] {
	return func(ctx context.Context, sub chan<- A) {
		defer close(sub)

		var (
			as = make(chan A)
			a  A
		)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		go fa(ctx, as)

		for a = range as {
			if predicate(a) {
				select {
				case <-done:
					return
				default:
					select {
					case <-done:
						return
					case sub <- a:
					}
				}
			}
		}
	}
}

func FilterMap[A, B any](fa warp.Event[A], f func(a A) warp.Nilable[B]) warp.Event[B] {
	return func(ctx context.Context, sub chan<- B) {
		defer close(sub)

		var (
			as = make(chan A)
			a  A
			nb warp.Nilable[B]
		)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		go fa(ctx, as)

		for a = range as {
			nb = f(a)
			if nilable.IsSome(nb) {
				select {
				case <-done:
					return
				default:
					select {
					case <-done:
						return
					case sub <- *nb:
					}
				}
			}
		}
	}
}

func plus1_[A any](_ A, n int) int {
	return n + 1
}

// Count creates an event that emits the number of times a source event is fired.
func Count[A any](fa warp.Event[A]) warp.Event[int] {
	return Fold(fa, 0, plus1_[A])
}

// CountAll returns the number of times a source event is fired in total.
func CountAll[A any](ctx context.Context, fa warp.Event[A]) int {
	return ReduceRight(ctx, fa, 0, plus1_[A])
}

// Take creates an event which emits the first n values observed from a source event.
func Take[A any](fa warp.Event[A], n int) warp.Event[A] {
	return func(ctx context.Context, sub chan<- A) {
		defer close(sub)

		var (
			as = make(chan A)
			a  A
			i  = 0
		)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		go fa(ctx, as)

		for a = range as {
			i++

			select {
			case <-done:
				return
			default:
				select {
				case <-done:
					return
				case sub <- a:
				}
			}

			if i == n {
				return
			}
		}
	}
}

// Until creates an event which emits values from an event until a predicate holds.
func Until[A any](fa warp.Event[A], predicate warp.Predicate[A]) warp.Event[A] {
	return func(ctx context.Context, sub chan<- A) {
		defer close(sub)

		var (
			as = make(chan A)
			a  A
		)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		go fa(ctx, as)

		for a = range as {
			select {
			case <-done:
				return
			default:
				if predicate(a) {
					return
				}
				select {
				case <-done:
					return
				case sub <- a:
				}
			}
		}
	}
}

// Once creates an event which emits values from an event for once and the last
// time when a predicate holds.
func Once[A any](fa warp.Event[A], predicate warp.Predicate[A]) warp.Event[A] {
	return func(ctx context.Context, sub chan<- A) {
		defer close(sub)

		var (
			as = make(chan A)
			a  A
		)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		go fa(ctx, as)

		for a = range as {
			select {
			case <-done:
				return
			default:
				if !predicate(a) {
					continue
				}
				select {
				case <-done:
				case sub <- a:
				}
				return
			}
		}
	}
}

// Fold creates an event which combines the values from a source event by applying
// a function starting with an initial value.
func Fold[A, B any](fa warp.Event[A], b B, f func(A, B) B) warp.Event[B] {
	return func(ctx context.Context, sub chan<- B) {
		defer close(sub)

		var (
			as     = make(chan A)
			a      A
			result = b
		)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		go fa(ctx, as)

		for a = range as {
			result = f(a, result)
			select {
			case <-done:
				return
			default:
				select {
				case <-done:
					return
				case sub <- result:
				}
			}
		}
	}
}

// Of creates an event which emits a single value.
func Of[A any](a A) warp.Event[A] {
	return func(ctx context.Context, sub chan<- A) {
		defer close(sub)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		select {
		case <-done:
			return
		default:
			select {
			case <-done:
			case sub <- a:
			}
		}
	}
}

// Interval creates an event which emits the current time periodically.
func Interval(dur time.Duration) warp.Event[time.Time] {
	return func(ctx context.Context, sub chan<- time.Time) {
		defer close(sub)

		var (
			ticker = time.NewTicker(dur)
			t      time.Time
			ok     bool
		)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		defer ticker.Stop()

	LOOP:
		for {
			select {
			case <-done:
				break LOOP
			default:
				select {
				case <-done:
					break LOOP
				case t, ok = <-ticker.C:
					if !ok {
						break LOOP
					}
					select {
					case <-done:
						break LOOP
					default:
						select {
						case <-done:
							break LOOP
						case sub <- t:
						}
					}
				}
			}
		}
	}
}

func FromIO[A any](io warp.IO[A]) warp.Event[A] {
	return Map(
		Take(Empty(), 1),
		func(_ struct{}) A {
			return io()
		},
	)
}

// From creates an event which emits multiple values sequentially from the supplied slice.
func From[A any](as []A) warp.Event[A] {
	return func(ctx context.Context, sub chan<- A) {
		defer close(sub)

		var (
			i = 0
			l = len(as)
		)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		for ; i < l; i++ {
			select {
			case <-done:
				return
			default:
				select {
				case <-done:
					return
				case sub <- as[i]:
				}
			}
		}
	}
}

var empty = struct{}{}

// Empty creates an event which emits an empty struct forever.
func Empty() warp.Event[struct{}] {
	return func(ctx context.Context, sub chan<- struct{}) {
		defer close(sub)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		for {
			select {
			case <-done:
				return
			default:
				select {
				case <-done:
					return
				case sub <- empty:
				}
			}
		}
	}
}

// After creates an event which emits a value after waiting for the specified duration.
func After[A any](dur time.Duration, a A) warp.Event[A] {
	return func(ctx context.Context, sub chan<- A) {
		defer close(sub)

		var (
			ticker = time.NewTicker(dur)
		)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		defer ticker.Stop()

	LOOP:
		for {
			select {
			case <-done:
				break LOOP
			default:
				select {
				case <-done:
					break LOOP
				case <-ticker.C:
					select {
					case <-done:
						break LOOP
					default:
						select {
						case <-done:
						case sub <- a:
						}
					}
					break LOOP
				}
			}
		}
	}
}

// MapNotNil creates an event which filters out any nil values by applying a function
// on each value received from some source event.
func MapNotNil[A, B any](fa warp.Event[A], f func(A) *B) warp.Event[B] {
	return func(ctx context.Context, sub chan<- B) {
		defer close(sub)

		var (
			as = make(chan A)
			a  A
			b  *B
		)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		go fa(ctx, as)

		for a = range as {
			select {
			case <-done:
				return
			default:
			}
			if b = f(a); b != nil {
				select {
				case <-done:
					return
				default:
					select {
					case <-done:
						return
					case sub <- *b:
					}
				}
			}
		}
	}
}

// A Last represents an event associated with its last value.
type Last[A any] struct {
	Now  A
	Last A
}

// WithLast creates an event which emits successive event values.
func WithLast[A any](fa warp.Event[A]) warp.Event[Last[A]] {
	return Fold(fa, Last[A]{}, func(a A, l Last[A]) Last[A] {
		return Last[A]{Now: a, Last: l.Now}
	})
}

// A Time represents an event value associated with some time.
type Time[A any] struct {
	Value A
	Time  time.Time
}

// WithTime creates an event which reports the current local time.
func WithTime[A any](fa warp.Event[A]) warp.Event[Time[A]] {
	return func(ctx context.Context, sub chan<- Time[A]) {
		defer close(sub)

		var (
			as = make(chan A)
			a  A
		)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		go fa(ctx, as)

		for a = range as {
			select {
			case <-done:
				return
			default:
				select {
				case <-done:
					return
				case sub <- Time[A]{Value: a, Time: time.Now()}:
				}
			}
		}
	}
}

func FromChannel[A any](source <-chan A) warp.Event[A] {
	return func(ctx context.Context, sub chan<- A) {
		defer close(sub)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		var a A

		for a = range source {
			select {
			case <-done:
				return
			default:
				select {
				case <-done:
					return
				case sub <- a:
				}
			}
		}
	}
}

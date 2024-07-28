package future

import (
	"context"
	"fmt"

	"github.com/onur1/gofp"
	"github.com/onur1/gofp/result"
	"github.com/onur1/ring"
)

type par[A any] struct {
	head int
	tail int
	b    *ring.Ring[A]
}

func newPar[A any](limit int) *par[A] {
	s := new(par[A])
	s.b = ring.NewRing[A](limit)
	s.head, s.tail = 0, 0
	return s
}

func (s *par[A]) Flush() []A {
	diff := s.head - s.tail
	v := make([]A, diff)
	for i := s.tail; i != s.head; i++ {
		v[diff-s.head+i] = s.b.Del(i)
	}
	return v
}

func (s *par[A]) Limit() int {
	return s.b.Size()
}

type indexed[A any] interface {
	Result() gofp.Result[A]
	Index() int
}

type indexedError[A any] struct {
	error
	index int
}

func (err indexedError[A]) Index() int {
	return err.index
}

func (err indexedError[A]) Error() string {
	return err.error.Error()
}

func newIndexedError[A any](index int, err error) indexed[A] {
	return indexedError[A]{error: err, index: index}
}

func (va indexedError[A]) Result() gofp.Result[A] {
	return result.Error[A](va.error)
}

type indexedValue[A any] struct {
	a     A
	index int
}

func (va indexedValue[A]) Index() int {
	return va.index
}

func (va indexedValue[A]) Result() gofp.Result[A] {
	return result.Ok(va.a)
}

func newIndexedValue[A any](index int, a A) indexed[A] {
	return indexedValue[A]{a: a, index: index}
}

// Parallel creates a Future which emits a stream of Results in parallel while keeping
// the sequence of output same as the order of input.
func Parallel[A any](fas gofp.Future[A], parallelism int) gofp.Future[A] {
	return func(ctx context.Context, sub chan<- gofp.Result[A]) {
		defer close(sub)

		var (
			cra = make(chan gofp.Result[A])
			fa  gofp.Result[A]
			ok  bool
		)

		var done <-chan struct{}

		if ctx != nil {
			done = ctx.Done()
		}

		go fas(ctx, cra)

		var (
			s       = newPar[gofp.Result[A]](parallelism)
			writes  = make(chan indexed[A])
			reads   = cra
			limit   = s.Limit()
			pending = 0
		)

		var res indexed[A]

	LOOP:
		for {
			select {
			case <-done:
				reads, done = nil, nil
				if s.head == s.tail {
					break LOOP
				}
			default:
				select {
				case <-done:
					reads, done = nil, nil
					if s.head == s.tail {
						break LOOP
					}
					break
				case fa, ok = <-reads:
					if !ok {
						if s.head == s.tail {
							break LOOP
						}
						reads = nil
						break
					}

					go func(fa gofp.Result[A], pos int) {
						a, err := fa(ctx)
						if err != nil {
							writes <- newIndexedError[A](pos, err)
						} else {
							writes <- newIndexedValue(pos, a)
						}
					}(fa, s.head)

					s.head += 1

					if (s.head - s.tail) < limit {
						break
					}

					reads = nil
				case res = <-writes:
					ra := res.Result()

					s.b.Put(res.Index(), ra)

					pending += 1

					for s.b.Get(s.tail) != nil {
						ra = s.b.Del(s.tail)

						s.tail += 1
						pending -= 1

						select {
						case <-done:
							break LOOP
						default:
							select {
							case <-done:
								break LOOP
							case sub <- ra:
							}
						}
					}

					if (s.head - s.tail) < limit {
						reads = cra
					}
				}
			}
		}

		parBuffer := s.Flush()
		if len(parBuffer) > 0 {
			fmt.Printf("parallel buffer: %+v", parBuffer)
		}

		close(writes)
	}
}

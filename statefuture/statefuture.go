// Package statefuture implements the StateFuture type.
package statefuture

import (
	"github.com/onur1/data"
	"github.com/onur1/data/future"
	"github.com/onur1/data/these"
)

func Succeed[S, A any](a A) data.StateFuture[S, A] {
	return func(s S) data.Future[data.These[A, S]] {
		return future.Succeed(these.Both(a, s))
	}
}

func Map[S, A, B any](fa data.StateFuture[S, A], f func(A) B) data.StateFuture[S, B] {
	return func(s S) data.Future[data.These[B, S]] {
		return future.Map(
			fa(s),
			func(tas data.These[A, S]) data.These[B, S] {
				a, s1 := tas()
				return these.Both(f(a), s1)
			},
		)
	}
}

func Ap[S, A, B any](fab data.StateFuture[S, func(A) B], fa data.StateFuture[S, A]) data.StateFuture[S, B] {
	return func(s S) data.Future[data.These[B, S]] {
		return future.Chain(fab(s), func(tabs data.These[func(A) B, S]) data.Future[data.These[B, S]] {
			ab, s1 := tabs()
			return future.Map(fa(s1), func(tas data.These[A, S]) data.These[B, S] {
				a, s2 := tas()
				return these.Both(ab(a), s2)
			})
		})
	}
}

func Chain[S, A, B any](ma data.StateFuture[S, A], f func(A) data.StateFuture[S, B]) data.StateFuture[S, B] {
	return func(s S) data.Future[data.These[B, S]] {
		return future.Chain(ma(s), func(tas data.These[A, S]) data.Future[data.These[B, S]] {
			a, s1 := tas()
			return f(a)(s1)
		})
	}
}

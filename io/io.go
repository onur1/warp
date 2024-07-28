// Package io implements the IO type.
package io

import "github.com/onur1/fpgo"

func Map[A, B any](fa fpgo.IO[A], f func(A) B) fpgo.IO[B] {
	return func() B {
		return f(fa())
	}
}

func Ap[A, B any](fab fpgo.IO[func(A) B], fa fpgo.IO[A]) fpgo.IO[B] {
	return func() B {
		return fab()(fa())
	}
}

func Chain[A, B any](ma fpgo.IO[A], f func(A) fpgo.IO[B]) fpgo.IO[B] {
	return func() B {
		return f(ma())()
	}
}

func ApFirst[A, B any](fa fpgo.IO[A], fb fpgo.IO[B]) fpgo.IO[A] {
	return Ap(Map(fa, func(a A) func(B) A {
		return func(_ B) A {
			return a
		}
	}), fb)
}

func ApSecond[A, B any](fa fpgo.IO[A], fb fpgo.IO[B]) fpgo.IO[B] {
	return Ap(Map(fa, func(_ A) func(B) B {
		return func(b B) B {
			return b
		}
	}), fb)
}

func ChainFirst[A, B any](ma fpgo.IO[A], f func(A) fpgo.IO[B]) fpgo.IO[A] {
	return Chain(ma, func(a A) fpgo.IO[A] {
		return Map(f(a), func(_ B) A {
			return a
		})
	})
}

func ChainRec[A, B any](init A, f func(A) fpgo.IO[func() (A, B, bool)]) fpgo.IO[B] {
	return func() B {
		var (
			a  A
			b  B
			ok bool
		)

		a, b, ok = f(init)()()

		for {
			if ok {
				break
			}
			a, b, ok = f(a)()()
		}

		return b
	}
}

func Of[A any](a A) fpgo.IO[A] {
	return func() A {
		return a
	}
}

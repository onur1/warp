// Package io implements the IO type.
package io

import (
	"github.com/onur1/data"
)

func Map[A, B any](fa data.IO[A], f func(A) B) data.IO[B] {
	return func() B {
		return f(fa())
	}
}

func Ap[A, B any](fab data.IO[func(A) B], fa data.IO[A]) data.IO[B] {
	return func() B {
		return fab()(fa())
	}
}

func Chain[A, B any](ma data.IO[A], f func(A) data.IO[B]) data.IO[B] {
	return func() B {
		return f(ma())()
	}
}

func ApFirst[A, B any](fa data.IO[A], fb data.IO[B]) data.IO[A] {
	return Ap(Map(fa, func(a A) func(B) A {
		return func(_ B) A {
			return a
		}
	}), fb)
}

func ApSecond[A, B any](fa data.IO[A], fb data.IO[B]) data.IO[B] {
	return Ap(Map(fa, func(_ A) func(B) B {
		return func(b B) B {
			return b
		}
	}), fb)
}

func ChainFirst[A, B any](ma data.IO[A], f func(A) data.IO[B]) data.IO[A] {
	return Chain(ma, func(a A) data.IO[A] {
		return Map(f(a), func(_ B) A {
			return a
		})
	})
}

func ChainRec[A, B any](init A, f func(A) data.IO[func() (A, B, bool)]) data.IO[B] {
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

func Of[A any](a A) data.IO[A] {
	return func() A {
		return a
	}
}

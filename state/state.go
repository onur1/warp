// Package state implements the State type.
package state

import "github.com/onur1/data"

// Of creates a state with the given return value.
func Of[S, A any](a A) data.State[S, A] {
	return func(s S) (A, S) {
		return a, s
	}
}

// Get returns a state.
func Get[S any]() data.State[S, S] {
	return func(s S) (S, S) {
		return s, s
	}
}

// Put resets a state with the given value.
func Put[S, A any](s S) data.State[S, A] {
	return func(_ S) (a A, _ S) {
		return a, s
	}
}

// Modify transforms a state by applying a function on it.
func Modify[S, A any](f func(S) S) data.State[S, A] {
	return func(s S) (a A, _ S) {
		return a, f(s)
	}
}

// Gets returns a value which depend on a state.
func Gets[S, A any](f func(S) A) data.State[S, A] {
	return func(s S) (A, S) {
		return f(s), s
	}
}

func Map[E, A, B any](fa data.State[E, A], f func(A) B) data.State[E, B] {
	return func(s1 E) (B, E) {
		a, s2 := fa(s1)
		return f(a), s2
	}
}

func Ap[E, A, B any](fab data.State[E, func(A) B], fa data.State[E, A]) data.State[E, B] {
	return func(s1 E) (B, E) {
		f, s2 := fab(s1)
		a, s3 := fa(s2)
		return f(a), s3
	}
}

func Chain[E, A, B any](ma data.State[E, A], f func(A) data.State[E, B]) data.State[E, B] {
	return func(s1 E) (B, E) {
		a, s2 := ma(s1)
		return f(a)(s2)
	}
}

// ApFirst creates a state by combining two states, keeping only the state
// of the first.
func ApFirst[E, A, B any](fa data.State[E, A], fb data.State[E, B]) data.State[E, A] {
	return Ap(Map(fa, func(a A) func(B) A {
		return func(_ B) A {
			return a
		}
	}), fb)
}

// ApSecond creates a state by combining two states, keeping only the state
// of the second.
func ApSecond[E, A, B any](fa data.State[E, A], fb data.State[E, B]) data.State[E, B] {
	return Ap(Map(fa, func(_ A) func(B) B {
		return func(b B) B {
			return b
		}
	}), fb)
}

// ChainFirst composes two states in sequence, using the return value of one state
// to determine the next state, keeping only the result of the first.
func ChainFirst[E, A, B any](ma data.State[E, A], f func(A) data.State[E, B]) data.State[E, A] {
	return Chain(ma, func(a A) data.State[E, A] {
		return Map(f(a), func(_ B) A {
			return a
		})
	})
}

func Evaluate[S, A any](ma data.State[S, A], s S) (a A) {
	a, _ = ma(s)
	return
}

func Execute[S, A any](ma data.State[S, A], s S) (s1 S) {
	_, s1 = ma(s)
	return
}

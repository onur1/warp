// Package these implements the These type.
package these

import (
	"github.com/onur1/data"
	"github.com/onur1/data/either"
)

func Left[E, A comparable](e E) data.These[E, A] {
	return data.These[E, A](either.Left[E, A](e))
}

func Right[E, A comparable](a A) data.These[E, A] {
	return data.These[E, A](either.Right[E](a))
}

func Both[E, A comparable](e E, a A) data.These[E, A] {
	return func() (E, A) {
		return e, a
	}
}

func Fold[E, A comparable, B any](
	fa data.These[E, A],
	onLeft func(E) B,
	onRight func(A) B,
	onBoth func(E, A) B,
) B {
	var (
		e E
		a A
	)
	e1, a1 := fa()
	if e1 != e && a1 != a {
		return onBoth(e1, a1)
	} else if a1 != a {
		return onRight(a1)
	}
	return onLeft(e1)
}

func Swap[E, A comparable](fa data.These[E, A]) data.These[A, E] {
	return Fold(fa, Right[A, E], Left[A, E], func(e E, a A) data.These[A, E] {
		return Both(a, e)
	})
}

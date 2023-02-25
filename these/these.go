// Package these implements the These type.
package these

import (
	"reflect"

	"github.com/onur1/data"
	"github.com/onur1/data/either"
)

func Left[E, A any](e E) data.These[E, A] {
	return data.These[E, A](either.Left[E, A](e))
}

func Right[E, A any](a A) data.These[E, A] {
	return data.These[E, A](either.Right[E](a))
}

func Both[E, A any](e E, a A) data.These[E, A] {
	return func() (E, A) {
		return e, a
	}
}

func Fold[E, A, B any](
	fa data.These[E, A],
	onLeft func(E) B,
	onRight func(A) B,
	onBoth func(E, A) B,
) B {
	e1, a1 := fa()
	if !reflect.ValueOf(e1).IsZero() {
		if !(reflect.ValueOf(a1).IsZero()) {
			return onBoth(e1, a1)
		} else {
			return onLeft(e1)
		}
	} else {
		return onRight(a1)
	}
}

func Swap[E, A any](fa data.These[E, A]) data.These[A, E] {
	return Fold(fa, Right[A, E], Left[A, E], func(e E, a A) data.These[A, E] {
		return Both(a, e)
	})
}

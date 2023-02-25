// Package either implements the Either type.
package either

import (
	"github.com/onur1/data"
)

func Right[E, A comparable](a A) data.Either[E, A] {
	return func() (e E, _ A) {
		return e, a
	}
}

func Left[E, A comparable](e E) data.Either[E, A] {
	return func() (_ E, a A) {
		return e, a
	}
}


// Package statefuture implements the StateFuture type.
package statefuture

import (
	"github.com/onur1/data"
	"github.com/onur1/data/future"
	"github.com/onur1/data/these"
)

// Succeed creates a statefuture with the given return value.
func Succeed[S, A comparable](a A) data.StateFuture[S, A] {
	return func(s S) data.Future[data.These[A, S]] {
		return future.Succeed(these.Both(a, s))
	}
}


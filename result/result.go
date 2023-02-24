// Package result implements the Result type and associated operations.
package result

// A Result represents a result of a computation that either yields a value
// of type A, or fails with an error.
type Result[A any] func() (A, error)

// Succeed creates a result which always returns a value.
func Succeed[A any](a A) Result[A] {
	return func() (A, error) {
		return a, nil
	}
}

// Fail creates a result which always fails with an error.
func Fail[A any](err error) Result[A] {
	return func() (A, error) {
		return *(new(A)), err
	}
}

// Map creates a result by applying a function on a succeeding result.
func Map[A, B any](fa Result[A], f func(A) B) Result[B] {
	a, err := fa()
	if err != nil {
		return Fail[B](err)
	}
	return Succeed(f(a))
}

// MapError creates a result by applying a function on a failing result.
func MapError[A any](fa Result[A], f func(error) error) Result[A] {
	_, err := fa()
	if err != nil {
		return Fail[A](f(err))
	}
	return fa
}

// Ap creates a result by applying a function contained in the first result
// on the value contained in the second result.
func Ap[A, B any](fab Result[func(A) B], fa Result[A]) Result[B] {
	var (
		err error
		ab  func(A) B
		a   A
	)

	ab, err = fab()
	if err != nil {
		return Fail[B](err)
	}

	a, err = fa()
	if err != nil {
		return Fail[B](err)
	}

	return Succeed(ab(a))
}

// Chain creates a result which combines two results in sequence, using the
// return value of one result to determine the next one.
func Chain[A, B any](ma Result[A], f func(A) Result[B]) Result[B] {
	a, err := ma()
	if err != nil {
		return Fail[B](err)
	}
	return f(a)
}

// Bimap creates a result by mapping a pair of functions over an error or a value
// contained in a result.
func Bimap[A, B any](fa Result[A], f func(error) error, g func(A) B) Result[B] {
	a, err := fa()
	if err != nil {
		return Fail[B](f(err))
	}
	return Succeed(g(a))
}

// ApFirst creates a result by combining two effectful computations, keeping
// only the result of the first.
func ApFirst[A, B any](fa Result[A], fb Result[B]) Result[A] {
	return Ap(Map(fa, func(a A) func(B) A {
		return func(_ B) A {
			return a
		}
	}), fb)
}

// ApFirst creates a result by combining two effectful computations, keeping
// only the result of the second.
func ApSecond[A, B any](fa Result[A], fb Result[B]) Result[B] {
	return Ap(Map(fa, func(_ A) func(B) B {
		return func(b B) B {
			return b
		}
	}), fb)
}

// Fold takes two functions and a result and returns a value by applying
// one of the supplied functions to the inner value.
func Fold[A, B any](ma Result[A], onError func(error) B, onSuccess func(A) B) B {
	a, err := ma()
	if err != nil {
		return onError(err)
	}
	return onSuccess(a)
}

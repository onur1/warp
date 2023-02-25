// Package nilable implements the Nilable type and associated operations.
package nilable

// A Nilable represents an optional value which is either some value or nil.
type Nilable[A any] *A

func IsNil[A any](a Nilable[A]) bool {
	return a == nil
}

func IsSome[A any](a Nilable[A]) bool {
	return a != nil
}

func Nil[A any]() Nilable[A] {
	return nil
}

func Some[A any](a A) Nilable[A] {
	return &a
}

func Map[A, B any](fa Nilable[A], f func(A) B) Nilable[B] {
	if fa == nil {
		return nil
	}
	return Some(f(*fa))
}

func Ap[A, B any](fab Nilable[func(A) B], fa Nilable[A]) Nilable[B] {
	if fab == nil {
		return nil
	}
	if fa == nil {
		return nil
	}
	return Some((*fab)(*fa))
}

func ApFirst[A, B any](fa Nilable[A], fb Nilable[B]) Nilable[A] {
	return Ap(Map(fa, func(a A) func(B) A {
		return func(_ B) A {
			return a
		}
	}), fb)
}

func ApSecond[A, B any](fa Nilable[A], fb Nilable[B]) Nilable[B] {
	return Ap(Map(fa, func(_ A) func(B) B {
		return func(b B) B {
			return b
		}
	}), fb)
}

func Chain[A, B any](ma Nilable[A], f func(A) Nilable[B]) Nilable[B] {
	if ma == nil {
		return nil
	}
	return f(*ma)
}

// FromNullable creates a nilable from a pointer to a value of type A.
func FromNullable[A any](ptr *A) data.Nilable[A] {
	if ptr != nil {
		return Some(*ptr)
	}
	return nil
}

// FromPredicate creates a nilable by testing a value against a predicate first.
func FromPredicate[A any](a A, predicate data.Predicate[A]) data.Nilable[A] {
	if predicate(a) {
		return Some(a)
	} else {
		return Nil[A]()
	}
}

// Attempt creates a nilable by running a function which returns a value,
// recovering with nil if a panic is thrown.
func Attempt[A any](f data.Lazy[A]) data.Nilable[A] {
	defer func() {
		recover()
	}()
	return Some(f())
}

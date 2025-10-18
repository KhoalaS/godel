package utils

import "fmt"

// Option[T] represents either Some(value) or None.
type Option[T any] struct {
	ok  bool
	val T
}

// Some constructs an Option[T] containing a value.
func Some[T any](v T) Option[T] {
	return Option[T]{ok: true, val: v}
}

// None constructs an empty Option[T].
func None[T any]() Option[T] {
	var zero T
	return Option[T]{ok: false, val: zero}
}

// FromPointer constructs an Option[T] from a *T.
// If p == nil -> None, otherwise Some(*p).
func FromPointer[T any](p *T) Option[T] {
	if p == nil {
		return None[T]()
	}
	return Some(*p)
}

// IsSome returns true if the Option contains a value.
func (o Option[T]) IsSome() bool { return o.ok }

// IsNone returns true if the Option is empty.
func (o Option[T]) IsNone() bool { return !o.ok }

// Unwrap returns the contained value if Some, otherwise panics.
func (o Option[T]) Unwrap() T {
	if !o.ok {
		panic("called Unwrap on None Option")
	}
	return o.val
}

// Expect returns the value if Some, otherwise panics with the provided message.
func (o Option[T]) Expect(msg string) T {
	if !o.ok {
		panic(msg)
	}
	return o.val
}

// UnwrapOr returns the contained value if Some, otherwise returns def.
func (o Option[T]) UnwrapOr(def T) T {
	if o.ok {
		return o.val
	}
	return def
}

// UnwrapOrElse returns the contained value if Some, otherwise calls f and returns its result.
func (o Option[T]) UnwrapOrElse(f func() T) T {
	if o.ok {
		return o.val
	}
	return f()
}

// Value returns (value, true) if Some, or (zero, false) if None.
// This matches common Go idiom.
func (o Option[T]) Value() (T, bool) {
	return o.val, o.ok
}

// String implements fmt.Stringer for debugging convenience.
func (o Option[T]) String() string {
	if o.ok {
		return fmt.Sprintf("Some(%v)", o.val)
	}
	return "None"
}

package utils

func FromAny[T any](value any) Option[T] {
	if val, ok := value.(T); ok {
		return Some(val)
	}
	return None[T]()
}

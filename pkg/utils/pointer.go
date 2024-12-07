package utils

// PointerValue returns the dereferenced value of the given pointer, or a default value if the pointer is nil.
func PointerValue[T Generic](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}
	return *value
}

// Pointer returns a pointer to the provided value of a generic type. It's useful for creating pointers to literals.
func Pointer[T Generic](value T) *T {
	return &value
}

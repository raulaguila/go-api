package utils

func PointerValue[T Generic](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}
	return *value
}

func Pointer[T Generic](value T) *T {
	return &value
}

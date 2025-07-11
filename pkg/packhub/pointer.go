package packhub

func MapValue(value *map[string]any, key string, defaultValue any) any {
	if value == nil {
		return defaultValue
	}

	if v, ok := (*value)[key]; ok {
		return v
	}

	return defaultValue
}

func PointerValue[T any](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}
	return *value
}

func Pointer[T any](value T) *T {
	return &value
}

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test for PointerValue function
func TestPointerValue(t *testing.T) {
	tests := []struct {
		name         string
		value        interface{}
		defaultValue interface{}
		expected     interface{}
	}{
		{name: "Non-nil pointer", value: Pointer(42), defaultValue: 10, expected: 42},
		{name: "Nil pointer", value: nil, defaultValue: 10, expected: 10},
		{name: "Non-nil string pointer", value: Pointer("test"), defaultValue: "default", expected: "test"},
		{name: "Nil string pointer", value: nil, defaultValue: "default", expected: "default"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.defaultValue.(type) {
			case int:
				assert.Equal(t, tc.expected, func() int {
					if tc.value == nil {
						return PointerValue[int](nil, tc.defaultValue.(int))
					} else {
						return PointerValue[int](tc.value.(*int), tc.defaultValue.(int))
					}
				}())
			case string:
				assert.Equal(t, tc.expected, func() string {
					if tc.value == nil {
						return PointerValue[string](nil, tc.defaultValue.(string))
					} else {
						return PointerValue[string](tc.value.(*string), tc.defaultValue.(string))
					}
				}())
			}
		})
	}
}

// Test for Pointer function
func TestPointer(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  interface{}
	}{
		{name: "Integer pointer", value: 42, want: Pointer(42)},
		{name: "String pointer", value: "test", want: Pointer("test")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.value.(type) {
			case int:
				got := Pointer(tc.value.(int))
				assert.Equal(t, tc.want, got)
				assert.Equal(t, tc.value, *got)
			case string:
				got := Pointer(tc.value.(string))
				assert.Equal(t, tc.want, got)
				assert.Equal(t, tc.value, *got)
			}
		})
	}
}

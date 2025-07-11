package packhub

import (
	"reflect"
	"testing"
)

func TestMapValue(t *testing.T) {
	tests := []struct {
		name         string
		input        *map[string]any
		key          string
		defaultValue any
		expected     any
	}{
		{
			name:         "NilMap",
			input:        nil,
			key:          "key",
			defaultValue: 123,
			expected:     123,
		},
		{
			name:         "NonNilPointer",
			input:        &map[string]any{"key": "test value"},
			key:          "key",
			defaultValue: 123,
			expected:     "test value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := MapValue(tt.input, tt.key, tt.defaultValue); result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPointerValue(t *testing.T) {
	tests := []struct {
		name         string
		input        *int
		defaultValue int
		expected     int
	}{
		{name: "NilPointer", input: nil, defaultValue: 10, expected: 10},
		{name: "NonNilPointer", input: func() *int { v := 5; return &v }(), defaultValue: 10, expected: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := PointerValue(tt.input, tt.defaultValue); result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPointer(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{name: "IntValue", input: 5, expected: func() *int { v := 5; return &v }()},
		//{name: "PointerToIntValue", input: func() *int { v := 10; return &v }(), expected: func() *int { v := 10; return &v }()},
		{name: "StructValue", input: struct{ A int }{A: 1}, expected: func() *struct{ A int } { v := struct{ A int }{A: 1}; return &v }()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := reflect.ValueOf(tt.expected).Interface()
			result := Pointer(tt.input)
			if !reflect.DeepEqual(*result, reflect.ValueOf(expected).Elem().Interface()) {
				t.Errorf("got %+v, want %+v", *result, expected)
			}
		})
	}
}

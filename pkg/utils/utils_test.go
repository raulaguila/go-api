package utils

import (
	"testing"
)

func TestPointerValue(t *testing.T) {
	tests := []struct {
		name         string
		value        *int
		defaultValue int
		expected     int
	}{
		{"nilValue", nil, 10, 10},
		{"nonNilValue", Pointer(20), 10, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PointerValue(tt.value, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestPointer(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected int
	}{
		{"pointerValue", 42, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Pointer(tt.value)
			if *result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, *result)
			}
		})
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected int
	}{
		{"sumEmpty", []int{}, 0},
		{"sumSingle", []int{5}, 5},
		{"sumMultiple", []int{1, 2, 3}, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sum(tt.numbers...)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestAvg(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected float64
	}{
		{"avgEmpty", []int{}, 0},
		{"avgSingle", []int{5}, 5},
		{"avgMultiple", []int{1, 2, 3}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Avg(tt.numbers...)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected int
	}{
		{"maxEmpty", []int{}, 0},
		{"maxSingle", []int{5}, 5},
		{"maxMultiple", []int{1, 2, 3}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Max(tt.numbers...)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected int
	}{
		{"minEmpty", []int{}, 0},
		{"minSingle", []int{5}, 5},
		{"minMultiple", []int{1, 2, 3}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Min(tt.numbers...)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

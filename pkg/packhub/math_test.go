package packhub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected int
	}{
		{"maxEmpty", []int{}, 0},
		{"maxSingle", []int{5}, 5},
		{"maxMultiple", []int{1, 2, 3, 5, 0, 4}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Max(tt.numbers...))
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
		{"minMultiple", []int{5, 7, 1, 2, 3}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Min(tt.numbers...))
		})
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected int
	}{
		{"empty", []int{}, 0},
		{"single", []int{5}, 5},
		{"positive", []int{1, 2, 3, 4, 5}, 15},
		{"negative", []int{-1, -2, -3, -4, -5}, -15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Sum(tt.numbers...))
		})
	}
}

func TestAvg(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []float64
		expected float64
	}{
		{"empty", []float64{}, 0},
		{"single", []float64{5}, 5},
		{"positive", []float64{1, 2, 3, 4, 5}, 3},
		{"negative", []float64{-1, -2, -3, -4, -5}, -3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Avg(tt.numbers...))
		})
	}
}

func TestHasNonZero(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected bool
	}{
		{"empty", []int{}, false},
		{"single", []int{5}, true},
		{"zero", []int{0}, false},
		{"positive", []int{1, 2, 3, 4, 5}, true},
		{"negative", []int{-1, -2, -3, -4, -5}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, HasNonZero(tt.numbers...))
		})
	}
}

func TestHasZero(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected bool
	}{
		{"empty", []int{}, false},
		{"single", []int{5}, false},
		{"zero", []int{0}, true},
		{"positive", []int{1, 2, 3, 4, 5}, false},
		{"negative", []int{-1, -2, -3, -4, -5}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, HasZero(tt.numbers...))
		})
	}
}

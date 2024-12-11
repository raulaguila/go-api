package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStruct[T Numbers] struct {
	name     string
	numbers  []T
	expected T
}

func TestMax(t *testing.T) {
	tests := []testStruct[int]{
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
	tests := []testStruct[int]{
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
	tests := []testStruct[int]{
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
	tests := []testStruct[float64]{
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
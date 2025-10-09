package packhub

import (
	"math"
	"testing"
)

func TestRandomFloat64Range(t *testing.T) {
	min, max := 1.5, 3.5
	for i := 0; i < 100; i++ {
		val := RandomFloat64(min, max)
		if val < min || val > max {
			t.Errorf("RandomFloat64(%v, %v) = %v; out of range", min, max, val)
		}
	}
}

func TestRandomFloat64DecimalPlaces(t *testing.T) {
	min, max := 0.0, 1.0
	places := 3
	for i := 0; i < 100; i++ {
		val := RandomFloat64(min, max, places)
		// Check decimal places
		diff := val * math.Pow(10, float64(places))
		if math.Abs(diff-math.Round(diff)) > 1e-9 {
			t.Errorf("RandomFloat64(%v, %v, %d) = %v; not rounded to %d decimal places", min, max, places, val, places)
		}
		if val < min || val > max {
			t.Errorf("RandomFloat64(%v, %v, %d) = %v; out of range", min, max, places, val)
		}
	}
}

func TestRandomFloat64MinEqualsMax(t *testing.T) {
	val := RandomFloat64(2.5, 2.5)
	if val != 2.5 {
		t.Errorf("RandomFloat64(2.5, 2.5) = %v; want 2.5", val)
	}
	val = RandomFloat64(2.5, 2.5, 2)
	if val != 2.5 {
		t.Errorf("RandomFloat64(2.5, 2.5, 2) = %v; want 2.5", val)
	}
}

func TestRandomFloat64NegativeRange(t *testing.T) {
	min, max := -5.0, -1.0
	for i := 0; i < 100; i++ {
		val := RandomFloat64(min, max)
		if val < min || val > max {
			t.Errorf("RandomFloat64(%v, %v) = %v; out of range", min, max, val)
		}
	}
}

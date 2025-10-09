package packhub

import (
	"math"
	"math/rand/v2"
)

func RandomFloat64(min, max float64, decimalPlaces ...int) float64 {
	value := min + rand.Float64()*(max-min)
	if len(decimalPlaces) == 0 {
		return value
	}
	return math.Round(value*math.Pow(10, float64(decimalPlaces[0]))) / math.Pow(10, float64(decimalPlaces[0]))
}

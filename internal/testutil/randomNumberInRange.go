package testutil

import (
	"math"
	"math/rand/v2"
)

func RandomIntInRange(min, max int) int {
	if min == max {
		return min
	}
	if min > max {
		tmp := max
		max = min
		min = tmp
	}
	return rand.IntN(max-min) + min
}

func RandomFloat64InRange(min, max float64) float64 {
	val := rand.Float64()*(max-min) + min
	return math.Round(val*100) / 100
}

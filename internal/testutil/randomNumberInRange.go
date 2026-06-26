package testutil

import "math/rand/v2"

func RandomNumberInRange(min, max int) int {
	return rand.IntN(max-min) + min
}

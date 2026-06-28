package testutil

import (
	"math/rand"
	"time"
)

func RandomTime() time.Time {
	randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000

	randomNow := time.Unix(randomTime, 0)

	return randomNow
}

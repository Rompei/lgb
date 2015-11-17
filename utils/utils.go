package utils

import (
	"math/rand"
	"time"
)

// CheckRate retuns boolean value from given rate
func CheckRate(rate int) bool {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100) < rate {
		return true
	}
	return false
}

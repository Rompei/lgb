package utils

import (
	"math/rand"
	"time"
)

// CheckRate retuns boolean value from given rate
func CheckRate(rate float64) bool {
	rand.Seed(time.Now().UnixNano())
	if float64(rand.Intn(100)) < rate {
		return true
	}
	return false
}

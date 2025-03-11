package redpackage

import (
	"math/rand"
	"time"
)

func DivideRedPackage(totalAmount float64, totalCount int) []float64 {
	result := make([]float64, totalCount)

	remainAmount := totalAmount

	remainCount := totalCount

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < totalCount-1; i++ {
		max := remainAmount / float64(remainCount) * 2

		amount := r.Float64()*(max-0.01) + 0.01

		result[i] = float64(int(amount*100)) / 100

		remainAmount -= amount
		remainCount--
	}
	result[totalCount-1] = remainAmount
	return result
}

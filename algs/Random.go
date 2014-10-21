package algs

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Return int uniformly in [a, b)
func RandInt(a, b int) int {
	return int(rand.Float64()*float64(b-a) + float64(a))
}

// Return real number uniformly in [0, 1)
func Random() float64 {
	return rand.Float64()
}

// Knuth shuffle algorithm
func Shuffle(a []interface{}) {
	var n int = len(a)
	for i := 0; i < n; i++ {
		r := i + rand.Intn(n-i)
		a[r], a[i] = a[i], a[r]
	}
}

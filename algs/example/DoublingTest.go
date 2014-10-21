package main

import (
	"algs"
	"fmt"
	"time"
)

func timeTrial(n int) time.Duration {
	var max int = 100000
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = algs.RandInt(0-max, max)
	}
	timer := algs.NewStopwatch()
	_ = algs.ThreeSum(a)
	return timer.ElapsedTime()
}

func main() {
	for n := 250; true; n += n {
		t := timeTrial(n)
		fmt.Println(n, t)
	}
}

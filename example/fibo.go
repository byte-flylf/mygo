// Determinism Counts! (Memoization)
package main

import (
	"fmt"
	"time"
)

func Fibonacci(n int64) int64 {
	var f int64

	if n <= 1 {
		f = 1
	} else {
		f = Fibonacci(n-1) + Fibonacci(n-2)
	}

	return f
}

var fibonacciLUT map[int64]int64

func MemoizedFibonacci(n int64) int64 {
	var f int64

	if stored, ok := fibonacciLUT[n]; ok {
		return stored
	}

	if n <= 1 {
		f = 1
	} else {
		f = MemoizedFibonacci(n-1) + MemoizedFibonacci(n-2)
	}

	fibonacciLUT[n] = f
	return f
}

func testFib(count int64, name string, f func(int64) int64) {
	begin := time.Now()

	var x int64
	for x = 0; x < count; x++ {
		f(x)
	}

	end := time.Now()

	msg := "%d iterations of %s took %v to complete.\n"
	fmt.Printf(msg, count, name, end.Sub(begin))
}

func init() {
	fibonacciLUT = make(map[int64]int64)
}

func main() {
	testFib(50, "Non-memoized", Fibonacci)
	testFib(50, "Memoized", MemoizedFibonacci)
}

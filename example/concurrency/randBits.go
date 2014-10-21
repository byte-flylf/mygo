package main

import (
	"fmt"
)

func RandBits() <-chan int {
	ch := make(chan int)
	go func() {
		for {
			select {
			case ch <- 0:
			case ch <- 1:
			}
		}
	}()
	return ch
}

func main() {
	ch := RandBits()
	for i := 10; i > 0; i-- {
		fmt.Print(<-ch)
	}
	fmt.Println()
}

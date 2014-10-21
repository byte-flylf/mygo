// par.go
package main

import (
	"fmt"
	"runtime"
	"strconv"
)

func main() {
	runtime.GOMAXPROCS(2)
	ch := make(chan int)
	n := 1000
	for i := 0; i < n; i++ {
		task(strconv.Itoa(i), ch, 100)
	}
	fmt.Printf("begin\n")
	for i := 0; i < n; i++ {
		<-ch
	}
}

func task(name string, ch chan int, max int) {
	go func() {
		i := 1
		for i <= max {
			fmt.Printf("%s %d\n", name, i)
			i++
		}
		ch <- 1
	}()
}

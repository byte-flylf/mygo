package main

import (
	"fmt"
)

func race() {
	wait := make(chan struct{})
	n := 0
	go func() {
		n++
		close(wait)
	}()
	n++
	<-wait
	fmt.Println(n)
}

func main() {
	race()
}

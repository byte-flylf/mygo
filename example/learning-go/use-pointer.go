package main

import (
	"fmt"
)

type X int

func (x X) m() {
	fmt.Println("in m")
}

func main() {
	var x X
	i := 3

	x.m()
}

// file: array-slice.go
//
// Chapter 1, Arrays, slice and maps, Page 14
package main

import (
	"fmt"
)

func main() {
	var arr [10]int
	arr[0] = 42
	arr[1] = 13
	fmt.Printf("The first element is %d\n", arr[0])

	{
		a := [...]int{1, 2, 3}
		fmt.Println(a)
	}

	{
		a := [2][2]int{[...]int{1, 2}, [...]int{3, 4}}
		fmt.Println(a)
	}

	{
		a := [2][2]int{{1, 2}, {3, 4}}
		fmt.Println(a)
	}

	{
		sl := make([]int, 10)
		fmt.Printf("sl, len=%d, cap=%d\n", len(sl), cap(sl))
	}

	{
		a := [...]int{1, 2, 3, 4, 5}
		sl := a[2:4]
		fmt.Printf("sl, len=%d, cap=%d\n", len(sl), cap(sl))
	}

	{
		var array [100]int
		slice := array[0:99]
		slice[98] = 'a'
		// slice[99] = 'a'
	}

	{
		var a = [...]int{0, 1, 2, 3, 4, 5, 6, 7}
		var s = make([]int, 6)
		n1 := copy(s, a[0:])
		fmt.Println(n1, s)
		n2 := copy(s, s[2:])
		fmt.Println(n2, s)
	}
}

package main

import "fmt"
import "unsafe"

type Data struct {
	Col1 byte
	Col2 int
	Col3 string
	Col4 int
}

func main() {
	var v Data

	fmt.Println(unsafe.Sizeof(v))

	fmt.Println("----")

	fmt.Println(unsafe.Alignof(v.Col1))
	fmt.Println(unsafe.Alignof(v.Col2))
	fmt.Println(unsafe.Alignof(v.Col3))
	fmt.Println(unsafe.Alignof(v.Col4))

	fmt.Println("----")

	fmt.Println(unsafe.Offsetof(v.Col1))
	fmt.Println(unsafe.Offsetof(v.Col2))
	fmt.Println(unsafe.Offsetof(v.Col3))
	fmt.Println(unsafe.Offsetof(v.Col4))

	fmt.Println("----")

	v.Col1 = 98
	v.Col2 = 77
	v.Col3 = "1234567890abcdef"
	v.Col4 = 23

	fmt.Println(unsafe.Sizeof(v))

	fmt.Println("----")

	x := unsafe.Pointer(&v)

	fmt.Println(*(*byte)(x))
	fmt.Println(*(*int)(unsafe.Pointer(uintptr(x) + unsafe.Offsetof(v.Col2))))
	fmt.Println(*(*string)(unsafe.Pointer(uintptr(x) + unsafe.Offsetof(v.Col3))))
	fmt.Println(*(*int)(unsafe.Pointer(uintptr(x) + unsafe.Offsetof(v.Col4))))
}

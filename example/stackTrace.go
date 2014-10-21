// 在 Go 中获取 stacktrace
// http://theo.im/blog/2014/07/21/Printing-stacktrace-in-Go/
package main

import (
	"fmt"
	"runtime"
)

func StackTrace(all bool) string {
	// Reserve 10K buffer at first
	buf := make([]byte, 10240)
	for {
		size := runtime.Stack(buf, all)
		// The size of the buffer may be not enough to hold the stacktrace,
		// so double the buffer size
		if size == len(buf) {
			buf = make([]byte, len(buf)<<1)
			continue
		}
		break
	}
	return string(buf)
}

func test() {
	a := 1
	b := 2
	fmt.Println(a + b)
	fmt.Println(StackTrace(true))
	return
}

func main() {
	test()
}

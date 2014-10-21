// Go中可以使用“+”合并字符串，但是这种合并方式效率非常低，每合并一次，就必须遍历一次字符串。
// Java中提供StringBuilder类来解决这个问题。Go中也有类似的机制，那就是Buffer
package main

import (
	"bytes"
	"fmt"
	"time"
)

const MAX = 100000

func main() {
	_, t := testBuf()
	fmt.Println("string buffer: ", t, "ns")

	_, t = testPlus()
	fmt.Println("string plus: ", t, "ns")
}

func testPlus() (s string, t int64) {
	start := time.Now().UnixNano()
	for i := 0; i < MAX; i++ {
		s += "a"
	}
	end := time.Now().UnixNano()
	return s, end - start
}

func testBuf() (s string, t int64) {
	var buf bytes.Buffer
	start := time.Now().UnixNano()
	for i := 0; i < MAX; i++ {
		buf.WriteString("a")
	}
	s = buf.String()
	end := time.Now().UnixNano()
	return s, end - start
}

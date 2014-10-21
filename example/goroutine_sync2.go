// golang中有2种方式同步程序，一种使用channel，另一种使用锁机制.
// 锁机制，更具体的是sync.WaitGroup，一种较为简单的同步方法集。
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
	}

	for i := 0; i < 100; i++ {
		go wg.Done()
	}
	fmt.Println("exit")
	wg.Wait()
}

func add(wg sync.WaitGroup) {
	wg.Add(1)
}

func done(wg sync.WaitGroup) {
	wg.Done()
}

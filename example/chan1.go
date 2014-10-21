// 测试:
// 当我们关闭一个带缓冲区的频道时，如果缓冲区中还有消息，
// 接收端是会继续接收完剩余消息呢？还是直接就丢弃剩余消息呢？

package main

import "fmt"

func main() {
	input := make(chan int, 10)
	wait  := make(chan int)

	for i := 0; i < 10; i ++ {
		input <- i
	}

	close(input)

	go func() {
		for {
			if i, ok := <- input; ok {
				fmt.Println(i)
			} else {
				break
			}
		}
		wait <- 1
	}()

	<-wait
}

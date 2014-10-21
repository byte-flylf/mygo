// 同步goroutine, 方法1：使用channel
//
// 代码思想:
// 新建NUM_OF_QUIT个goroutine，这些个goroutine里面新建1个chan bool，
// 通过这个channel来接受退出的信号，这些channel在新建的时候，已经发给了handle_exit。
// 在handle_exit这个goroutine里面，1方面监控由系统发过来的退出信号，然后再通知其他的goroutin优雅地退出；
// 另一方面通过slice收集其他goroutine发过来的channel。
// handle_exit通知其他的goroutine优雅退出后，再发信号给main进程主动退出。

package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

const NUM_OF_QUIT int = 10

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	done := make(chan bool)
	receive_channel := make(chan chan bool)
	finish := make(chan bool)

	for i := 0; i < NUM_OF_QUIT; i++ {
		go do_while_select(i, receive_channel, finish)
	}

	go handle_exit(done, receive_channel, finish)

	<-done
	os.Exit(0)

}
func handle_exit(done chan bool, receive_channel chan chan bool, finish chan bool) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	chan_slice := make([]chan bool, 0)
	for {
		select {
		case <-sigs:
			for _, v := range chan_slice {
				v <- true
			}
			for i := 0; i < len(chan_slice); i++ {
				<-finish
			}
			done <- true
			runtime.Goexit()
		case single_chan := <-receive_channel:
			log.Println("the single_chan is ", single_chan)
			chan_slice = append(chan_slice, single_chan)
		}
	}
}
func do_while_select(num int, rece chan chan bool, done chan bool) {
	quit := make(chan bool)
	rece <- quit
	for {
		select {
		case <-quit:
			done <- true
			runtime.Goexit()
		default:
			//简单输出
			log.Println("the ", num, "is running")
		}
	}
}

// 信号量控制程序退出
package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    signalChan := make(chan os.Signal, 1)
    exitChan := make(chan int)
    go func() {
        <-signalChan
        exitChan <- 1
    }()
    signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

    fmt.Println("waiting...")
    <-exitChan
    fmt.Println("succ to exit")
}

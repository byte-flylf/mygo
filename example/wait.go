// 测试waitgroup
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

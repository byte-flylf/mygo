package main

import (
    "fmt"
    "time"
     . "algs"
)

func timeTrial(n int) time.Duration {
    max := 1000000
    a := make([]int, 0)
    for i:= 0; i < n; i++ {
        a = append(a, int(RandInt(-max, max)))
    }
    timer := NewStopwatch()
    _ = ThreeSum(a)
    return timer.ElapsedTime()
}

func doubleRatio() {
    prev := timeTrial(125)
    for n := 250; ; n += n {
        time := timeTrial(n)
        fmt.Println(int64(time)/int64(prev))
        prev = time
    }
}

func main() {
    doubleRatio()
}

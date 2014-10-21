package main

// 权重轮询调度算法
import (
    "fmt"
    "time"
)

var slaveDns = map[int]map[string]interface {} {
    0: {"connectstring": "root@tcp(172.16.0.164:3306)/shiqu_tools?charset=utf8", "weight": 2},
    1: {"connectstring": "root@tcp(172.16.0.165:3306)/shiqu_tools?charset=utf8", "weight": 4},
    2: {"connectstring": "root@tcp(172.16.0.166:3306)/shiqu_tools?charset=utf8", "weight": 8},
}

var i int = -1
var cw int = 0
var gcd int = 2

func getDns() string {
    for {
        i = (i+1) % len(slaveDns)
        if i == 0 {
            cw = cw - gcd
            if cw <= 0 {
                cw = getMaxWeight()
                if cw == 0 {
                    return ""
                }
            }
        }
        if weight, _ := slaveDns[i]["weight"].(int); weight >= cw {
            return slaveDns[i]["connectstring"].(string)
        }
    }
}

func getMaxWeight() int {
    max := 0
    for _, v := range slaveDns {
        if weight, _ := v["weight"].(int); weight >= max {
            max = weight
        }
    }
    return max
}

func main() {
    note := map[string]int{}
    s_time := time.Now().Unix()
    for i := 0; i < 100; i++ {
        s := getDns()
        fmt.Println(s)
        if note[s] != 0 {
            note[s]++
        } else {
            note[s] = 1
        }
    }
    e_time := time.Now().Unix()
    fmt.Println("total time: ", e_time - s_time)
    fmt.Println("------------------------------")
    for k, v := range note {
        fmt.Println(k, " ", v)
    }
}

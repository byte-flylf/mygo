// 生成一千万个随机数, 用来测试
package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
    "time"
)

func init() {
    rand.Seed(time.Now().UTC().UnixNano())
}

const MAXN = 10000000

func main() {
    file, err := os.Create("/tmp/data.txt")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()

    w := bufio.NewWriter(file)
    for i := 0; i < MAXN; i++ {
        w.WriteString(fmt.Sprintf("%d\n", rand.Intn(MAXN)))
    }
    w.Flush()
}

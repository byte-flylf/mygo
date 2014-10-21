// 随机密码生成器
package main

import (
    "fmt"
    "flag"
    "time"
    "math/rand"
)

var pwdCount *int = flag.Int("n", 10, "How many passwd")
var pwdLength *int = flag.Int("l", 20, "How long a passwd")

var pwdCodes = [...]byte{
    'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
    'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
    'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
    'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
    '`', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-', '=',
    '~', '!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '_', '+',
    '[', ']', '\\',';', '\'', ',', '.', '/',
    '{', '}', '|', ':', '"',  '<', '>', '?', }

func main() {
    flag.Parse()

    timens := int64(time.Now().Nanosecond())
    rand.Seed(timens)
    for i := 0; i < *pwdCount; i++ {
        pwd := make([]byte, *pwdLength)
        for j := 0; j < *pwdLength; j++ {
            pwd[j] = pwdCodes[rand.Intn(*pwdCount)]
        }
        fmt.Println(string(pwd))
    }
}

package main

import (
    "bufio"
    "fmt"
    "net"
)

func main() {
    l, err := net.Listen("tcp", "127.0.0.1:8053")
    if err != nil {
        fmt.Printf("Failure to listen: %s\n", err.Error())
    }
    for {
        if c, err := l.Accept(); err == nil {
            go echo(c)
        }
    }
}

func echo(c net.Conn) {
    defer c.Close()
    line, err := bufio.NewReader(c).ReadString('\n')
    if err != nil {
        fmt.Printf("Failure to read: %s\n", err.Error())
        return
    }
    fmt.Printf("receive %d bytes\n", len(line))
    _, err = c.Write([]byte(line))
    if err != nil {
        fmt.Printf("Failure to write: %s\n", err.Error())
    }
}

// 测试goprotobuf的使用
package main

import (
	"./lm"
	proto "code.google.com/p/goprotobuf/proto"
	"fmt"
	"os"
)

func main() {
	msg := &lm.Helloworld{Id: proto.Int32(101), Str: proto.String("hello")}
	path := string("/tm/log.txt")
	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("failed: %s\n", err)
		return
	}

	defer f.Close()
	buffer, err := proto.Marshal(msg)
	f.Write(buffer)
}

//  unique ID generator
// http://blog.cloudflare.com/go-at-cloudflare
package main

import (
	"crypto/sha1"
	"fmt"
	"time"
)

func main() {
	id := make(chan string)
	go func() {
		h := sha1.New()
		c := []byte(time.Now().String())
		for {
			h.Write(c)
			id <- fmt.Sprintf("%x", h.Sum(nil))
		}
	}()

	for i := 0; i < 20; i++ {
		fmt.Println(<-id)
	}
}

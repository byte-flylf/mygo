// Collect data from multiple backends
// version 1
package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Backend interface {
	Query(q string) string
}

type MyBackend string

func (b MyBackend) Query(q string) string {
	time.Sleep(time.Duration(rand.Intn(100)))
	return fmt.Sprintf("%s/%s", b, q)
}

type Skillex struct{}

func (s Skillex) Query(q string) string {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return "wub-wub-wub"
}

func QueryAll(q string, backends ...Backend) []string {
	results := []string{}
	for _, backend := range backends {
		r := backend.Query(q)
		results = append(results, r)
	}
	return results
}

func main() {
	b1 := MyBackend("server-1")
	b2 := MyBackend("server-2")
	b3 := Skillex{}

	began := time.Now()
	results := QueryAll("dubstep", b1, b2, b3)
	fmt.Println(results)
	fmt.Println(time.Since(began))
}

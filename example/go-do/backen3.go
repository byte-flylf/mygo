// Collect data from multiple backends
// version 2: Let's fire queries concurrently.
// version 3: Replicate backends
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

type Replicas []Backend

func (r Replicas) Query(q string) string {
	c := make(chan string, len(r))
	for _, backend := range r {
		go func(b Backend) { c <- b.Query(q) }(backend)
	}
	return <-c
}

func QueryAll(q string, backends ...Backend) []string {
	c := make(chan string, len(backends))
	for _, backend := range backends {
		go func(b Backend) { c <- b.Query(q) }(backend)
	}

	results := []string{}
	for i := 0; i < cap(c); i++ {
		results = append(results, <-c)
	}
	return results
}

func main() {
	r1 := Replicas{MyBackend("foo1"), MyBackend("foo2"), MyBackend("foo3")}
	r2 := Replicas{MyBackend("bar1"), MyBackend("bar2")}
	r3 := Replicas{Skillex{}, Skillex{}, Skillex{}, Skillex{}, Skillex{}}

	began := time.Now()
	results := QueryAll("dubstep", r1, r2, r3)
	fmt.Println(results)
	fmt.Println(time.Since(began))
}

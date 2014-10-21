// Pipelined data processing
// We're subscribing to an event publisher, converting messages to enriched data models, and feeding them into a data store.
package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Msg string

// isten function simulates an infinite stream of messages, pushing them down an output channel.
func Listen(out chan Msg) {
	for {
		time.Sleep(time.Duration(rand.Intn(250)) * time.Millisecond)
		if rand.Intn(10) < 6 {
			out <- "foo"
		} else {
			out <- "bar"
		}
	}
}

// Enrich stage reads a single message from the input channel, processes it, and pushes the result down the output channel.
func Enrich(in, out chan Msg) {
	for {
		msg := <-in
		msg = "* " + msg + " *"
		out <- msg
	}
}

func Filter(in, out chan Msg) {
	for {
		msg := <-in
		if msg == "bar" {
			continue
		}
		out <- msg
	}
}

// Store stage simulates writing the message somewhere.
func Store(in chan Msg) {
	for {
		msg := <-in
		fmt.Println(msg)
	}
}

// Using channels to pass ownership of a message between stages makes the program naturally concurrent.
// It also cleanly separates the business logic from transport semantics: total separation of concerns.
// Note that because the channels are unbuffered, you get automatic backpressure,
// which (in my experience) is generally what you want.

func main() {
	// build the infrastructure
	toFilter := make(chan Msg)
	toEnricher := make(chan Msg)
	toStore := make(chan Msg)

	// launch the actors
	go Listen(toFilter)
	go Filter(toFilter, toEnricher)
	// concurrency
	go Enrich(toEnricher, toStore)
	go Enrich(toEnricher, toStore)
	go Enrich(toEnricher, toStore)
	go Store(toStore)

	time.Sleep(1 * time.Second)
}

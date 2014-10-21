// Pipelined data processing
// We're subscribing to an event publisher, converting messages to enriched data models, and feeding them into a data store.
// version 2: modify the Listener to start an HTTP server, to generate those Msgs.
package main

import (
	"fmt"
	"net/http"
)

type Msg struct {
	Data string
	Done chan bool
}

// isten function simulates an infinite stream of messages, pushing them down an output channel.
func Listen(out chan *Msg) {
	h := func(w http.ResponseWriter, r *http.Request) {
		msg := &Msg{
			Data: r.FormValue("data"),
			Done: make(chan bool),
		}
		out <- msg

		success := <-msg.Done // wait for done signal
		if !success {
			w.Write([]byte(fmt.Sprintf("aborted: %s", msg.Data)))
			return
		}
		w.Write([]byte(fmt.Sprintf("OK: %s", msg.Data)))
	}

	http.HandleFunc("/incoming", h)
	fmt.Println("listening on : 8080")
	http.ListenAndServe(":8080", nil) // block
}

// Enrich stage reads a single message from the input channel, processes it, and pushes the result down the output channel.
func Enrich(in, out chan *Msg) {
	for {
		msg := <-in
		msg.Data = "* " + msg.Data + " *"
		out <- msg
	}
}

func Filter(in, out chan *Msg) {
	for {
		msg := <-in
		if msg.Data == "bar" {
			msg.Done <- false
			continue
		}
		out <- msg
	}
}

// Store stage simulates writing the message somewhere.
func Store(in chan *Msg) {
	for {
		msg := <-in
		fmt.Println(msg)
		msg.Done <- true
	}
}

// Using channels to pass ownership of a message between stages makes the program naturally concurrent.
// It also cleanly separates the business logic from transport semantics: total separation of concerns.
// Note that because the channels are unbuffered, you get automatic backpressure,
// which (in my experience) is generally what you want.

func main() {
	// build the infrastructure
	toFilter := make(chan *Msg)
	toEnricher := make(chan *Msg)
	toStore := make(chan *Msg)

	// launch the actors
	go Listen(toFilter)
	go Filter(toFilter, toEnricher)
	// concurrency
	go Enrich(toEnricher, toStore)
	go Enrich(toEnricher, toStore)
	go Enrich(toEnricher, toStore)
	go Store(toStore)

	select {} // block forever without spinning the CPU
}

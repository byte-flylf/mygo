// https://gist.github.com/icub3d/6108906
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/golang/groupcache"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

// This instances address and pool.
var addr string
var pool *groupcache.HTTPPool

// This instances group.
var dns *groupcache.Group

func init() {
	flag.StringVar(&addr, "addr", "127.0.0.1:50000",
		"the addr:port this groupcache instance should run on.")
}

func main() {
	flag.Parse()

	start()
	console()
}

func start() {
	pool = groupcache.NewHTTPPool("http://" + addr)

	dns = groupcache.NewGroup("dns", 64<<20, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			// Lookup the adddress based on the hostname (key).
			addrs, err := net.LookupHost(key)
			if err != nil {
				return err
			}

			// We'll just store the first.
			dest.SetString(addrs[0])
			return nil
		}))

	// Run in the background so we can use the console.
	go http.ListenAndServe(addr, nil)
}

func console() {
	scanner := bufio.NewScanner(os.Stdin)
	quit := false

	for !quit {
		fmt.Print("gc> ")

		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		parts := strings.Split(line, " ")
		cmd := parts[0]
		args := parts[1:]

		switch cmd {
		case "peers":
			pool.Set(args...)
		case "stats":
			stats := dns.CacheStats(groupcache.MainCache)
			fmt.Println("Bytes:    ", stats.Bytes)
			fmt.Println("Items:    ", stats.Items)
			fmt.Println("Gets:     ", stats.Gets)
			fmt.Println("Hits:     ", stats.Hits)
			fmt.Println("Evictions:", stats.Evictions)
		case "get":
			var data []byte
			err := dns.Get(nil, args[0],
				groupcache.AllocatingByteSliceSink(&data))

			if err != nil {
				fmt.Println("get error:", err)
				continue
			}

			fmt.Print(args[0], ":")
			io.Copy(os.Stdout, bytes.NewReader(data))
			fmt.Println()
		case "quit":
			quit = true
		default:
			fmt.Println("unrecognized command:", cmd, args)
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println("reading stdin:", err)
	}
}

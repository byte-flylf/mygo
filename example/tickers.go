// https://gobyexample.com/tickers
package main

import "time"
import "fmt"

func main() {
	ticker := time.NewTicker(time.Microsecond * 500)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()
	time.Sleep(time.Microsecond * 1500)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

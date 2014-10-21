package main
 
import (
	"code.google.com/p/go-tour/wc"
	"strings"
)
 
func WordCount(s string) map[string]int {
	words := make(map[string]int)
	
	for _, word := range strings.Fields(s) {
		words[word] += 1
	}
 
	return words
}
 
func main() {
	wc.Test(WordCount)
}

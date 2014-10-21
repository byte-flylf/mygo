// Go by Example: Sorting
package main

import "sort"
import "fmt"

func main() {
    fruits := []string {"peach", "banana", "kiwi"}
    sort.Strings(fruits)
    fmt.Println(fruits)
}

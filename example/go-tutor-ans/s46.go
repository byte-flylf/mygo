package main
 
import "fmt"
 
// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	var f1, f2 = 0, 1
	f := func() int {
		f0 := f1
		f1, f2 = f2, f1+f2
		return f0
	}
	return f
}
 
func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}

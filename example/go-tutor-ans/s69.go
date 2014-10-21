package main
 
import (
	"code.google.com/p/go-tour/tree"
	"fmt"
)
 
// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	// Walk all the way left first
	if t.Left != nil {
		// Use go Walk if we don't care about sorting
		Walk(t.Left, ch)
	}
	
	// When we can't go left anymore, send value
	ch <- t.Value
	
	// Walk right
	if t.Right != nil {
		Walk(t.Right, ch)
	}
}
 
// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch := make(chan int, 10)
	ch2 := make(chan int, 10)
	ts1 := make([]int, 10)
	ts2 := make([]int, 10)
	
	go Walk(t1, ch)
	go Walk(t2, ch2)
	for i := 0; i < 10; i++ {
		ts1[i] = <-ch
		ts2[i] = <-ch2
	}
	
	for i := range ts1 {
		if ts1[i] != ts2[i] {
			return false
		}
	}
 
	return true
}
 
func main() {
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}

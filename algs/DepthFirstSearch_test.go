package algs

import (
	"github.com/bmizerany/assert"
	"sort"
	"testing"
)

func TestDepthFirstSearch(t *testing.T) {
	G := NewAdjGraphForFile("./algs4-data/tinyG.txt")
	var search Search
	var slice []int
	// test 1
	search = NewDepthFirstSearch(G, 0)
	assert.NotEqual(t, search.Count(), G.V())
	slice = make([]int, 0)
	for v := 0; v < G.V(); v++ {
		if search.Marked(v) {
			slice = append(slice, v)
		}
	}
	sort.Ints(slice)
	assert.Equal(t, slice, []int{0, 1, 2, 3, 4, 5, 6})
	// test 2
	search = NewDepthFirstSearch(G, 9)
	assert.NotEqual(t, search.Count(), G.V())
	slice = make([]int, 0)
	for v := 0; v < G.V(); v++ {
		if search.Marked(v) {
			slice = append(slice, v)
		}
	}
	sort.Ints(slice)
	assert.Equal(t, slice, []int{9, 10, 11, 12})
}

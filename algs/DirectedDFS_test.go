package algs

import (
	"github.com/bmizerany/assert"
	"sort"
	"testing"
)

func dfsAdj(dfs *DirectedDFS, dg *Digraph) []int {
	out := make([]int, 0)
	for v := 0; v < dg.V(); v++ {
		if dfs.Marked(v) {
			out = append(out, v)
		}
	}
	sort.Ints(out)
	return out
}

func TestDirectedDFS(t *testing.T) {
	dg := NewDigraphForFile("./algs4-data/tinyDG.txt")
	assert.Equal(t, dg.V(), 13)

	var dfs *DirectedDFS
	dfs = NewDirectedDFS(dg, 1)
	assert.Equal(t, dfsAdj(dfs, dg), []int{1})

	dfs = NewDirectedDFS(dg, 2)
	assert.Equal(t, dfsAdj(dfs, dg), []int{0, 1, 2, 3, 4, 5})

	dfs = NewDirectedDFSFromMultSrc(dg, 1, 2, 6)
	assert.Equal(t, dfsAdj(dfs, dg), []int{0, 1, 2, 3, 4, 5, 6, 8, 9, 10, 11, 12})
}

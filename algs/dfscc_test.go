package algs

import (
	"github.com/bmizerany/assert"
	"sort"
	"testing"
)

func components(g Graph, c CC) [][]int {
	uf := make([][]int, c.Count())
	for i := 0; i < len(uf); i++ {
		uf[i] = make([]int, 0)
	}

	for v := 0; v < g.V(); v++ {
		id := c.ID(v)
		uf[id] = append(uf[id], v)
	}

	for i := 0; i < len(uf); i++ {
		sort.Ints(uf[i])
	}
	return uf
}

func TestDfsCC(t *testing.T) {
	var g Graph
	var c CC

	g = NewAdjGraphForFile("./algs4-data/tinyG.txt")
	c = NewDfsCC(g)
	assert.Equal(t, c.Count(), 3)
	assert.Equal(t, components(g, c), [][]int{{0, 1, 2, 3, 4, 5, 6}, {7, 8}, {9, 10, 11, 12}})
}

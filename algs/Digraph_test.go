package algs

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestDigraph(t *testing.T) {
	dg := NewDigraphForFile("./algs4-data/tinyDG.txt")
	assert.Equal(t, dg.V(), 13)
	assert.Equal(t, dg.E(), 22)
	assert.Equal(t, dg.Adj(0), []int{5, 1})
	assert.Equal(t, dg.Adj(1), []int{})
	assert.Equal(t, dg.Adj(2), []int{0, 3})
	assert.Equal(t, dg.Adj(3), []int{5, 2})
	assert.Equal(t, dg.Adj(4), []int{3, 2})
	assert.Equal(t, dg.Adj(5), []int{4})
	assert.Equal(t, dg.Adj(6), []int{9, 4, 8, 0})
	assert.Equal(t, dg.Adj(7), []int{6, 9})
	assert.Equal(t, dg.Adj(8), []int{6})
	assert.Equal(t, dg.Adj(9), []int{11, 10})
	assert.Equal(t, dg.Adj(10), []int{12})
	assert.Equal(t, dg.Adj(11), []int{4, 12})
	assert.Equal(t, dg.Adj(12), []int{9})
}

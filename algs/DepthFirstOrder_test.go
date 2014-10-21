package algs

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestDepthFirstOrder(t *testing.T) {
	g := NewDigraphForFile("./algs4-data/tinyDAG.txt")
	order := NewDepthFirstOrder(g)

	assert.Equal(t, order.Pre(), []int{0, 5, 4, 1, 6, 9, 11, 12, 10, 2, 3, 7, 8})
	assert.Equal(t, order.Post(), []int{4, 5, 1, 12, 11, 10, 9, 6, 0, 3, 2, 7, 8})
	assert.Equal(t, order.ReversePost(), []int{8, 7, 2, 3, 0, 6, 9, 10, 11, 12, 1, 5, 4})
}

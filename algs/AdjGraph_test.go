package algs

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestAdjGraph(t *testing.T) {
	var G Graph
	G = NewAdjGraphForFile("./algs4-data/tinyCG.txt")
	assert.Equal(t, G.V(), 6)
	assert.Equal(t, G.E(), 8)
	assert.Equal(t, G.Adj(0), []int{2, 1, 5})
	assert.Equal(t, G.Adj(1), []int{0, 2})
	assert.Equal(t, G.Adj(2), []int{0, 1, 3, 4})
	assert.Equal(t, G.Adj(3), []int{5, 4, 2})
	assert.Equal(t, G.Adj(4), []int{3, 2})
	assert.Equal(t, G.Adj(5), []int{3, 0})
}

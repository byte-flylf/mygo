package algs

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestDirectedCycle(t *testing.T) {
	dg := NewDigraphForFile("./algs4-data/tinyDG.txt")
	assert.Equal(t, dg.V(), 13)

	var c *DirectedCycle
	c = NewDirectedCycle(dg)
	assert.Equal(t, c.HasCycle(), true)
	assert.Equal(t, c.Cycle(), []int{3, 5, 4, 3})

	dg = NewDigraphForFile("./algs4-data/tinyDAG.txt")
	c = NewDirectedCycle(dg)
	assert.Equal(t, c.HasCycle(), false)

}

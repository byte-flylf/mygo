package algs

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestCycle(t *testing.T) {
	var g Graph

	g = NewAdjGraphForFile("./algs4-data/tinyG.txt")
	c := NewCycle(g)
	assert.Equal(t, c.HasCycle(), true)
	assert.Equal(t, c.Cycle(), []int{3, 4, 5, 3})

	g = NewAdjGraphForFile("./algs4-data/mediumG.txt")
	c = NewCycle(g)
	assert.Equal(t, c.HasCycle(), true)
	assert.Equal(t, c.Cycle(), []int{15, 0, 225, 15})
}

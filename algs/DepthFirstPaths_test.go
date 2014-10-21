package algs

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestDepthFirstPaths(t *testing.T) {
	G := NewAdjGraphForFile("./algs4-data/tinyCG.txt")
	var path Pather
	path = NewDepthFirstPaths(G, 0)

	assert.Equal(t, path.HasPathTo(0), true)
	assert.Equal(t, path.PathTo(0), []int{0})

	assert.Equal(t, path.HasPathTo(1), true)
	assert.Equal(t, path.PathTo(1), []int{0, 2, 1})

	assert.Equal(t, path.HasPathTo(2), true)
	assert.Equal(t, path.PathTo(2), []int{0, 2})

	assert.Equal(t, path.HasPathTo(3), true)
	assert.Equal(t, path.PathTo(3), []int{0, 2, 3})

	assert.Equal(t, path.HasPathTo(4), true)
	assert.Equal(t, path.PathTo(4), []int{0, 2, 3, 4})

	assert.Equal(t, path.HasPathTo(5), true)
	assert.Equal(t, path.PathTo(5), []int{0, 2, 3, 5})
}

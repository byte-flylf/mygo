package algs

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestWeightedQuickUnionUF(t *testing.T) {
	uf := NewWeightedQuickUnionUF(10)
	uf.Union(4, 3)
	uf.Union(3, 8)
	uf.Union(6, 5)
	uf.Union(9, 4)
	uf.Union(2, 1)
	uf.Union(8, 9)
	uf.Union(5, 0)
	uf.Union(7, 2)
	uf.Union(6, 1)
	uf.Union(1, 0)
	uf.Union(6, 7)
	assert.Equal(t, uf.id, []int{6, 2, 6, 4, 4, 6, 6, 2, 4, 4})

	uf = NewWeightedQuickUnionUF(10)
	uf.Union(8, 6)
	uf.Union(6, 3)
	uf.Union(7, 2)
	uf.Union(5, 0)
	uf.Union(8, 1)
	uf.Union(0, 7)
	uf.Union(8, 4)
	uf.Union(4, 5)
	uf.Union(2, 9)
	assert.Equal(t, uf.id, []int{5, 8, 7, 8, 8, 8, 8, 5, 8, 8})

}

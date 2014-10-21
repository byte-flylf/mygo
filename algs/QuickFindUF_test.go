package algs

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestQuickFindUF(t *testing.T) {
	uf := NewQuickFindUF(10)
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
	assert.Equal(t, uf.id, []int{1, 1, 1, 8, 8, 1, 1, 1, 8, 8})

	uf = NewQuickFindUF(10)
	uf.Union(2, 4)
	uf.Union(3, 8)
	uf.Union(9, 8)
	uf.Union(5, 0)
	uf.Union(8, 1)
	uf.Union(7, 8)
	assert.Equal(t, uf.id, []int{0, 1, 4, 1, 4, 0, 6, 1, 1, 1})

	// test 3
	uf = NewQuickFindUF(10)
	uf.Union(8, 5)
	uf.Union(8, 6)
	uf.Union(3, 9)
	uf.Union(8, 4)
	uf.Union(4, 0)
	uf.Union(1, 2)
	assert.Equal(t, uf.id, []int{0, 2, 2, 9, 0, 0, 0, 7, 0, 9})
}

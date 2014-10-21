package algs

import (
   "testing"
   "github.com/bmizerany/assert"
)

func TestQuickUnionUF(t *testing.T) {
    uf := NewQuickUnionUF(10)
    uf.Union(4, 3)
    assert.Equal(t, uf.id, []int{0, 1, 2, 3, 3, 5, 6, 7, 8, 9})
    uf.Union(3, 8)
    assert.Equal(t, uf.id, []int{0, 1, 2, 8, 3, 5, 6, 7, 8, 9})
    uf.Union(6, 5)
    assert.Equal(t, uf.id, []int{0, 1, 2, 8, 3, 5, 5, 7, 8, 9})
    uf.Union(9, 4)
    assert.Equal(t, uf.id, []int{0, 1, 2, 8, 3, 5, 5, 7, 8, 8})
    uf.Union(2, 1)
    assert.Equal(t, uf.id, []int{0, 1, 1, 8, 3, 5, 5, 7, 8, 8})
    uf.Union(8, 9)
    uf.Union(5, 0)
    assert.Equal(t, uf.id, []int{0, 1, 1, 8, 3, 0, 5, 7, 8, 8})
    uf.Union(7, 2)
    assert.Equal(t, uf.id, []int{0, 1, 1, 8, 3, 0, 5, 1, 8, 8})
    uf.Union(6, 1)
    assert.Equal(t, uf.id, []int{1, 1, 1, 8, 3, 0, 5, 1, 8, 8})
    uf.Union(1, 0)
    uf.Union(6, 7)
    assert.Equal(t, uf.id, []int{1, 1, 1, 8, 3, 0, 5, 1, 8, 8})
    // test 2
    uf = NewQuickUnionUF(10)
    uf.Union(8, 6)
    uf.Union(6, 3)
    uf.Union(7, 2)
    uf.Union(5, 0)
    uf.Union(8, 1)
    uf.Union(0, 7)
    uf.Union(8, 4)
    uf.Union(4, 5)
    uf.Union(2, 9)
    assert.Equal(t, uf.id, []int{2, 4, 9, 1, 2, 0, 3, 2, 6, 9})
}

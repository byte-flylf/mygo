package algs

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestPercolation(t *testing.T) {
	p := NewPercolation(5)
	p.Open(1, 1)
	assert.Equal(t, p.IsOpen(1, 1), true)
	p.Open(3, 4)
	assert.Equal(t, p.IsOpen(3, 4), true)
	p.Open(4, 5)
	assert.Equal(t, p.IsOpen(5, 4), false)

	// test is full
	p = NewPercolation(5)
	p.Open(1, 1)
	p.Open(2, 1)
	p.Open(3, 1)
	p.Open(4, 1)
	p.Open(5, 1)
	assert.Equal(t, p.IsFull(1, 1), true)
	assert.Equal(t, p.IsFull(2, 1), true)
	assert.Equal(t, p.IsFull(3, 1), true)
	assert.Equal(t, p.IsFull(4, 1), true)
	assert.Equal(t, p.IsFull(5, 1), true)
	assert.Equal(t, p.IsFull(1, 2), false)
	assert.Equal(t, p.IsFull(2, 3), false)
	assert.Equal(t, p.IsFull(4, 5), false)
	assert.Equal(t, p.IsFull(1, 3), false)
	assert.Equal(t, p.IsFull(1, 5), false)

}

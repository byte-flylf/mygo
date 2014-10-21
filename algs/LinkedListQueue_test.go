package algs

import (
    "testing"
    "github.com/bmizerany/assert"
)

func TestLinkedListQueue(t *testing.T) {
    queue := NewLinkedListQueue()
    var n int = 8
    for i := 0; i < n; i++ {
        queue.Enqueue(i)
    }
    assert.Equal(t, queue.Size(), n)
    assert.Equal(t, queue.IsEmpty(), false)
    for i := 0; i < n; i++ {
        x := queue.Dequeue()
        j, _ := x.(int)
        assert.Equal(t, j, i)
    }
    assert.Equal(t, queue.IsEmpty(), true)
}

package algs

import (
    "testing"

    "github.com/bmizerany/assert"
)

func TestHeapSort(t *testing.T) {
    a := []string{"", "S", "O", "R", "T", "E", "X", "A", "M", "P", "L", "E"}
    HeapSort(a)
    assert.Equal(t, a, []string{"", "A", "E", "E", "L", "M", "O", "P", "R", "S", "T", "X"})
}

package algs

import (
   "testing"

    "github.com/bmizerany/assert"
)


func TestMergeSort(t *testing.T) {
    var a []string
    a = []string{"S", "O", "R", "T", "E", "X", "A", "M", "P", "L", "E"}
    MergeSort(a)
    assert.Equal(t, a, []string{"A", "E", "E", "L", "M", "O", "P", "R", "S", "T", "X"})

    a = []string{"S", "O", "R", "T", "E", "X", "A", "M", "P", "L", "E"}
    MergeXSort(a)
    assert.Equal(t, a, []string{"A", "E", "E", "L", "M", "O", "P", "R", "S", "T", "X"})

    a = []string{"S", "O", "R", "T", "E", "X", "A", "M", "P", "L", "E"}
    MergeBU(a)
    assert.Equal(t, a, []string{"A", "E", "E", "L", "M", "O", "P", "R", "S", "T", "X"})
}

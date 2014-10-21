package algs

import (
   "testing"
   "sort"

    "github.com/bmizerany/assert"
)


func TestSort(t *testing.T) {
    var a sort.StringSlice
    a = sort.StringSlice{"S", "O", "R", "T", "E", "X", "A", "M", "P", "L", "E"}
    SelectionSort(a)
    assert.Equal(t, a, sort.StringSlice{"A", "E", "E", "L", "M", "O", "P", "R", "S", "T", "X"})

    a = sort.StringSlice{"S", "O", "R", "T", "E", "X", "A", "M", "P", "L", "E"}
    InsertionSort(a)
    assert.Equal(t, a, sort.StringSlice{"A", "E", "E", "L", "M", "O", "P", "R", "S", "T", "X"})

    a = sort.StringSlice{"S", "O", "R", "T", "E", "X", "A", "M", "P", "L", "E"}
    ShellSort(a)
    assert.Equal(t, a, sort.StringSlice{"A", "E", "E", "L", "M", "O", "P", "R", "S", "T", "X"})

    a = sort.StringSlice{"S", "O", "R", "T", "E", "X", "A", "M", "P", "L", "E"}
    QuickSort(a)
    assert.Equal(t, a, sort.StringSlice{"A", "E", "E", "L", "M", "O", "P", "R", "S", "T", "X"})
}

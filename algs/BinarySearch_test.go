package algs

import (
    "testing"
    "github.com/bmizerany/assert"
)

func TestBinarySearch(t *testing.T) {
    arr := []int{10, 11, 12, 16, 18, 23, 29, 33, 48, 54, 57, 68, 77, 84, 98}
    assert.Equal(t, BinarySearch(23, arr), 5)
    assert.Equal(t, BinarySearch(50, arr), -1)
}

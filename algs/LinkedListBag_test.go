package algs

import (
   "testing"
    "github.com/bmizerany/assert"
)

func TestLinkedListBag(t *testing.T) {
    bag := NewLinkedListBag()
    bag.Add("a")
    bag.Add("b")
    bag.Add("c")
    bag.Add("d")
    assert.Equal(t, bag.Size(), 4)
    assert.Equal(t, bag.IsEmpty(), false)
    bag.Add("e")
    assert.Equal(t, bag.Size(), 5)
}
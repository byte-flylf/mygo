package algs

import (
    "testing"
    "github.com/bmizerany/assert"
)


func TestMaxpq(t *testing.T) {
    pq := NewMaxPQ(4)
    pq.Insert("P")
    pq.Insert("Q")
    pq.Insert("E")
    pq.DelMax()
    pq.Insert("X")
    pq.Insert("A")
    pq.Insert("M")
    pq.DelMax()
    pq.Insert("P")
    pq.Insert("L")
    pq.Insert("E")
    pq.DelMax()
    assert.Equal(t, pq.pq, []string{"", "P", "M", "L", "A", "E", "E", "", "", ""})
}

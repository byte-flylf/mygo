package algs

import (
	//"fmt"
	"github.com/bmizerany/assert"
	"strings"
	"testing"
)

func TestBinarySearchST(t *testing.T) {
	s := "S E A R C H E X A M P L E"
	parts := strings.Split(s, " ")
	st := NewBinarySearchST()
	for i, s := range parts {
		st.Put(s, i)
	}
	n := st.Size()
	assert.Equal(t, n, 10)
	assert.Equal(t, st.IsEmpty(), false)
	assert.Equal(t, st.Contains("F"), false)
	assert.Equal(t, st.Contains("M"), true)
	assert.Equal(t, st.Min(), "A")
	assert.Equal(t, st.Max(), "X")
	assert.Equal(t, st.Get("A"), 8)
	assert.Equal(t, st.Get("W"), nil)
	assert.Equal(t, st.Floor("G"), "E")
	assert.Equal(t, st.Ceiling("M"), "M")
	assert.Equal(t, st.Ceiling("O"), "P")
	assert.Equal(t, st.Select(3), "H")
	assert.Equal(t, st.Rank("H"), 3)

	// test Search
    var keys []string
	keys = []string{}
	for _, x := range st.RangeKeys("F", "T") {
		k, _ := x.(string)
		keys = append(keys, k)
	}
	assert.Equal(t, keys, []string{"H", "L", "M", "P", "R", "S"})
	assert.Equal(t, st.RangeSize("F", "T"), 6)
	// Keys
	keys = []string{}
	for _, x := range st.Keys() {
		k, _ := x.(string)
		keys = append(keys, k)
	}
	assert.Equal(t, keys, []string{"A", "C", "E", "H", "L", "M", "P", "R", "S", "X"})

	// test DeleteMin
	st.DeleteMin()
	keys = []string{}
	for _, node := range(st.Keys()) {
		k, _ := node.(string)
		keys = append(keys, k)
	}
	assert.Equal(t, keys, []string{"C", "E", "H", "L", "M", "P", "R", "S", "X"})
	// test DeleteMax
	st.DeleteMax()
	keys = []string{}
	for _, node := range(st.Keys()) {
		k, _ := node.(string)
		keys = append(keys, k)
	}
	assert.Equal(t, keys, []string{"C", "E", "H", "L", "M", "P", "R", "S"})
	// test Delete
	st2 := NewBinarySearchST()
	for i, s := range parts {
		st2.Put(s, i)
	}
	st2.Delete("E")
	keys = []string{}
	for _, node := range(st2.Keys()) {
		k, _ := node.(string)
		keys = append(keys, k)
	}
	assert.Equal(t, keys, []string{"A", "C",  "H", "L", "M", "P", "R", "S", "X"})
}

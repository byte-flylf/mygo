package algs

import (
	//"fmt"
	"github.com/bmizerany/assert"
	"strings"
	"testing"
)

func TestBST(t *testing.T) {
	s := "S E A R C H E X A M P L E"
	parts := strings.Split(s, " ")
	st := NewBST()
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

	// test PreOrder
	var keys []string
	var vals []int
	var nums []int
	for _, node := range st.PreOrder() {
		k, _ := node.key.(string)
		v, _ := node.val.(int)
		keys = append(keys, k)
		vals = append(vals, v)
		nums = append(nums, node.n)
	}
	assert.Equal(t, keys, []string{"S", "E", "A", "C", "R", "H", "M", "L", "P", "X"})
	assert.Equal(t, vals, []int{0, 12, 8, 4, 3, 5, 9, 11, 10, 7})
	assert.Equal(t, nums, []int{10, 8, 2, 1, 5, 4, 3, 1, 1, 1})
	// test InOrder
	keys, vals, nums = []string{}, []int{}, []int{}
	for _, node := range st.InOrder() {
		k, _ := node.key.(string)
		v, _ := node.val.(int)
		keys = append(keys, k)
		vals = append(vals, v)
		nums = append(nums, node.n)
	}
	assert.Equal(t, keys, []string{"A", "C", "E", "H", "L", "M", "P", "R", "S", "X"})
	assert.Equal(t, vals, []int{8, 4, 12, 5, 11, 9, 10, 3, 0, 7})
	assert.Equal(t, nums, []int{2, 1, 8, 4, 1, 3, 1, 5, 10, 1})
	// test PostOrder
	keys, vals, nums = []string{}, []int{}, []int{}
	for _, node := range st.PostOrder() {
		k, _ := node.key.(string)
		keys = append(keys, k)
	}
	assert.Equal(t, keys, []string{"C", "A", "L", "P", "M", "H", "R", "E", "X", "S"})
	// test Search
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
	keys, vals, nums = []string{}, []int{}, []int{}
	for _, node := range st.PreOrder() {
		k, _ := node.key.(string)
		v, _ := node.val.(int)
		keys = append(keys, k)
		vals = append(vals, v)
		nums = append(nums, node.n)
	}
	assert.Equal(t, keys, []string{"S", "E", "C", "R", "H", "M", "L", "P", "X"})
	assert.Equal(t, vals, []int{0, 12, 4, 3, 5, 9, 11, 10, 7})
	assert.Equal(t, nums, []int{9, 7, 1, 5, 4, 3, 1, 1, 1})
	// test DeleteMax
	st.DeleteMax()
	keys, vals, nums = []string{}, []int{}, []int{}
	for _, node := range st.PreOrder() {
		k, _ := node.key.(string)
		v, _ := node.val.(int)
		keys = append(keys, k)
		vals = append(vals, v)
		nums = append(nums, node.n)
	}
	assert.Equal(t, keys, []string{"S", "E", "C", "R", "H", "M", "L", "P"})
	assert.Equal(t, vals, []int{0, 12, 4, 3, 5, 9, 11, 10})
	assert.Equal(t, nums, []int{8, 7, 1, 5, 4, 3, 1, 1})
	// test Delete
	st2 := NewBST()
	for i, s := range parts {
		st2.Put(s, i)
	}
	st2.Delete("E")
	keys, vals, nums = []string{}, []int{}, []int{}
	for _, node := range st2.PreOrder() {
		k, _ := node.key.(string)
		v, _ := node.val.(int)
		keys = append(keys, k)
		vals = append(vals, v)
		nums = append(nums, node.n)
	}
	assert.Equal(t, keys, []string{"S", "H", "A", "C", "R", "M", "L", "P", "X"})
	assert.Equal(t, vals, []int{0, 5, 8, 4, 3, 9, 11, 10, 7})
	assert.Equal(t, nums, []int{9, 7, 2, 1, 4, 3, 1, 1, 1})
}

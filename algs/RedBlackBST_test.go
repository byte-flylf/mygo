package algs

import (
    "fmt"
	"github.com/bmizerany/assert"
	"strings"
	"testing"
)

func TestRedBlackBST(t *testing.T) {
	s := "S E A R C H E X A M P L E"
	parts := strings.Split(s, " ")
	st := NewRedBlackBST()
	for i, s := range parts {
		st.Put(s, i)
	}
	assert.Equal(t, st.Size(), 10)
	assert.Equal(t, st.IsEmpty(), false)
	assert.Equal(t, st.Contains("F"), false)
	assert.Equal(t, st.Contains("M"), true)
	assert.Equal(t, st.Min(), "A")
	assert.Equal(t, st.Max(), "X")
	assert.Equal(t, st.Get("A"), 8)
	assert.Equal(t, st.Get("W"), nil)

	// Keys
    keys := []string{}
    vals := []int{}
	for _, x := range st.Keys() {
		k, _ := x.(string)
		keys = append(keys, k)
        v, _ := st.Get(x).(int)
        vals = append(vals, v)
	}
	assert.Equal(t, keys, []string{"A", "C", "E", "H", "L", "M", "P", "R", "S", "X"})
	assert.Equal(t, vals, []int{8, 4, 12, 5, 11, 9, 10, 3, 0, 7})

    // select
    keys = []string{}
    for i := 0;  i < st.Size(); i++ {
        k, _ := st.Select(i).(string)
        keys = append(keys, k)
    }
	assert.Equal(t, keys, []string{"A", "C", "E", "H", "L", "M", "P", "R", "S", "X"})

    // rank, floor, ceiling
    ranks := []int{}
    floors := []string{}
    ceilings := []string{}
    for i := 'A'; i <= 'Z'; i++ {
        s := fmt.Sprintf("%c", i)
        ranks = append(ranks, st.Rank(s))
        x, _ := st.Floor(s).(string)
        floors = append(floors, x)
        y, _ := st.Ceiling(s).(string)
        ceilings = append(ceilings, y)
    }
	assert.Equal(t, ranks, []int{0, 1, 1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 6, 6, 6, 7, 7, 8, 9, 9, 9, 9, 9, 10, 10})
	assert.Equal(t, floors, []string{"A", "A", "C", "C", "E", "E", "E", "H", "H", "H", "H", "L", "M", "M", "M", "P", "P", "R",
                        "S", "S", "S", "S", "S", "X", "X", "X", })
    assert.Equal(t, ceilings, []string{"A", "C", "C", "E", "E", "H", "H", "H", "L", "L", "L", "L", "M", "P", "P", "P",
                    "R", "R", "S", "X", "X", "X", "X", "X", "", ""})

    //  test range search and range count
    from := []string{ "A", "Z", "X", "0", "B", "C" }
    to := []string { "Z", "A", "X", "Z", "G", "L" }
    vals = []int{}
    result := [][]string{}
    for i := 0; i < len(from); i++ {
        vals = append(vals, st.RangeSize(from[i], to[i]))
        row := make([]string, 0)
        for _, k := range(st.RangeKeys(from[i], to[i])) {
            s, _ := k.(string)
            row = append(row, s)
        }
        result = append(result, row)
    }
	assert.Equal(t, vals, []int{10, 0, 1, 10, 2, 4})
    assert.Equal(t, result, [][]string{{"A", "C", "E", "H", "L", "M", "P", "R", "S", "X"}, {}, {"X"},
                    {"A", "C", "E", "H", "L", "M", "P", "R", "S", "X"}, {"C", "E"}, {"C", "E", "H", "L"}})
    /*
    A-Z (10) : A C E H L M P R S X 
    Z-A ( 0) : 
    X-X ( 1) : X 
    0-Z (10) : A C E H L M P R S X 
    B-G ( 2) : C E 
    C-L ( 4) : C E H L 
    */

    //  delete the smallest keys
    for i := 0; i < st.Size() / 2; i++ {
        st.DeleteMin()
    }
    keys = []string{}
    vals = []int{}
	for _, x := range st.Keys() {
		k, _ := x.(string)
		keys = append(keys, k)
        v, _ := st.Get(x).(int)
        vals = append(vals, v)
	}
	assert.Equal(t, keys, []string{"H", "L", "M", "P", "R", "S", "X"})
	assert.Equal(t, vals, []int{5, 11, 9, 10, 3, 0, 7})

    // delete all the remaining keys
    for !st.IsEmpty() {
        st.Delete(st.Select(st.Size()/2))
    }

    // After adding back N keys
	for i, s := range parts {
		st.Put(s, i)
	}
    keys = []string{}
    vals = []int{}
	for _, x := range st.Keys() {
		k, _ := x.(string)
		keys = append(keys, k)
        v, _ := st.Get(x).(int)
        vals = append(vals, v)
	}
	assert.Equal(t, keys, []string{"A", "C", "E", "H", "L", "M", "P", "R", "S", "X"})
	assert.Equal(t, vals, []int{8, 4, 12, 5, 11, 9, 10, 3, 0, 7})

}

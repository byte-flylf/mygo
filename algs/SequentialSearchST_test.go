package algs

import (
	"github.com/bmizerany/assert"
	"strings"
	"testing"
)

func TestSequentialSearchST(t *testing.T) {
	s := "S E A R C H E X A M P L E"
	parts := strings.Split(s, " ")
	st := NewSequentialSearchST()
	for i, s := range parts {
		st.Put(s, i)
	}
	n := st.Size()
	assert.Equal(t, n, 10)
	assert.Equal(t, st.IsEmpty(), false)
	assert.Equal(t, st.Contains("F"), false)
	assert.Equal(t, st.Contains("M"), true)
	assert.Equal(t, st.Get("A"), 8)
	assert.Equal(t, st.Get("W"), nil)

    st.Delete("R")
    assert.Equal(t, st.Size(), 9)
    assert.Equal(t, st.Get("R"), nil)

    st.Put("A", 80)
    assert.Equal(t, st.Get("A"), 80)
}

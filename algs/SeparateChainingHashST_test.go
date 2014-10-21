package algs

import (
	"testing"
    "github.com/bmizerany/assert"
)

func TestSeparateChainingHashST(t *testing.T) {
    st := NewSeparateChainingHashST(2)
    st.Put("S", 6)
    st.Put("E", 10)
    st.Put("A", 4)
    st.Put("R", 14)
    st.Put("C", 5)
    st.Put("H", 4)
    st.Put("E", 10)
    st.Put("X", 15)
    st.Put("A", 4)
    assert.Equal(t, st.Size(), 7)
    assert.Equal(t, st.Get("A"), 4)
    st.Put("M", 1)
    st.Put("P", 14)
    st.Put("L", 6)
    st.Put("E", 10)
    assert.Equal(t, st.Get("E"), 10)
    assert.Equal(t, st.Size(), 10)
    st.Delete("A")
    assert.Equal(t, st.Size(), 9)
    assert.Equal(t, st.Contains("A"), false)
}

package algs

import (
    "testing"
    "github.com/bmizerany/assert"
)

func TestLinkedListStack(t *testing.T) {
    st := NewLinkedListStack()
    var n int = 8
    for i := 0; i < n; i++ {
        st.Push(i)
    }
    assert.Equal(t, st.Size(), n)
    assert.Equal(t, st.IsEmpty(), false)
    for i := 0; i < n; i++ {
        x := st.Pop()
        j, _ := x.(int)
        assert.Equal(t, j, n-1-i)
    }
    assert.Equal(t, st.IsEmpty(), true)
}

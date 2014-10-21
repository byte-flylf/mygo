// Symbol table implementation with sequential search in an
// unordered linked list of key-value pairs.
package algs

import (
    "github.com/sdming/kiss/gotype"
)

// a helper linked list data type
type STNode struct {
    key interface{}
    val interface{}
    next *STNode
}

type SequentialSearchST struct {
    n int  // number of key-value pairs
    first *STNode     // the linked list of key-value pairs
}

func NewSequentialSearchST() *SequentialSearchST {
    return new(SequentialSearchST)
}

// return number of key-value pairs
func (st *SequentialSearchST) Size() int {
    return st.n
}

// is the symbol table empty?
func (st *SequentialSearchST) IsEmpty() bool {
    return st.Size() == 0
}

// does this symbol table contain the given key?
func (st *SequentialSearchST) Contains(key string) bool {
    return st.Get(key) != nil
}

//  return the value associated with the key, or null if the key is not present
func (st *SequentialSearchST) Get(key interface{}) interface{} {
    for x := st.first; x != nil; x = x.next {
        if cmp := gotype.Compare(key, x.key); cmp == 0 {
            return x.val
        }
    }
    return nil
}

// add a key-value pair, replacing old key-value pair if key is already present
func (st *SequentialSearchST) Put(key interface{}, val interface{}) {
    for x := st.first; x != nil; x = x.next {
        if cmp := gotype.Compare(key, x.key); cmp == 0 {
            x.val = val
            return
        }
    }
    st.first = &STNode{key, val, st.first}
    st.n++
}

//  remove key-value pair with given key (if it's in the table)
func (st *SequentialSearchST) Delete(key interface{}) {
    st.first = st.delete(st.first, key)
    return
}

func (st *SequentialSearchST) delete(x *STNode, key interface{}) *STNode {
    if x == nil {
        return nil
    }
    if cmp := gotype.Compare(key, x.key); cmp == 0 {
        st.n--
        return x.next
    }
    x.next = st.delete(x.next, key)
    return x
}

// return all keys
func (st *SequentialSearchST) Keys() []interface{} {
    queue := make([]interface{}, 0)
    for p := st.first; p != nil; p = p.next {
        queue = append(queue, p.key)
    }
    return queue
}

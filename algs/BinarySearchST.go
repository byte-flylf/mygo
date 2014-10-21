// Symbol table implementation with binary search in an ordered array.
package algs

import (
    "errors"
    "github.com/sdming/kiss/gotype"
)

type BinarySearchST struct {
    keys []interface{}
    vals []interface{}
    n int
}

const INIT_CAPACITY = 2

func NewBinarySearchST() *BinarySearchST {
    st := BinarySearchST{}
    st.keys = make([]interface{}, INIT_CAPACITY)
    st.vals = make([]interface{}, INIT_CAPACITY)
    st.n = 0
    return &st
}

func (st *BinarySearchST) resize(capacity int) {
    if capacity <  len(st.keys) {
        return
    }
    tempk := make([]interface{}, capacity)
    tempv := make([]interface{}, capacity)
    copy(tempk, st.keys)
    copy(tempv, st.vals)
    st.keys, st.vals = tempk, tempv
}

func (st *BinarySearchST) Contains(key interface{}) bool {
    return st.Get(key) != nil
}

func (st *BinarySearchST) Size() int {
    return st.n
}

func (st *BinarySearchST) Get(key interface{}) interface {} {
    if st.IsEmpty() {
        return nil
    }
    i := st.Rank(key)
    if i < st.n && gotype.Compare(st.keys[i], key) == 0 {
        return st.vals[i]
    }
    return nil
}

func (st *BinarySearchST) IsEmpty() bool {
    return st.Size() == 0
}

// return the number of keys in the table that are smaller than given key
func (st *BinarySearchST) Rank(key interface{}) int {
    var lo, hi, m, cmp int
    lo, hi = 0, st.n-1
    for lo <= hi {  // <=, not <
        m = lo + (hi - lo)/2
        cmp = gotype.Compare(key, st.keys[m])
        if cmp < 0 {
            hi = m - 1
        } else if cmp > 0 {
            lo = m + 1
        } else {
            return m
        }
    }
    return lo
}

// Search for key. Update value if found; grow table if new. 
func (st *BinarySearchST) Put(key interface{}, val interface{}) {
    if val == nil {
        st.Delete(key)
        return
    }
    i := st.Rank(key)
    if i < st.n && gotype.Compare(key, st.keys[i]) == 0 {
        st.vals[i] = val
        return
    }
    // insert new key-value pair
    if st.n == len(st.keys) {
        st.resize(2 * st.n)
    }
    for j := st.n; j > i; j-- {
        st.keys[j] = st.keys[j-1]
        st.vals[j] = st.vals[j-1]
    }
    st.keys[i]= key
    st.vals[i] = val
    st.n++
}

// Remove the key-value pair if present
func (st *BinarySearchST) Delete(key interface{}) {
    if st.IsEmpty() {
        return
    }
    i := st.Rank(key)
    if i == st.n || gotype.Compare(key, st.keys[i]) != 0 {
        return
    }
    for j := i; j < st.n - 1; j++ {
        st.keys[j] = st.keys[j+1]
        st.vals[j] = st.vals[j+1]
    }
    st.n--
    st.keys[st.n], st.vals[st.n] = nil, nil

    // resize if 1/4 full
    if st.n > 0 && st.n == len(st.keys)/4 {
        st.resize(len(st.keys)/2)
    }
}

// delete the minimum key and its associated valuE
func (st *BinarySearchST) DeleteMin() error {
    if st.IsEmpty() {
        return errors.New("Symbol table underflow error")
    }
    st.Delete(st.Min())
    return nil
}

// delete the maximum key and its associated value
func (st *BinarySearchST) DeleteMax() error {
    if st.IsEmpty() {
        return errors.New("Symbol table underflow error")
    }
    st.Delete(st.Max())
    return nil
}

// Ordered symbol table methods
func (st *BinarySearchST) Min() interface{} {
    if st.IsEmpty() {
        return nil
    }
    return st.keys[0]
}

func (st *BinarySearchST) Max() interface{} {
    if st.IsEmpty() {
        return nil
    }
    return st.keys[st.n-1]
}

func (st *BinarySearchST) Select(k int) interface{} {
    if k < 0 || k > st.n {
        return nil
    }
    return st.keys[k]
}

func (st *BinarySearchST) Floor(key interface{}) interface{} {
    i := st.Rank(key)
    if i < st.n && gotype.Compare(key, st.keys[i]) == 0 {
        return st.keys[i]
    }
    if i == 0 {
        return nil
    }
    return st.keys[i-1]
}

func (st *BinarySearchST) Ceiling(key interface{}) interface{} {
    i := st.Rank(key)
    if i == st.n {
        return nil
    }
    return st.keys[i]
}

func (st *BinarySearchST) RangeSize(lo, hi interface{}) int {
    if gotype.Compare(lo, hi) > 0 {
        return 0
    }
    if st.Contains(hi) {
        return st.Rank(hi) - st.Rank(lo) + 1
    } else {
        return st.Rank(hi) - st.Rank(lo)
    }
}

func (st *BinarySearchST) RangeKeys(lo interface{}, hi interface{}) []interface{} {
    queue := make([]interface{}, 0, 100)
    if lo == nil || hi == nil {
        return nil
    }
    if gotype.Compare(lo, hi) > 0 {
        return nil
    }
    for i := st.Rank(lo); i < st.Rank(hi); i++ {
        queue = append(queue, st.keys[i])
    }
    if st.Contains(hi) {
        queue = append(queue, st.keys[st.Rank(hi)])
    }
    return queue
}

func (st *BinarySearchST) Keys() []interface{} {
    return st.RangeKeys(st.Min(), st.Max())
}

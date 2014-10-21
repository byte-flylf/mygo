package algs

import (
    "encoding/binary"
    "github.com/dgryski/dgohash"
)

// Symbol table implementation with linear probing hash table.
type LinearProbingHashST struct {
    n int
    m int
    keys []string
    vals []interface{}
}

// create linear proving hash table of given capacity
func NewLinearProbingHashST(capacity int) *LinearProbingHashST {
    st := new(LinearProbingHashST)
    st.m = capacity
    st.keys = make([]string, capacity)
    st.vals = make([]interface{}, capacity)
    return st
}

func (st *LinearProbingHashST) Size() int {
    return st.n
}

func (st *LinearProbingHashST) IsEmpty() bool {
    return st.Size() == 0
}

// does a key-value pair with the given key exist in the symbol table?
func (st *LinearProbingHashST) Contains(key string) bool {
    return st.Get(key) !=  nil
}

//  hash value between 0 and M-1
func (st *LinearProbingHashST) hash(key string) int {
    h := dgohash.NewJava32()
    h.Write([]byte(key))
    bsum := h.Sum(nil)
    s := binary.BigEndian.Uint32(bsum)
    res :=  int(s & uint32(st.m -1))
    //fmt.Println("key", key, "hash", res)
    return res
}

//  resize the hash table to the given capacity by re-hashing all of the keys
func (st *LinearProbingHashST) resize(capacity int) {
    var temp *LinearProbingHashST
    temp = NewLinearProbingHashST(capacity)
    for i := 0; i < st.m; i++ {
        if st.keys[i] != "" {
            temp.Put(st.keys[i], st.vals[i])
        }
    }
    st.keys = temp.keys
    st.vals = temp.vals
    st.m = capacity
}

// insert the key-value pair into the symbol table
func (st *LinearProbingHashST) Put(key string, val interface{}) {
    if (val == nil) {
        st.Delete(key)
    }
    // double table size if 50% full
    if st.n >= st.m / 2 {
        st.resize(2*st.m)
    }

    var i int
    for i = st.hash(key); st.keys[i] != ""; i = (i+1) % st.m {
        if st.keys[i] == key {
            st.vals[i] = val
            return
        }
    }
    st.keys[i] = key
    st.vals[i] = val
    st.n++
}

// return the value associated with the given key, null if no such value
func (st *LinearProbingHashST) Get(key string) interface{} {
    for i := st.hash(key); st.keys[i] != ""; i = (i+1) % st.m {
        if st.keys[i] == key {
            return st.vals[i]
        }
    }
    return nil
}

// delete the key (and associated value) from the symbol table
func (st *LinearProbingHashST) Delete(key string) {
    if !st.Contains(key) {
        return
    }
    i := st.hash(key)
    for key != st.keys[i] {
        i = (i+1) % st.m
    }
   st.keys[i] = ""
   st.vals[i] = nil
   i = (i+1) % st.m
   for st.keys[i] != "" {
       keyToRedo := st.keys[i]
       valToRedo := st.vals[i]
       st.keys[i], st.vals[i] = "", nil
       st.n--
       st.Put(keyToRedo, valToRedo)
       i = (i+1) % st.m
   }
   st.n--
   if st.n > 0 && st.n == st.m / 8 {
       st.resize(st.m / 2)
   }
}

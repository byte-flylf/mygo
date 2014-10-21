package algs

import (
    "encoding/binary"
    "github.com/dgryski/dgohash"
)


// A symbol table implemented with a separate-chaining hash table.
type SeparateChainingHashST struct {
	n  int                   // number of key-value pairs
	m  int                  //  hash table size
	st []*SequentialSearchST //  array of linked-list symbol tables
}

//  create separate chaining hash table with M lists
func NewSeparateChainingHashST(m int) *SeparateChainingHashST {
	st := make([]*SequentialSearchST, m)
	for i := 0; i < m; i++ {
		st[i] = NewSequentialSearchST()
	}
	return &SeparateChainingHashST{0, m, st}
}

// resize the hash table to have the given number of chains b rehashing all of the keys
func (self *SeparateChainingHashST) resize(chains int) {
    var temp *SeparateChainingHashST
    temp = NewSeparateChainingHashST(chains)

    for i := 0; i < self.m; i++ {
        for _, key := range(self.st[i].Keys()) {
            k := key.(string)
            v := self.st[i].Get(key).(int)
            temp.Put(k, v)
        }
    }
    self.m = chains
    self.st = temp.st
}

//  hash value between 0 and M-1
func (st *SeparateChainingHashST) hash(key string) uint32 {
    h := dgohash.NewJava32()
    h.Write([]byte(key))
    bsum := h.Sum(nil)
    s := binary.BigEndian.Uint32(bsum)
    res :=  s & uint32(st.m -1)
    return res
}

func (st *SeparateChainingHashST) Size() int {
    return st.n
}

func (st *SeparateChainingHashST) IsEmpty() bool {
    return st.Size() == 0
}

func (st *SeparateChainingHashST) Contains(key string) bool {
    return st.Get(key) != nil
}

func (self *SeparateChainingHashST) Get(key string) interface{} {
    i := self.hash(key)
    return self.st[i].Get(key)
}

func (st *SeparateChainingHashST) Put(key string,  val int) {
    // double table size if average length of list >= 10
    //if st.n >= st.m * 10 {
    if st.n >= st.m * 2 {
        st.resize(2*st.m)
    }

    i := st.hash(key)
    if (!st.st[i].Contains(key)) {
        st.n++
    }
    st.st[i].Put(key, val)
}

func (self *SeparateChainingHashST) Delete(key string) {
    i := self.hash(key)
    if self.st[i].Contains(key) {
        self.n--
    }
    self.st[i].Delete(key)

    if self.m > INIT_CAPACITY && self.n <= 2*self.m {
        self.resize(self.m / 2)
    }
}

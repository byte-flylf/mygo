package algs

import (
    "github.com/sdming/kiss/gotype"
)

// The MaxPQ class represents a priority queue of generic keys.
type MaxPQ struct {
    pq []string    //store items at indices 1 to N
    n int
}

// Initializes an empty priority queue with the given initial capacity
func NewMaxPQ(capacity int) *MaxPQ {
    p := MaxPQ{}
    p.pq = make([]string, capacity + 1)
    p.pq[0] = ""
    p.n = 0
    return &p
}

// Returns the number of keys on the priority queue.
func (self *MaxPQ) Size() int {
    return self.n
}

// Is the priority queue empty?
func (self *MaxPQ) IsEmpty() bool {
    return self.n == 0
}


// Returns a largest key on the priority queue.
func (self *MaxPQ) max() string {
    return self.pq[1]
}

// add a new key to the priority queue
func (self *MaxPQ) Insert(v string) {
    // double size of array if necessary
    if self.n >= len(self.pq) - 1 {
        self.resize(2*len(self.pq))
    }

    self.n++
    self.pq[self.n] = v 
    self.swim(self.n)
}

// Removes and returns a largest key on the priority queue
func (self *MaxPQ) DelMax() string {
    max := self.pq[1]
    // exch(1, N--)
    self.pq[1], self.pq[self.n] = self.pq[self.n], self.pq[1]
    self.n--

    self.sink(1)
    self.pq[self.n+1] = "" //  to avoid loiterig and help with garbage collection
    if self.n > 0 && self.n == (len(self.pq) -1)/4 {
        self.resize(len(self.pq)/2)
    }

    return max
}

// Helper functions to restore the heap invariant.
func (self *MaxPQ) swim(k int) {
    for k > 1 && gotype.Compare(self.pq[k/2], self.pq[k]) < 0 {
        self.pq[k/2], self.pq[k] = self.pq[k], self.pq[k/2]
        k = k/2
    }
}

func (self *MaxPQ) sink(k int) {
    for 2*k <= self.n {
        j := 2*k
        if j < self.n && gotype.Compare(self.pq[j], self.pq[j+1]) < 0 {
            j++
        }
        if gotype.Compare(self.pq[k], self.pq[j]) >= 0 {
            break
        }
        self.pq[k], self.pq[j] = self.pq[j], self.pq[k]
        k = j;
    }
}

// helper function to double the size of the heap array
func (self *MaxPQ) resize(capacity int) {
    //  assert capacity > self.n
    temp := make([]string, capacity)
    for i := 1; i <= self.n; i++ {
        temp[i] = self.pq[i]
    }
    self.pq = temp
}


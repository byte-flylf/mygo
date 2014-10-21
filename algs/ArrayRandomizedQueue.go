package algs

// array implementation of a randomized queue
// interface: randomized_queue.go
type ArrayRandomizedQueue struct {
    n int
    q []interface{}
}

func NewArrayRandomizedQueue() *ArrayRandomizedQueue {
    var queue *ArrayRandomizedQueue 
    queue = new(ArrayRandomizedQueue)
    queue.q = make([]interface{}, 0)
    return queue
}

// is the queue empty
func (self *ArrayRandomizedQueue) IsEmpty() bool {
    return self.n == 0
}

// return the number of items on the queue
func (self *ArrayRandomizedQueue) Size() int {
    return self.n
}

// add the item
func (self *ArrayRandomizedQueue) Enqueue(item interface{}) {
    self.q = append(self.q, item)
    self.n++
}

// delete and return a random item
func (self *ArrayRandomizedQueue) Dequeue() interface{} {
    if self.IsEmpty() {
        return nil
    }
    i := RandInt(0, self.n)
    self.q[i], self.q[0] = self.q[0], self.q[i]
    res := self.q[0]
    self.q = self.q[1:]
    return res
}

// return (but do not delete) a random item
func (self *ArrayRandomizedQueue) Sample() interface{} {
    if self.IsEmpty() {
        return nil
    }
    return self.q[RandInt(0, self.n)]
}


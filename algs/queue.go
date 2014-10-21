package algs

// A FIFO queue is a collection that is based on the first-in-first-out (FIFO) policy.
type Queue interface {
    Enqueue(item interface{})
    Dequeue() (item interface{})
    IsEmpty() bool
    Size() int
}

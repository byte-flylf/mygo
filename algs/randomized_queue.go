package algs

// Randomized queue
// A randomized queue is similar to a stack or queue, except that the item removed is chosen uniformly at random from items in the data structure.
type RandomizedQueue interface {
    IsEmpty() bool // is the queue empty?
    Size() int     // return the number of items on the queue
    Enqueue(item interface{})   // add the item
    Dequeue() interface{}    //  delete and return a random item
    Sample() interface{}     // return (but do not delete) a random item
}

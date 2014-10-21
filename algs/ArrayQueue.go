package algs

// Queue implementation with a resizing array.
type ArrayQueue struct {
	items []interface{} // queue elements
	n     int           // number of elements on queue
}

func NewArrayQueue() *ArrayQueue {
	return &ArrayQueue{items: []interface{}{}, n: 0}
}

func (self *ArrayQueue) Enqueue(item interface{}) {
	self.items = append(self.items, item)
	self.n++
}

func (self *ArrayQueue) Dequeue() interface{} {
	if self.IsEmpty() {
		return nil
	}
	item := self.items[0]
	self.items = self.items[1:]
	self.n--
	return item
}

func (self *ArrayQueue) IsEmpty() bool {
	return self.Size() == 0
}

func (self *ArrayQueue) Size() int {
	return self.n
}

package algs

// Doubly linked list implementation of a Deque.
// interface: deque.go
type LinkedListDeque struct {
    N int
    sentinel *DequeNode
}

type DequeNode struct {
    item interface{}
    prev *DequeNode
    next *DequeNode
}

// construct an empty deque
func NewLinkedListDeque() *LinkedListDeque {
    var s *DequeNode = new(DequeNode)
    s.item = nil
    s.prev, s.next = s, s
    var d *LinkedListDeque = new(LinkedListDeque)
    d.sentinel = s
    d.N = 0
    return d
}

// is the deque empty?
func (self *LinkedListDeque) IsEmpty() bool {
    return self.Size() == 0
}

// return the number of items on the deque
func (self *LinkedListDeque) Size() int {
    return self.N
}

// insert the item at the front
func (self *LinkedListDeque) AddFirst(item interface{}) {
    if item == nil {
        return
    }
    node := new(DequeNode)
    node.item = item
    node.prev = self.sentinel
    node.next = self.sentinel.next
    self.sentinel.next.prev = node
    self.sentinel.next = node
    self.N++
}

// insert the item at the end
func (self *LinkedListDeque) AddLast(item interface{}) {
    if item == nil {
        return
    }
    node := new(DequeNode)
    node.item = item
    node.prev = self.sentinel.prev
    node.next = self.sentinel
    self.sentinel.prev.next = node
    self.sentinel.prev = node
    self.N++
}

// delete and return item at the front
func (self *LinkedListDeque) RemoveFirst() interface{} {
    if self.IsEmpty() {
        return nil
    }
    res := self.sentinel.next.item
    self.sentinel.next = self.sentinel.next.next
    self.sentinel.next.prev = self.sentinel
    self.N--
    return res
}

// delete and return the item at the end
func (self *LinkedListDeque) RemoveLast() interface{} {
    if self.IsEmpty() {
        return nil
    }
    res := self.sentinel.prev.item
    // This is why we initialize sentinel.previous = sentinel;
    self.sentinel.prev = self.sentinel.prev.prev
    self.sentinel.prev.next = self.sentinel
    return res
}


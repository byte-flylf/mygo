package algs

type LinkedListQueue struct {
    first *Node
    last *Node
    n int
}

func NewLinkedListQueue() *LinkedListQueue {
    return &LinkedListQueue{first: nil, last: nil, n: 0}
}

func (self *LinkedListQueue) IsEmpty() bool {
    return self.first == nil
}

func (self *LinkedListQueue) Size() int {
    return self.n
}

func (self *LinkedListQueue) Enqueue(item interface{}) {
    oldlist := self.last
    self.last = &Node{item, nil}
    if self.IsEmpty() {
        self.first = self.last
    } else {
        oldlist.next = self.last
    }
    self.n++
}

func (self *LinkedListQueue) Dequeue() interface{} {
    if self.IsEmpty() {
        return nil
    }
    item := self.first.item
    self.first = self.first.next
    if self.first == nil {
        self.last = nil
    }
    return item
}

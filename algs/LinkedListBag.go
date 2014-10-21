package algs

type Node struct {
    item interface {}
    next *Node
}

type LinkedListBag struct {
    first *Node
    n int
}

func NewLinkedListBag() *LinkedListBag {
    return &LinkedListBag{nil, 0}
}

func (self *LinkedListBag) Size() int {
    return self.n
}

func (self *LinkedListBag) IsEmpty() bool {
    return self.first == nil
}

func (self *LinkedListBag) Add(v interface{}) {
    self.first = &Node{item:  v, next: self.first}
    self.n++
}

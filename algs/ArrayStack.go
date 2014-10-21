package algs

// Stack implementation with a resizing array.
type ArrayStack struct {
    items []interface{}
    n int
}

func NewArrayStack() *ArrayStack {
    return &ArrayStack{items: []interface{}{}, n: 0}
}

func (self *ArrayStack) Push(v interface{}) {
    self.items = append(self.items, v)
    self.n++
}

func (self *ArrayStack) Pop() interface{} {
    if self.IsEmpty() {
        return nil
    }
    last := self.Size() - 1
    item := self.items[last]
    self.items = self.items[:last]
    self.n--
    return item
}

func (self *ArrayStack) IsEmpty() bool {
    return self.Size() == 0
}

func (self *ArrayStack) Size() int {
    return len(self.items)
}

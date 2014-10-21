package algs


//  Bag implementation with a resizing array.
type ArrayBag struct {
	items []interface{} // array of items
	n     int           // number of elements on stack
}

func NewArrayBag() *ArrayBag {
	return &ArrayBag{items: []interface{}{}, n: 0}
}

func (self *ArrayBag) IsEmpty() bool {
	return self.Size() == 0
}

func (self *ArrayBag) Size() int {
	return self.n
}

func (self *ArrayBag) Add(v interface{}) {
	self.items = append(self.items, v)
	self.n++
}

func (self *ArrayBag) ChannelIterator() (chan interface{}) {
    ch := iterSlice(self.items)
    return ch
}

package algs

type LinkedListStack struct {
	first *Node
	n     int
}

func NewLinkedListStack() *LinkedListStack {
	return &LinkedListStack{}
}

func (me *LinkedListStack) Push(item interface{}) {
	oldfirst := me.first
	me.first = &Node{
		item: item,
		next: oldfirst,
	}
	me.n++
}

func (me *LinkedListStack) Pop() interface{} {
	if me.IsEmpty() {
		return nil
	}

	item := me.first.item
	me.first = me.first.next
	me.n--

	return item
}

func (me *LinkedListStack) IsEmpty() bool {
	return me.first == nil
}

func (me *LinkedListStack) Size() int {
	return me.n
}

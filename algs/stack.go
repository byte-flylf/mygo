package algs

//  A pushdown stack is a collection that is based on the last-in-first-out (LIFO) policy
type Stack interface {
	Push(item interface{})
	Pop() interface{}
	IsEmpty() bool
	Size() int
}

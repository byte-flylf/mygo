package algs

import (
	"errors"
	"fmt"
	"github.com/sdming/kiss/gotype"
)

// binary search tree node
type BstNode struct {
	key   interface{} // sorted by key
	val   interface{} // associated data
	n     int         // number of nodes in subtree
	left  *BstNode    // left and right subtrees
	right *BstNode
}

func NewBstNode(key interface{}, val interface{}) *BstNode {
	return &BstNode{key, val, 1, nil, nil}
}

func (n *BstNode) String() {
	fmt.Sprintf("(%s, %s, %d, %v, %v)", n.key, n.val, n.n, n.left, n.right)
}

// binaray search tree
// A symbol table implemented with a binary search tree.
// Data files:  http://algs4.cs.princeton.edu/32bst/tinyST.txt
type BST struct {
	root *BstNode // root of BST
}

func NewBST() *BST {
	return &BST{nil}
}

// is the symbol table empty
func (st *BST) IsEmpty() bool {
	return st.Size() == 0
}

// return number of key-value pairs in BST
func (st *BST) Size() int {
	return st.size(st.root)
}

// return number of key-value pairs in BST rooted at x
func (st *BST) size(x *BstNode) int {
	if x == nil {
		return 0
	}
	return x.n
}

// does there exist a key-value pair with given key
func (st *BST) Contains(key interface{}) bool {
	return st.Get(key) != nil
}

// return value associated with the given key, or null if no such key exists
func (st *BST) Get(key interface{}) interface{} {
	return st.get(st.root, key)
}

func (st *BST) get(x *BstNode, key interface{}) interface{} {
	if x == nil {
		return nil
	}
	cmp := gotype.Compare(key, x.key)
	if cmp < 0 {
		return st.get(x.left, key)
	} else if cmp > 0 {
		return st.get(x.right, key)
	}
	return x.val
}

func (st *BST) Put(key, val interface{}) {
	if val == nil {
		st.Delete(key)
		return
	}
	st.root = st.put(st.root, key, val)
}

func (st *BST) put(x *BstNode, key interface{}, val interface{}) *BstNode {
	if x == nil {
		return NewBstNode(key, val)
	}
	cmp := gotype.Compare(key, x.key)
	if cmp < 0 {
		x.left = st.put(x.left, key, val)
	} else if cmp > 0 {
		x.right = st.put(x.right, key, val)
	} else {
		x.val = val
	}
	x.n = 1 + st.size(x.left) + st.size(x.right)
	return x
}

func (st *BST) DeleteMin() error {
	if st.IsEmpty() {
		return errors.New("Symbol table underflow")
	}
	st.root = st.deletemin(st.root)
	return nil
}

func (st *BST) deletemin(x *BstNode) *BstNode {
	if x.left == nil {
		return x.right
	}
	x.left = st.deletemin(x.left)
	x.n = st.size(x.left) + st.size(x.right) + 1
	return x
}

func (st *BST) DeleteMax() error {
	if st.IsEmpty() {
		return errors.New("Symbol table underflow")
	}
	st.root = st.deletemax(st.root)
	return nil
}

func (st *BST) deletemax(x *BstNode) *BstNode {
	if x.right == nil {
		return x.left
	}
	x.right = st.deletemax(x.right)
	x.n = st.size(x.left) + st.size(x.right) + 1
	return x
}

func (st *BST) Delete(key interface{}) {
	st.root = st.delete(st.root, key)
}

func (st *BST) delete(x *BstNode, key interface{}) *BstNode {
	if x == nil {
		return nil
	}
	cmp := gotype.Compare(key, x.key)
	if cmp < 0 {
		x.left = st.delete(x.left, key)
	} else if cmp > 0 {
		x.right = st.delete(x.right, key)
	} else {
		if x.right == nil {
			return x.left
		}
		if x.left == nil {
			return x.right
		}
		t := x
		x = st.min(t.right)
		x.right = st.deletemin(t.right)
		x.left = t.left
	}
	x.n = st.size(x.left) + st.size(x.right) + 1
	return x
}

func (st *BST) Min() interface{} {
	if st.IsEmpty() {
		return nil
	}
	return st.min(st.root).key
}

func (st *BST) min(x *BstNode) *BstNode {
	if x.left == nil {
		return x
	} else {
		return st.min(x.left)
	}
}

func (st *BST) Max() interface{} {
	if st.IsEmpty() {
		return nil
	}
	return st.max(st.root).key
}

func (st *BST) max(x *BstNode) *BstNode {
	if x.right == nil {
		return x
	} else {
		return st.max(x.right)
	}
}

func (st *BST) Floor(key interface{}) interface{} {
	x := st.floor(st.root, key)
	if x == nil {
		return nil
	} else {
		return x.key
	}
}

func (st *BST) floor(x *BstNode, key interface{}) *BstNode {
	if x == nil {
		return nil
	}
	cmp := gotype.Compare(key, x.key)
	if cmp == 0 {
		return x
	}
	if cmp < 0 {
		return st.floor(x.left, key)
	}
	t := st.floor(x.right, key)
	if t != nil {
		return t
	}
	return x
}

func (st *BST) Ceiling(key interface{}) interface{} {
	x := st.ceiling(st.root, key)
	if x == nil {
		return nil
	}
	return x.key
}

func (st *BST) ceiling(x *BstNode, key interface{}) *BstNode {
	if x == nil {
		return nil
	}
	cmp := gotype.Compare(key, x.key)
	if cmp == 0 {
		return x
	}
	if cmp < 0 {
		t := st.ceiling(x.left, key)
		if t != nil {
			return t
		} else {
			return x
		}
	}
	return st.ceiling(x.right, key)
}

// Return key of rank k.
func (st *BST) Select(k int) interface{} {
	if k < 0 || k >= st.Size() {
		return nil
	}
	x := st.selec(st.root, k)
	return x.key
}

func (st *BST) selec(x *BstNode, k int) *BstNode {
	if x == nil {
		return nil
	}
	t := st.size(x.left)
	if t > k {
		return st.selec(x.left, k)
	} else if t < k {
		return st.selec(x.right, k-t-1)
	} else {
		return x
	}
}

func (st *BST) Rank(key interface{}) int {
	return st.rank(key, st.root)
}

//  Number of keys in the subtree less than x.key.
func (st *BST) rank(key interface{}, x *BstNode) int {
	if x == nil {
		return 0
	}
	cmp := gotype.Compare(key, x.key)
	if cmp < 0 {
		return st.rank(key, x.left)
	} else if cmp > 0 {
		return 1 + st.size(x.left) + st.rank(key, x.right)
	} else {
		return st.size(x.left)
	}
}

func (st *BST) InOrder() []*BstNode {
	queue := make([]*BstNode, 0, 1000)
	var inorder func(x *BstNode)
	inorder = func(x *BstNode) {
		if x == nil {
			return
		}
		inorder(x.left)
		queue = append(queue, x)
		inorder(x.right)
	}
	inorder(st.root)
	return queue
}

func (st *BST) PreOrder() []*BstNode {
	// 以下代码，不使用闭包，存在bug
	// func preorder(x *Node, queue []*BstNode)
	//     ....
	//     queue = append(queue, x)
	//     preorder(x.left)
	//     preorder(x.right)
	// 解决方案，闭包 or 函数每次都重新返回queue
	queue := make([]*BstNode, 0, 1000) // cap必须足够大，减少创建slice的次数
	var preorder func(x *BstNode)
	preorder = func(x *BstNode) {
		if x == nil {
			return
		}
		queue = append(queue, x)
		preorder(x.left)
		preorder(x.right)
		return
	}
	preorder(st.root)
	return queue
}

func (st *BST) PostOrder() []*BstNode {
	queue := make([]*BstNode, 0, 1000)
	var postorder func(*BstNode)
	postorder = func(x *BstNode) {
		if x == nil {
			return
		}
		postorder(x.left)
		postorder(x.right)
		queue = append(queue, x)
	}
	postorder(st.root)
	return queue
}

// Range search
func (st *BST) RangeKeys(lo interface{}, hi interface{}) []interface{} {
	queue := make([]interface{}, 0)
	var keys func(*BstNode, interface{}, interface{})
	keys = func(x *BstNode, lo interface{}, hi interface{}) {
		if x == nil {
			return
		}
		cmplo := gotype.Compare(lo, x.key)
		cmphi := gotype.Compare(hi, x.key)
		if cmplo < 0 {
			keys(x.left, lo, hi)
		}
		if cmplo <= 0 && cmphi >= 0 {
			queue = append(queue, x.key)
		}
		if cmphi > 0 {
			keys(x.right, lo, hi)
		}
	}
	keys(st.root, lo, hi)
	return queue
}

// Range search
func (st *BST) RangeSize(lo interface{}, hi interface{}) int {
	if gotype.Compare(lo, hi) > 0 {
		return 0
	}
	if st.Contains(hi) {
		return st.Rank(hi) - st.Rank(lo) + 1
	} else {
		return st.Rank(hi) - st.Rank(lo)
	}
}

func (st *BST) Keys() []interface{} {
	return st.RangeKeys(st.Min(), st.Max())
}

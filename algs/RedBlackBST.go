package algs

import (
	"errors"
	"fmt"
	"github.com/sdming/kiss/gotype"
)

const (
    RED = true
    BLACK  = false
)

// BST helper node data type
type RBNode struct {
    key interface{}    // key
    val interface{}    // associated data
    color bool         // color of parent link
    n int              // subtree count
    left *RBNode       // links to left and right subtrees
    right *RBNode
}

func NewRBNode(key interface{}, val interface{}, color bool) *RBNode {
    return &RBNode{key, val,color, 1, nil, nil}
}

func (nd *RBNode) String() string {
    return fmt.Sprintf("(%v, %v, %v, %d, %p, %p)", nd.key, nd.val, nd.color, nd.n, nd.left, nd.right)
}

// A symbol table implemented using a left-leaning red-black BST.
// This is the 2-3 version.
type RedBlackBST struct {
    root *RBNode  // root of the BST
}

func NewRedBlackBST() *RedBlackBST {
    return &RedBlackBST{nil}
}

// is node x red; false if x is null 
func (st *RedBlackBST) isRed( x *RBNode) bool {
    if x == nil {
        return false
    }
    return x.color == RED
}

//number of node in subtree rooted at x; 0 if x is null
func (st *RedBlackBST) size(x *RBNode) int {
    if x == nil {
        return 0
    }
    return x.n
}

//  return number of key-value pairs in this symbol table
func (st *RedBlackBST) Size() int {
    return st.size(st.root)
}

// is this symbol table empty?
func (st *RedBlackBST) IsEmpty() bool {
    return st.root == nil
}

// ****************************
// Standard BST search
// ****************************
// value associated with the given key; null if no such key
func (st *RedBlackBST) Get(key interface{}) interface{} {
    return st.get(st.root, key)
}

// value associated with the given key in subtree rooted at x; null if no such key
func (st *RedBlackBST) get(x *RBNode, key interface{}) interface{} {
    var cmp int
    for x != nil {
        cmp = gotype.Compare(key, x.key)
        if cmp < 0 {
            x = x.left
        } else if cmp > 0 {
            x = x.right
        } else {
            return x.val
        }
    }
    return nil
}

// is there a key-value pair with the given key?
func (st *RedBlackBST) Contains(key interface{}) bool {
    return st.Get(key) != nil
}

// is there a key-value pair with the given key in the subtree rooted at x?
func (st *RedBlackBST) contains(x *RBNode, key interface{}) bool {
    return st.get(x, key) != nil
}

// ****************************
// Red-black insertion
// ****************************
// insert the key-value pair; overwrite the old value with the new value
// if the key is already present
func (st *RedBlackBST) Put(key interface{}, val interface{}) {
    st.root = st.put(st.root, key, val)
    st.root.color = BLACK
}

//  insert the key-value pair in the subtree rooted at h
func (st *RedBlackBST) put(h *RBNode, key interface{}, val interface{}) *RBNode {
    if h == nil {
        return NewRBNode(key, val, RED);
    }
    cmp := gotype.Compare(key, h.key)
    if cmp < 0 {
        h.left = st.put(h.left, key, val)
    } else if cmp > 0 {
        h.right = st.put(h.right, key, val)
    } else {
        h.val = val
    }
    // fix-up any right-leaning links
    if st.isRed(h.right) && !st.isRed(h.left) {
        h = st.rotateLeft(h)
    }
    if st.isRed(h.left) && st.isRed(h.left.left) {
        h = st.rotateRight(h)
    }
    if st.isRed(h.left) && st.isRed(h.right) {
        st.flipColors(h)
    }
    h.n = st.size(h.left) + st.size(h.right) + 1
    return h
}

// ****************************
// Red-black deletion
// ****************************
// delete the key-value pair with the minimum key
func (st *RedBlackBST) DeleteMin() error {
    if st.IsEmpty() {
        return errors.New("BST underflow")
    }
    //  if both children of root are black, set root to red
    if !st.isRed(st.root.left) && !st.isRed(st.root.right) {
        st.root.color = RED
    }
    st.root = st.deleteMin(st.root)
    if !st.IsEmpty() {
        st.root.color = BLACK
    }
    return nil
}

// delete the key-value pair with the minimum key rooted at h
func (st *RedBlackBST) deleteMin(h *RBNode) *RBNode {
    if h.left == nil {
        return nil
    }
    if !st.isRed(h.left) && !st.isRed(h.left.left) {
        h = st.moveRedLeft(h)
    }

    h.left = st.deleteMin(h.left)
    return st.balance(h)
}

//  delete the key-value pair with the maximum key
func (st *RedBlackBST) DeleteMax() error {
    if st.IsEmpty() {
        return errors.New("BST underflow")
    }

    // if both children of root are black, set root to red
    if !st.isRed(st.root.left) && !st.isRed(st.root.right) {
        st.root.color = RED
    }

    st.root = st.deleteMax(st.root)
    if !st.IsEmpty() {
        st.root.color = BLACK
    }
    return nil
}

// delete the key-value pair with the given key
func (st *RedBlackBST) deleteMax(h *RBNode) *RBNode {
    if st.isRed(h.left) {
        h = st.rotateRight(h)
    }
    if h.right == nil {
        return nil
    }
    if !st.isRed(h.right) && !st.isRed(h.right.left) {
        h = st.moveRedRight(h)
    }
    h.right = st.deleteMax(h.right)

    return st.balance(h)
}

// delete the key-value pair with the given key
func (st *RedBlackBST) Delete(key interface{}) {
    if !st.Contains(key) {
        fmt.Println("symbol table does not contain ", key)
        return
    }

    // if both children of root are black, set root to red
    if !st.isRed(st.root.left) && !st.isRed(st.root.right) {
        st.root.color = RED
    }
    st.root = st.delet(st.root, key)
    if !st.IsEmpty() {
        st.root.color = BLACK
    }
}

//  delete the key-value pair with the given key rooted at h
func (st *RedBlackBST) delet(h *RBNode, key interface{}) *RBNode {
    if gotype.Compare(key, h.key) < 0 {
        if !st.isRed(h.left) && !st.isRed(h.left.left) {
            h = st.moveRedLeft(h)
        }
        h.left = st.delet(h.left, key)
    } else {
        if st.isRed(h.left) {
            h = st.rotateRight(h)
        }
        if gotype.Compare(key, h.key) == 0 && h.right == nil {
            return nil
        }
        if !st.isRed(h.right) && !st.isRed(h.right.left) {
            h = st.moveRedRight(h)
        }
        if gotype.Compare(key, h.key) == 0 {
            x := st.min(h.right)
            h.key = x.key
            h.val = x.val
            h.right = st.deleteMin(h.right)
        } else {
            h.right = st.delet(h.right, key)
        }
    }
    return st.balance(h)
}

//***********************************************
//   red-black tree helper functions
//***********************************************

//  make a left-leaning link lean to the right
func (st *RedBlackBST) rotateRight(h *RBNode) *RBNode {
    x := h.left
    h.left = x.right
    x.right = h
    x.color = h.color
    x.right.color = RED
    x.n = h.n
    h.n = st.size(h.left) + st.size(h.right) + 1
    return x
}

//  make a right-leaning link lean to the left
func (st *RedBlackBST) rotateLeft(h *RBNode) *RBNode {
    x := h.right
    h.right = x.left
    x.left = h
    x.color = h.color
    h.color = RED
    x.n = h.n
    h.n = st.size(h.left) + st.size(h.right) + 1
    return x
}

//  flip the colors of a node and its two children
func (st *RedBlackBST) flipColors(h *RBNode) {
    h.color = !h.color
    h.left.color = !h.left.color
    h.right.color = !h.right.color
}

// Assuming that h is red and both h.left and h.left.left
// are black, make h.left or one of its children red.
func (st *RedBlackBST) moveRedLeft(h *RBNode) *RBNode {
    st.flipColors(h)
    if st.isRed(h.right.left) {
        h.right = st.rotateRight(h.right)
        h = st.rotateLeft(h)
    }
    return h
}

//  Assuming that h is red and both h.right and h.right.left
// are black, make h.right or one of its children re
func (st *RedBlackBST) moveRedRight(h *RBNode) *RBNode {
    //assert (h != null);
    //assert isRed(h) && !isRed(h.right) && !isRed(h.right.left);
    st.flipColors(h)
    if st.isRed(h.left.left) {
        h = st.rotateRight(h)
    }
    return h
}

//  restore red-black tree invariant
func (st *RedBlackBST) balance(h *RBNode) *RBNode {
    if st.isRed(h.right) {
        h = st.rotateLeft(h)
    }
    if st.isRed(h.left) && st.isRed(h.left.left) {
        h = st.rotateRight(h)
    }
    if st.isRed(h.left) && st.isRed(h.right) {
        st.flipColors(h)
    }
    h.n = st.size(h.left) + st.size(h.right) + 1
    return h
}

//  height of tree; 0 if empty
func (st *RedBlackBST) Height() int {
    return st.height(st.root)
}

func (st *RedBlackBST) height(x *RBNode) int {
    if x == nil {
        return -1
    }
    i := st.height(x.left)
    j := st.height(x.right)
    if i > j {
        return 1 + i
    } else {
        return 1 + j
    }
}

// the smallest key; null if no such key
func (st *RedBlackBST) Min() interface{} {
    if st.IsEmpty() {
        return nil
    }
    return st.min(st.root).key
}

//  the smallest key in subtree rooted at x; null if no such key
func (st *RedBlackBST) min(x *RBNode) *RBNode {
    if x.left == nil {
        return x
    } else {
        return st.min(x.left)
    }
}

//  the largest key; null if no such key
func (st *RedBlackBST) Max() interface{} {
    if st.IsEmpty() {
        return nil
    }
    return st.max(st.root).key
}

//  the largest key in the subtree rooted at x; null if no such key
func (st *RedBlackBST) max(x *RBNode) *RBNode {
    if x.right == nil {
        return x
    }
    return st.max(x.right)
}

//  the largest key less than or equal to the given key
func (st *RedBlackBST) Floor(key interface{}) interface{} {
    x := st.floor(st.root, key)
    if x == nil {
        return nil
    } else {
        return x.key
    }
}

//  the largest key in the subtree rooted at x less than or equal to the given key
func (st *RedBlackBST) floor(x *RBNode, key interface{}) *RBNode {
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

//  the smallest key greater than or equal to the given key
func (st *RedBlackBST) Ceiling(key interface{}) interface{} {
    x := st.ceiling(st.root, key)
    if x == nil {
        return nil
    } else {
        return x.key
    }
}

//  the smallest key in the subtree rooted at x greater than or equal to the given key
func (st *RedBlackBST) ceiling(x *RBNode, key interface{}) *RBNode {
    if x == nil {
        return nil
    }
    cmp := gotype.Compare(key, x.key)
    if cmp == 0 {
        return x
    }
    if cmp > 0 {
        return st.ceiling(x.right, key)
    }
    t := st.ceiling(x.left, key)
    if t != nil {
        return t
    } else {
        return x
    }
}

//  the key of rank k
func (st *RedBlackBST) Select(k int) interface{} {
    if k < 0 || k >= st.Size() {
        return nil
    }
    x := st.selec(st.root, k)
    return x.key
}

//  the key of rank k in the subtree rooted at x
func (st * RedBlackBST) selec(x *RBNode, k int) *RBNode {
    t := st.size(x.left)
    if t > k {
        return st.selec(x.left, k)
    } else if t < k {
        return st.selec(x.right, k-t-1)
    }
    return x
}

//  number of keys less than key
func (st *RedBlackBST) Rank(key interface{}) int {
    return st.rank(key, st.root)
}

//  number of keys less than key in the subtree rooted at x
func (st *RedBlackBST) rank(key interface{}, x *RBNode) int {
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

// all of the keys, as an Iterable
func (st *RedBlackBST) Keys() []interface{} {
    return st.RangeKeys(st.Min(), st.Max())
}

// Range search
func (st *RedBlackBST) RangeKeys(lo interface{}, hi interface{}) []interface{} {
	queue := make([]interface{}, 0, 1000)
	var keys func(*RBNode, interface{}, interface{})
	keys = func(x *RBNode, lo interface{}, hi interface{}) {
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
func (st *RedBlackBST) RangeSize(lo interface{}, hi interface{}) int {
	if gotype.Compare(lo, hi) > 0 {
		return 0
	}
	if st.Contains(hi) {
		return st.Rank(hi) - st.Rank(lo) + 1
	} else {
		return st.Rank(hi) - st.Rank(lo)
	}
}

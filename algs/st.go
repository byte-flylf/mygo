package algs

// Symbol table.
// A symbol table A symbol table is a data structure that associates a value with a key.
// It supports two primary operations: insert (put) a new pair into the table and search for (get) the value associated with a given key
type ST interface {
	Put(key interface{}, val interface{})
	Get(key interface{}) interface{}
	Delete(key interface{})
	Contains(key interface{}) bool
	IsEmpty() bool
	Size() int
	Keys() []interface{} // public Iterable<Key> keys()
}

// Ordered symbol tables
// In typical applications, keys are Comparable objects,
// so the option exists of using the code a.compareTo(b) to compare two keys a and b.
// Several symbol-table implementations take advantage of order among the keys that is implied by Comparable
// to provide efficient implementations of the put() and get() operations. More important, in such implementations,
// we can think of the symbol table as keeping the keys in order and consider a significantly expanded API that defines numerous
// natural and useful operations involving relative key order. For applications where keys are Comparable,
type OrderST interface {
	ST

	Min() interface{}
	Max() interface{}
    // find the largest key that is less than or equal to the given key
	Floor(key interface{}) interface{}
    // find the smallest key that is greater than or equal to the given key
	Ceiling(key interface{}) interface{}
    // find the number of keys less than a given key
	Rank(key interface{}) int
    // find the key with a given rank
	Select(k int) interface{}
	DeleteMin() error
	DeleteMax() error

	RangeSize(lo interface{}, hi interface{}) int           // Range Count
	RangeKeys(lo interface{}, hi interface{}) []interface{} // Range Search
}

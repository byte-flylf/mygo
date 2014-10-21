package algs

//  Weighted quick-union (without path compression)
//  interface: uf.go
type WeightedQuickUnionUF struct {
	id    []int // id[i] = parent of i
	sz    []int //  sz[i] = number of objects in subtree rooted at i
	count int   //  number of components
}

// Create an empty union-find data structure with N isolated components 0 through N-1.
func NewWeightedQuickUnionUF(n int) *WeightedQuickUnionUF {
	arr := make([]int, n)
	sz := make([]int, n)
	for i := range arr {
		arr[i] = i
		sz[i] = 1
	}
	return &WeightedQuickUnionUF{id: arr, sz: sz, count: n}
}

// Return the number of components.
func (self *WeightedQuickUnionUF) Count() int {
	return self.count
}

// Return the component identifier for component containing p.
// O(lgN)
func (self *WeightedQuickUnionUF) Find(p int) int {
	return self.root(p)
}

func (self *WeightedQuickUnionUF) root(p int) int {
	for p != self.id[p] {
		p = self.id[p]
	}
	return p
}

// Are objects p and q in the same component
func (self *WeightedQuickUnionUF) Connected(p, q int) bool {
	return self.root(p) == self.root(q)
}

// Merge components containing p and q.
// O(lgN)
func (self *WeightedQuickUnionUF) Union(p, q int) {
	i, j := self.root(p), self.root(q)
	if i == j {
		return
	}
	if self.sz[i] < self.sz[j] {
		self.id[i] = j
		self.sz[j] += self.sz[i]
	} else {
		self.id[j] = i
		self.sz[i] += self.sz[j]
	}
	self.count--
}

package algs

// Quick-union algorithm.
type QuickUnionUF struct {
	id    []int // id[i] = parent of i
	count int   // number of components
}

func NewQuickUnionUF(n int) *QuickUnionUF {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = i
	}
	return &QuickUnionUF{id: arr, count: n}
}

func (self *QuickUnionUF) Count() int {
	return self.count
}

func (self *QuickUnionUF) Find(p int) int {
	return self.root(p)
}

func (self *QuickUnionUF) root(p int) int {
	for p != self.id[p] {
		p = self.id[p]
	}
	return p
}

func (self *QuickUnionUF) Connected(p, q int) bool {
	return self.Find(p) == self.Find(q)
}

func (self *QuickUnionUF) Union(p, q int) {
	i, j := self.root(p), self.root(q)
	if i == j {
		return
	}
	self.id[i] = j
	self.count--
}

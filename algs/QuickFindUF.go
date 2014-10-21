package algs

// Quick-find algorithm.  interface: uf.go
type QuickFindUF struct {
    id []int     // id[i] = component identifier of i
    count int    // number of components
}

func NewQuickFindUF(n int) *QuickFindUF {
    arr := make([]int, n)
    for i := range arr {
        arr[i] = i
    }
    return &QuickFindUF{id: arr, count: n}
}

func (self *QuickFindUF) Count() int {
    return self.count
}

func (self *QuickFindUF) Find(p int) int {
    return self.id[p]
}

func (self *QuickFindUF) Connected(p, q int) bool {
    return self.id[p] == self.id[q]
}

func (self *QuickFindUF) Union(p, q int) {
    if self.Connected(p, q) {
        return
    }
    pid := self.id[p]
    for i := range(self.id) {
        if self.id[i] == pid {
            self.id[i] = self.id[q]
        }
    }
    self.count--
}

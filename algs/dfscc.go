package algs

type DfsCC struct {
	marked []bool
	id     []int
	size   []int
	count  int
}

// This implementation uses depth-first search.
func NewDfsCC(G Graph) *DfsCC {
	cc := new(DfsCC)
	cc.marked = make([]bool, G.V())
	cc.id = make([]int, G.V())
	cc.size = make([]int, G.V())
	for v := 0; v < G.V(); v++ {
		if !cc.marked[v] {
			cc.dfs(G, v)
			cc.count++
		}
	}

	return cc
}

// depth-first search
func (cc *DfsCC) dfs(G Graph, v int) {
	cc.marked[v] = true
	cc.id[v] = cc.count
	cc.size[cc.count]++
	for _, w := range G.Adj(v) {
		if !cc.marked[w] {
			cc.dfs(G, w)
		}
	}
}

// Returns the component id of the connected component containing vertex
func (cc DfsCC) ID(v int) int {
	return cc.id[v]
}

// Returns the number of vertices in the connected component containing vertex
func (cc DfsCC) Size(v int) int {
	return cc.size[cc.id[v]]
}

// Returns the number of vertices in the connected component containing vertex
func (cc DfsCC) Count() int {
	return cc.count
}

// Are vertices v and w in the same connected component?
func (cc DfsCC) Connected(v, w int) bool {
	return cc.ID(v) == cc.ID(w)
}

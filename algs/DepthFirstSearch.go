package algs

/*
*  The DepthFirstSearch class represents a data type for
*  determining the vertices connected to a given source vertex s
*  in an undirected graph. For versions that find the paths, see
*  {@link DepthFirstPaths} and {@link BreadthFirstPaths}.
*
*  This implementation uses depth-first search.
*  For additional documentation, see <a href="/algs4/41graph">Section 4.1</a> of
*  <i>Algorithms, 4th Edition</i> by Robert Sedgewick and Kevin Wayne.
 */
type DepthFirstSearch struct {
	marked []bool // marked[v] = is there an s-v path?
	count  int    // number of vertices connected to s
}

func NewDepthFirstSearch(G Graph, s int) *DepthFirstSearch {
	self := new(DepthFirstSearch)
	self.marked = make([]bool, G.V())
	self.dfs(G, s)
	return self
}

// depth first search from v
func (self *DepthFirstSearch) dfs(G Graph, v int) {
	self.count++
	self.marked[v] = true
	for _, w := range G.Adj(v) {
		if !self.marked[w] {
			self.dfs(G, w)
		}
	}
}

// Is there a path between the source vertex 's' and vertex 'v'?
func (self *DepthFirstSearch) Marked(v int) bool {
	return self.marked[v]
}

// Returns the number of vertices connected to the source vertex
func (self *DepthFirstSearch) Count() int {
	return self.count
}

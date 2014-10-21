package algs

//  The DepthFirstPaths class represents a data type for finding
//  paths from a source vertex to every other vertex
//  in an undirected graph.
type DepthFirstPaths struct {
	marked []bool // marked[v] = is there an s-v path?
	edgeTo []int  // edgeTo[v] = last edge on s-v path
	s      int    // source vertex
}

// Computes a path between 's' and every other vertex in graph
func NewDepthFirstPaths(G Graph, s int) *DepthFirstPaths {
	self := new(DepthFirstPaths)
	self.s = s
	self.edgeTo = make([]int, G.V())
	self.marked = make([]bool, G.V())
	self.dfs(G, s)
	return self
}

// depth first search from v
func (self *DepthFirstPaths) dfs(G Graph, v int) {
	self.marked[v] = true
	for _, w := range G.Adj(v) {
		if !self.marked[w] {
			self.edgeTo[w] = v
			self.dfs(G, w)
		}
	}
}

// Is there a path between the source vertex and vertex 'v'
func (self *DepthFirstPaths) HasPathTo(v int) bool {
	return self.marked[v]
}

// Returns a path between the source vertex and vertex 'v'
func (self *DepthFirstPaths) PathTo(v int) []int {
	if !self.HasPathTo(v) {
		return nil
	}
	out := make([]int, 0)
	for x := v; x != self.s; x = self.edgeTo[x] {
		out = append([]int{x}, out...)
	}
	out = append([]int{self.s}, out...)
	return out
}

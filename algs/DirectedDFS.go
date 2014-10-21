package algs

// The DirectedDFS class represents a data type for
//  determining the vertices reachable from a given source vertex
type DirectedDFS struct {
	marked []bool // marked[v] = true if v is reachable from source
	count  int    // number of vertices reachable from s
}

func NewDirectedDFS(g *Digraph, s int) *DirectedDFS {
	p := new(DirectedDFS)
	p.marked = make([]bool, g.V())
	p.dfs(g, s)
	return p
}

// Computes the vertices in digraph  that are connected to any of the source vertices
func NewDirectedDFSFromMultSrc(g *Digraph, sources ...int) *DirectedDFS {
	p := new(DirectedDFS)
	p.marked = make([]bool, g.V())
	for _, v := range sources {
		if !p.marked[v] {
			p.dfs(g, v)
		}
	}
	return p
}

func (p *DirectedDFS) dfs(g *Digraph, v int) {
	p.count++
	p.marked[v] = true
	for _, w := range g.Adj(v) {
		if !p.marked[w] {
			p.dfs(g, w)
		}
	}
}

// Returns the number of vertices reachable from the source vertex
func (p *DirectedDFS) Count() int {
	return p.count
}

// Is there a directed path from the source vertex
func (p *DirectedDFS) Marked(v int) bool {
	return p.marked[v]
}

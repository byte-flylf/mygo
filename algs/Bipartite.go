package algs

//  The Bipartite class represents a data type for
//  determining whether an undirected graph is bipartite or whethe
//  it has an odd-length cycle.
type Bipartite struct {
	isBipartite bool   // is the graph bipartite?
	color       []bool // color[v] gives vertices on one side of bipartition
	marked      []bool // marked[v] = true if v has been visited in DFS
	edgeTo      []int  // edgeTo[v] = last edge on path to v
}

func NewBipartite(g Graph) *Bipartite {
	p := new(Bipartite)
	p.isBipartite = true
	p.color = make([]bool, g.V())
	p.marked = make([]bool, g.V())
	p.edgeTo = make([]int, g.V())

	for v := 0; v < g.V(); v++ {
		if !p.marked[v] {
			p.dfs(g, v)
		}
	}

	return p
}

func (b *Bipartite) dfs(g Graph, v int) {
	b.marked[v] = true
	for _, w := range g.Adj(v) {
		if !b.marked[w] {
			b.edgeTo[w] = v
			b.color[w] = !b.color[v]
			b.dfs(g, w)
		} else if b.color[w] == b.color[v] {
			b.isBipartite = false
		}
	}
}

func (b Bipartite) IsBipartite() bool {
	return b.isBipartite
}

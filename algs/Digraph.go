package algs

// represents a directed graph of vertices
type Digraph struct {
	v   int
	e   int
	adj [][]int
}

func NewDigraph(v int) *Digraph {
	if v < 0 {
		panic("Number of vertices in a Digraph must be nonnegative")
	}
	p := new(Digraph)
	p.v = v
	p.e = 0
	p.adj = make([][]int, v)
	for i := range p.adj {
		p.adj[i] = make([]int, 0)
	}

	return p
}

func NewDigraphForFile(file string) *Digraph {
	var numbers []int
	var err error

	numbers, err = ReadInts(file)
	if err != nil {
		panic(err)
	}

	v := numbers[0]
	G := NewDigraph(v)

	for i := 0; i < numbers[1]; i++ {
		v := numbers[2+i*2]
		w := numbers[2+i*2+1]
		G.AddEdge(v, w)
	}
	return G
}

// Initializes a new digraph that is a deep copy of G
func CopyDigraph(g *Digraph) *Digraph {
	g2 := NewDigraph(g.v)
	g2.e = g.e
	for v := 0; v < g.v; v++ {
		reverse := make([]int, len(g.adj[v]))
		copy(reverse, g.adj[v])
	}
	return g2
}

// Returns the number of vertices in the digraph.
func (g Digraph) V() int {
	return g.v
}

// Returns the number of edges in the digraph.
func (g Digraph) E() int {
	return g.e
}

// Adds the directed edge v->w to the digraph.
// @param v the tail vertex
// @param w the head vertex
func (g *Digraph) AddEdge(v, w int) {
	if v < 0 || v >= g.v {
		panic("vertex v out of bounds")
	}
	if w < 0 || w >= g.v {
		panic("vertex w out of bounds")
	}
	g.adj[v] = append(g.adj[v], w)
	g.e++
}

// Returns the vertices adjacent from vertex v in the digraph.
func (g *Digraph) Adj(v int) []int {
	//return g.adj[v]
	out := make([]int, len(g.adj[v]))
	for i := 0; i < len(g.adj[v]); i++ {
		out[i] = g.adj[v][len(g.adj[v])-1-i]
	}
	return out
}

// Returns the reverse of the digraph
func (g *Digraph) Reverse() *Digraph {
	r := NewDigraph(g.V())
	for v := 0; v < g.V(); v++ {
		for _, w := range g.Adj(v) {
			r.AddEdge(w, v)
		}
	}
	return r
}

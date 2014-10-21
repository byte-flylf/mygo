package algs

import (
	"fmt"
)

// The Graph class represents an undirected graph of vertices
// named 0 through V-1
// It supports the following two primary operations: add an edge to the graph,
// iterate over all of the vertices adjacent to a vertex. It also provides
// methods for returning the number of vertices V and the number
// edges <em>E</em>. Parallel edges and self-loops are permitted.
// This implementation uses an adjacency-lists representation, which
// is a vertex-indexed array of {@link Bag} objects.
// All operations take constant time (in the worst case) except
// iterating over the vertices adjacent to a given vertex, which takes
// time proportional to the number of such vertices.
//
// For additional documentation, see <a href="http://algs4.cs.princeton.edu/41undirected">Section 4.1</a> of
// Algorithms, 4th Edition by Robert Sedgewick and Kevin Wayne.
type AdjGraph struct {
	v   int
	e   int
	adj []Bager
}

// Initializes an empty graph with V vertices and 0 edges
func NewAdjGraph(v int) *AdjGraph {
	if v < 0 {
		panic("Number of vertices must be nonnegative")
	}

	g := new(AdjGraph)
	g.v = v
	g.adj = make([]Bager, v)
	for i := 0; i < v; i++ {
		g.adj[i] = NewArrayBag() // NewLinkedListBag
	}
	return g
}

// Initializes a graph from an input stream.
func NewAdjGraphForFile(file string) *AdjGraph {
	var numbers []int
	var err error

	numbers, err = ReadInts(file)
	if err != nil {
		panic(err)
	}

	v := numbers[0]
	G := NewAdjGraph(v)

	for i := 0; i < numbers[1]; i++ {
		v := numbers[2+i*2]
		w := numbers[2+i*2+1]
		G.AddEdge(v, w)
	}
	return G
}

// Returns the number of vertices in the graph.
func (self *AdjGraph) V() int {
	return self.v
}

// Returns the number of edges in the graph.
func (self *AdjGraph) E() int {
	return self.e
}

// Adds the undirected edge v-w to the graph.
func (self *AdjGraph) AddEdge(v, w int) {
	if v < 0 || v >= self.v {
		panic("IndexOutOfBoundsException")
	}
	if w < 0 || w >= self.v {
		panic("IndexOutOfBoundsException")
	}
	self.e++
	self.adj[v].Add(w)
	self.adj[w].Add(v)
}

// Returns the vertices adjacent to vertex
func (self *AdjGraph) Adj(v int) []int {
	if v < 0 || v >= self.v {
		panic("IndexOutOfBoundsException")
	}
	out := make([]int, 0)
	for w := range self.adj[v].ChannelIterator() {
		i, _ := w.(int)
		out = append(out, i)
	}
	ReverseIntSlice(out)
	return out
}

func (self *AdjGraph) String() string {
	s := fmt.Sprintf("%d vertices, %d edges\n", self.v, self.e)
	for i := 0; i < self.v; i++ {
		s = s + fmt.Sprintf("%d : ", i)
		for _, w := range self.Adj(i) {
			s = s + fmt.Sprintf("%v ", w)
		}
		s += "\n"
	}
	return s
}

package algs

// Undirected graph data type.
type Graph interface {
	V() int           // number of vertices
	E() int           // number of edges
	AddEdge(v, w int) // add edge v-w to this graph
	Adj(v int) []int  // vertices adjacent to v
}

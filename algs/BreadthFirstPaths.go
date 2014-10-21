package algs

import (
	"math"
)

const INFINITY int = math.MaxInt64

// This implementation uses breadth-first search.
type BreadthFirstPaths struct {
	marked []bool // marked[v] = is there an s-v path
	edgeTo []int  // edgeTo[v] = previous edge on shortest s-v path
	distTo []int  // distTo[v] = number of edges shortest s-v path
}

func NewBreadthFirstPaths(G Graph, s int) *BreadthFirstPaths {
	p := new(BreadthFirstPaths)
	p.marked = make([]bool, G.V())
	p.edgeTo = make([]int, G.V())
	p.distTo = make([]int, G.V())
	p.bfs(G, s)
	return p
}

// breadth-first search from a single source
func (p *BreadthFirstPaths) bfs(G Graph, s int) {
	q := NewArrayQueue()
	for v := 0; v < G.V(); v++ {
		p.distTo[v] = INFINITY
	}
	p.distTo[s] = 0
	p.marked[s] = true
	q.Enqueue(s)

	for !q.IsEmpty() {
		v := q.Dequeue()
		i, _ := v.(int)
		for _, w := range G.Adj(i) {
			if !p.marked[w] {
				p.edgeTo[w] = i
				p.distTo[w] = p.distTo[i] + 1
				p.marked[w] = true
				q.Enqueue(w)
			}
		}
	}
}

func (p *BreadthFirstPaths) HasPathTo(v int) bool {
	return p.marked[v]
}

func (p *BreadthFirstPaths) DistTo(v int) int {
	return p.distTo[v]
}

// Returns a path between the source vertex and vertex 'v'
func (p *BreadthFirstPaths) PathTo(v int) []int {
	if !p.HasPathTo(v) {
		return nil
	}

	out := make([]int, 0)
	var x int
	for x = v; p.distTo[x] != 0; x = p.edgeTo[x] {
		out = append([]int{x}, out...)
	}
	out = append([]int{x}, out...)
	return out
}

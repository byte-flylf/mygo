package algs

/*
* The Cycle class represents a data type for
* determining whether an undirected graph has a cycle.
* hasCycle operation determines whether the graph has
* a cycle and, if so, the cycle operation returns one.
*
* This implementation uses depth-first search.
* The constructor takes time proportional to  V+E (in the worst case)
* where V is the number of vertices and <em>Eem> is the number of edges.
* Afterwards, the hasCycle operation takes constant time;
* the cycle operation takes time proportional
* to the length of the cycle.
*
 */

type Cycle struct {
	marked []bool
	edgeTo []int
	cycle  Stack
}

func NewCycle(g Graph) *Cycle {
	c := new(Cycle)
	if c.hasSelfLoop(g) {
		return c
	}
	if c.hasParallelEdges(g) {
		return c
	}

	c.marked = make([]bool, g.V())
	c.edgeTo = make([]int, g.V())
	for v := 0; v < g.V(); v++ {
		if !c.marked[v] {
			c.dfs(g, -1, v)
		}
	}

	return c
}

// does this graph have a self loop?
// side effect: initialize cycle to be self loop
func (c *Cycle) hasSelfLoop(g Graph) bool {
	for v := 0; v < g.V(); v++ {
		for _, w := range g.Adj(v) {
			if v == w {
				c.cycle = NewArrayStack()
				c.cycle.Push(v)
				c.cycle.Push(v)
				return true
			}
		}
	}
	return false
}

// does this graph have two parallel edges?
// side effect: initialize cycle to be two parallel edges
func (c *Cycle) hasParallelEdges(g Graph) bool {
	c.marked = make([]bool, g.V())

	for v := 0; v < g.V(); v++ {
		for _, w := range g.Adj(v) {
			if c.marked[w] {
				c.cycle = NewArrayStack()
				c.cycle.Push(v)
				c.cycle.Push(w)
				c.cycle.Push(v)
				return true
			}
			c.marked[w] = true
		}

		for _, w := range g.Adj(v) {
			c.marked[w] = false
		}
	}

	return false
}

func (c *Cycle) dfs(g Graph, u, v int) {
	c.marked[v] = true

	for _, w := range g.Adj(v) {
		// short circuit if cycle already found
		if c.cycle != nil {
			return
		}

		if !c.marked[w] {
			c.edgeTo[w] = v
			c.dfs(g, v, w)
		} else if w != u {
			c.cycle = NewArrayStack()
			for x := v; x != w; x = c.edgeTo[x] {
				c.cycle.Push(x)
			}
			c.cycle.Push(w)
			c.cycle.Push(v)
		}
	}
}

func (c *Cycle) HasCycle() bool {
	return c.cycle != nil
}

func (c *Cycle) Cycle() []int {
	out := make([]int, c.cycle.Size())

	for i := 0; i < len(out); i++ {
		out[i] = c.cycle.Pop().(int)
	}
	return out
}

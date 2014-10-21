package algs

// a data type for determining whether a digraph has a directed cycle.
type DirectedCycle struct {
	marked  []bool
	edgeTo  []int
	onStack []bool
	cycle   Stack
}

func NewDirectedCycle(g *Digraph) *DirectedCycle {
	p := new(DirectedCycle)
	p.marked = make([]bool, g.V())
	p.edgeTo = make([]int, g.V())
	p.onStack = make([]bool, g.V())
	p.cycle = NewArrayStack()
	for v := 0; v < g.V(); v++ {
		if !p.marked[v] {
			p.dfs(g, v)
		}
	}

	return p
}

func (p *DirectedCycle) dfs(g *Digraph, v int) {
	p.onStack[v] = true
	p.marked[v] = true
	for _, w := range g.Adj(v) {
		if p.HasCycle() {
			return
		}

		if !p.marked[w] {
			p.edgeTo[w] = v
			p.dfs(g, w)
		} else if p.onStack[w] {
			for x := v; x != w; x = p.edgeTo[x] {
				p.cycle.Push(x)
			}
			p.cycle.Push(w)
			p.cycle.Push(v)
		}
	}
	p.onStack[v] = false
}

// Does the digraph have a directed cycle?
func (p *DirectedCycle) HasCycle() bool {
	return !p.cycle.IsEmpty()
}

// Returns a directed cycle if the digraph has a directed cycle
func (p *DirectedCycle) Cycle() []int {
	out := make([]int, p.cycle.Size())
	for i := 0; !p.cycle.IsEmpty(); i++ {
		item := p.cycle.Pop()
		out[i] = item.(int)
	}
	return out
}

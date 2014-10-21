package algs

// Compute preorder and postorder for a digraph or edge-weighted digraph.
type DepthFirstOrder struct {
	marked      []bool
	pre         Queue
	post        Queue
	reversePost Stack
}

func NewDepthFirstOrder(g *Digraph) *DepthFirstOrder {
	p := new(DepthFirstOrder)
	p.marked = make([]bool, g.V())
	p.pre = NewArrayQueue()
	p.post = NewArrayQueue()
	p.reversePost = NewArrayStack()

	for v := 0; v < g.V(); v++ {
		if !p.marked[v] {
			p.dfs(g, v)
		}
	}

	return p
}

func (p *DepthFirstOrder) dfs(g *Digraph, v int) {
	p.pre.Enqueue(v)

	p.marked[v] = true
	for _, w := range g.Adj(v) {
		if !p.marked[w] {
			p.dfs(g, w)
		}
	}
	p.post.Enqueue(v)
	p.reversePost.Push(v)
}

func (p *DepthFirstOrder) Pre() []int {
	out := make([]int, p.pre.Size())
	for i := 0; !p.pre.IsEmpty(); i++ {
		item := p.pre.Dequeue()
		out[i] = item.(int)
	}
	return out
}

func (p *DepthFirstOrder) Post() []int {
	out := make([]int, p.post.Size())
	for i := 0; !p.post.IsEmpty(); i++ {
		item := p.post.Dequeue()
		out[i] = item.(int)
	}
	return out
}

func (p *DepthFirstOrder) ReversePost() []int {
	out := make([]int, p.reversePost.Size())
	for i := 0; !p.reversePost.IsEmpty(); i++ {
		item := p.reversePost.Pop()
		out[i] = item.(int)
	}
	return out
}

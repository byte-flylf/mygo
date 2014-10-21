package algs

// A graph, implemented using an adjacency matrix.
// Parallel edges are disallowed; self-loops are allowd.
type AdjMatrixGraph struct {
    v int
    e int
    adj [][]bool
}

func NewAdjMatrixGraph(v int) *AdjMatrixGraph {
    g := new(AdjMatrixGraph)
    g.v = v
    g.adj = make([][]bool, v)
    for i := 0; i < v; i++ {
        g.adj[i] = make([]bool, v)
    }
    return g
}

func (self *AdjMatrixGraph) V() int {
    return self.v
}

func (self *AdjMatrixGraph) E() int {
    return self.e
}

func (self *AdjMatrixGraph) AddEdge(v, w int) {
    if !self.adj[v][w] {
        self.e++
    }
    self.adj[v][w] = true
    self.adj[w][v] = true
}

func (self *AdjMatrixGraph) Contains(v, w int) bool {
    return self.adj[v][w]
}

// return list of neighbors of v
func (self *AdjMatrixGraph) Adj(v int)  (chan interface{}) {
    return self.ChannelIterator(v)
}

func (self *AdjMatrixGraph) ChannelIterator(v int) (chan interface{}) {
    ch := make(chan interface{})
    go func() {
        for idx, b := range self.adj[v] {
            if b  {
                ch <- idx
            }
        }
        close(ch)
    } ()
    return ch
}

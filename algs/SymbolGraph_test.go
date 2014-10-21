package algs

import (
	"github.com/bmizerany/assert"
	"sort"
	"testing"
)

func adj(G Graph, source string, sg *SymbolGraph) []string {
	dest := make([]string, 0)
	if sg.Contains(source) {
		s := sg.Index(source)
		for _, v := range G.Adj(s) {
			dest = append(dest, sg.Name(v))
		}
	}
	sort.Strings(dest)
	return dest
}

func TestSymbolGraph(t *testing.T) {
	filename := "./algs4-data/routes.txt"
	sg := NewSymbolGraph(filename, " ")

	g := sg.G()
	assert.Equal(t, adj(g, "JFK", sg), []string{"ATL", "MCO", "ORD"})
	assert.Equal(t, adj(g, "LAX", sg), []string{"LAS", "PHX"})
}

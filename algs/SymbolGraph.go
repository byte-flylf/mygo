package algs

import (
	"strings"
)

/*
 *  The SymbolGraph class represents an undirected graph, where the
 *  vertex names are arbitrary strings.
 *  By providing mappings between string vertex names and integers,
 *  it serves as a wrapper around the Graph data type,
 *  which assumes the vertex names are integers
 *  between 0 and V - 1.
 *  It also supports initializing a symbol graph from a file.
 */

type SymbolGraph struct {
	st   map[string]int // string -> index
	keys []string       // index  -> string
	g    Graph
}

func NewSymbolGraph(filename string, delimiter string) *SymbolGraph {
	p := new(SymbolGraph)
	p.st = make(map[string]int)

	lines, err := ReadLines(filename)
	if err != nil {
		panic("fail to read file")
	}

	for _, line := range lines {
		parts := strings.Split(line[:len(line)-1], delimiter)
		for _, v := range parts {
			if _, ok := p.st[v]; !ok {
				p.st[v] = len(p.st)
			}
		}
	}

	p.keys = make([]string, len(p.st))
	for k, v := range p.st {
		p.keys[v] = k
	}

	p.g = NewAdjGraph(len(p.st))
	for _, line := range lines {
		parts := strings.Split(line[:len(line)-1], delimiter)
		v := p.st[parts[0]]
		for _, s := range parts[1:] {
			w := p.st[s]
			p.g.AddEdge(v, w)
		}
	}

	return p
}

func (sg SymbolGraph) Contains(s string) bool {
	if _, ok := sg.st[s]; ok {
		return true
	}
	return false
}

func (sg SymbolGraph) Index(s string) int {
	if sg.Contains(s) {
		return sg.st[s]
	}
	return -1
}

func (sg SymbolGraph) Name(v int) string {
	return sg.keys[v]
}

func (sg SymbolGraph) G() Graph {
	return sg.g
}

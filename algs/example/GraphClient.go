// Typical graph-processing code.
//
//  $ ./GraphClient algs4-data/tinyG.txt
// 13 vertices, 13 edges
// 0 : 5 1 2 6
// 1 : 0
// 2 : 0
// 3 : 4 5
// 4 : 3 6 5
// 5 : 0 4 3
// 6 : 4 0
// 7 : 8
// 8 : 7
// 9 : 12 10 11
// 10 : 9
// 11 : 12 9
// 12 : 9 11
// vertex of maximum degree = 4
// average degree = 2
// number of self loops = 0
//

package main

import (
	"fmt"
	"os"

	. "algs"
)

// degree of v
func degree(G Graph, v int) int {
	d := 0
	for _ = range G.Adj(v) {
		d++
	}
	return d
}

// maximum degree
func maxDegree(G Graph) int {
	max := 0
	for v := 0; v < G.V(); v++ {
		if degree(G, v) > max {
			max = degree(G, v)
		}
	}
	return max
}

// average degree
func avgDegree(G Graph) int {
	return 2 * G.E() / G.V()
}

// number of self-loops
func numberOfSelfLoops(G Graph) int {
	count := 0
	for i := 0; i < G.V(); i++ {
		for w := range G.Adj(i) {
			if i == w {
				count++
			}
		}
	}
	return count / 2
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <input>", os.Args[0])
		os.Exit(0)
	}
	G := NewAdjGraphForFile(os.Args[1])
	fmt.Println(G)

	fmt.Println("vertex of maximum degree =", maxDegree(G))
	fmt.Println("average degree =", avgDegree(G))
	fmt.Println("number of self loops =", numberOfSelfLoops(G))
}

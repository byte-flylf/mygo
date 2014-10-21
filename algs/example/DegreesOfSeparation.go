// uses breadth-first search to find the degree of separation between two individuals in a social network
package main

import (
	"algs"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("usage: %s filename source", os.Args[0])
		return
	}

	sg := algs.NewSymbolGraph(os.Args[1], " ")
	g := sg.G()

	source := os.Args[2]
	if !sg.Contains(source) {
		fmt.Printf("%s not in database.", source)
		return
	}
	s := sg.Index(source)
	bfs := algs.NewBreadthFirstPaths(g, s)

	var sink string
	for {
		if _, err := fmt.Scanf("%s", &sink); err != nil {
			fmt.Println("fail to read stdin: ", err)
			return
		}
		if !sg.Contains(sink) {
			fmt.Printf("%s not in database.", sink)
			return
		}
		t := sg.Index(sink)
		if bfs.HasPathTo(t) {
			for _, v := range bfs.PathTo(t) {
				fmt.Println("    " + sg.Name(v))
			}
		}

	}
}

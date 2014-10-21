// Reads in a set of key-value pairs from a two-column CSV file
// specified on the command line; then, reads in keys from standard
// input and prints out corresponding values.
//
// % java LookupCSV amino.csv 0 3  
// input: TTA
// output: Leucine

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"algs"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("usage: %s <in> <key> <val>", os.Args[0])
		os.Exit(0)
	}

	keyField, _ := strconv.Atoi(os.Args[2])
	valField, _ := strconv.Atoi(os.Args[3])
	lines, err := algs.ReadLines(os.Args[1])
	if err != nil {
		fmt.Println("readlines fail", err)
		os.Exit(1)
	}

	var st *algs.LinearProbingHashST
	st = algs.NewLinearProbingHashST(10)
	for _, line := range lines {
        tokens := strings.Split(line[:len(line)-1], ",")    // strip '\n'
		key := tokens[keyField]
		val := tokens[valField]
        fmt.Println(key, val)
		st.Put(key, val)
	}

	r := bufio.NewReader(os.Stdin)
	for {
		s, err := r.ReadString('\n')
		if err != nil {
			break
		}
        s = s[:len(s)-1]
        if st.Contains(s) {
			fmt.Println(st.Get(s))
		} else {
			fmt.Println("Not found")
		}
	}
}

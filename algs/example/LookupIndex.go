package main

import (
    "bufio"
	"fmt"
	"os"
    "strings"

	"algs"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("usage: %s <filename> <separator>\n", os.Args[0])
		os.Exit(0)
	}

	separator := os.Args[2]
	var lines []string
	lines, _ = algs.ReadLines(os.Args[1])
	st := algs.NewLinearProbingHashST(10)
	ts := algs.NewLinearProbingHashST(10)
	for _, line := range lines {
		fields := strings.Split(line[:len(line)-1], separator)
        key := fields[0]
        for _, val := range(fields[1:])  {
            var q1 []string
            if !st.Contains(key) {
                q1 = []string {val}
            } else {
                q1 = st.Get(key).([]string)
                q1 = append(q1, val)
            }
            st.Put(key, q1)

            var q2 []string
            if !ts.Contains(val) {
                q2 = []string {key}
            } else {
                q2 = st.Get(key).([]string)
                q2 = append(q2, key)
            }
            ts.Put(val, q2)
        }
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
		}
		if ts.Contains(s) {
			fmt.Println(ts.Get(s))
		}
	}
}

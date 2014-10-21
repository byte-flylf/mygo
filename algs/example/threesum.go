/*  A program with cubic running time. Read in N integers
 *  and counts the number of triples that sum to exactly 0
 *  (ignoring integer overflow).
 *  Data files:   http://algs4.cs.princeton.edu/14analysis/1Kints.txt
 *                http://algs4.cs.princeton.edu/14analysis/2Kints.txt
 *                http://algs4.cs.princeton.edu/14analysis/4Kints.txt
 *                http://algs4.cs.princeton.edu/14analysis/8Kints.txt
 *                http://algs4.cs.princeton.edu/14analysis/16Kints.txt
 *                http://algs4.cs.princeton.edu/14analysis/32Kints.txt
 *                http://algs4.cs.princeton.edu/14analysis/1Mints.txt
 *
 */
package main

import (
	"algs"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <input>\n", os.Args[0])
		os.Exit(0)
	}
	a, err := algs.ReadInts(os.Args[1])
	if err != nil {
		log.Fatal("ReadInts: ", err)
	}
	timer := algs.NewStopwatch()
	var cnt int = algs.ThreeSum(a)
	fmt.Println("elapsed time = ", timer.ElapsedTime())
	fmt.Println(cnt)
}

//  Sort N random real numbers, T times using the two
//  algorithms specified on the command line.
// $> ./SortCompare Insertion Selection 1000 100
package main

import (
    "fmt"
    "os"
    "sort"
    "strconv"

    . "algs"
)

func timef(alg string,  slice sort.Float64Slice) float64 {
    sw := NewStopwatch()
    if alg == "Insertion" {
        InsertionSort(slice)
    } else if alg == "Selection" {
        SelectionSort(slice)
    } else if alg == "Shell" {
        ShellSort(slice)
    } else {
        fmt.Printf("invalid sort", alg)
        os.Exit(1)
    }
    return float64(sw.ElapsedTime())
}

func timeRandomInput(alg string, N int, T int) float64 {
    var total float64 = 0.0
    var slice sort.Float64Slice = make([]float64, N)
    for t := 0; t < T; t++ {
        for i := 0; i < N; i++ {
            slice[i] = Random();
        }
        total += timef(alg, slice)
    }
    return total
}

func main() {
    if len(os.Args) != 5 {
        fmt.Println("usage: sortcompare <sort1_name> <sort2_name> <n> <times>)")
        os.Exit(1)
    }
    N, _ := strconv.Atoi(os.Args[3])
    T, _ := strconv.Atoi(os.Args[4])
    time1 := timeRandomInput(os.Args[1], N, T)
    time2 := timeRandomInput(os.Args[2], N, T)
    fmt.Printf("For %d random doubles\n  %s is", N, os.Args[1])
    fmt.Printf(" %.4f times faster than %s\n", time2 / time1, os.Args[2])
}

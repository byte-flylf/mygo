// Filling a Slice Using Command-line Flags in Go
// http://lawlessguy.wordpress.com/2013/07/23/filling-a-slice-using-command-line-flags-in-go-golang/
package main

import (
    "flag"
    "fmt"
    "strconv"
)

type intslice []int

func (i *intslice) String() string {
    return fmt.Sprintf("%d", *i)
}

func (i *intslice) Set(value string) error {
    fmt.Printf("%s\n", value)
    tmp, err := strconv.Atoi(value)
    if err != nil {
        *i = append(*i, -1)
    } else {
        *i = append(*i, tmp)
    }
    return nil
}

var myints intslice

func main() {
    flag.Var(&myints, "i", "List of intergers")
    flag.Parse()
    if flag.NFlag() == 0 {
        flag.PrintDefaults()
    } else {
        fmt.Println("Here are the values in 'myints'")
        for i := 0; i < len(myints); i++ {
            fmt.Printf("%d\n", myints[i])
        }
    }
}

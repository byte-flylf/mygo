// Read in a list of words from input and print out
// the most frequently occurring word.
package main

import (
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"
    "unicode"
    . "algs"
)

//  Data files:   http://algs4.cs.princeton.edu/31elementary/tinyTale.txt
//                http://algs4.cs.princeton.edu/31elementary/tale.txt
//                http://algs4.cs.princeton.edu/31elementary/leipzig100K.txt
//                http://algs4.cs.princeton.edu/31elementary/leipzig300
//                http://algs4.cs.princeton.edu/31elementary/leipzig1M.txt
func main() {
    if len(os.Args) < 3 {
        fmt.Println("usage: frequency_counter <minlen> <input> [st_version]")
        os.Exit(1)
    }
    minlen,  err := strconv.Atoi(os.Args[1])
    if err != nil {
        log.Fatal("Atoi fail", err)
    }

    var st ST
    var lines []string
    st = NewBST()
    if len(os.Args) > 3 {
        if os.Args[3] == "bin" {
            fmt.Println("NewBinarySearchST")
            st = NewBinarySearchST()
        }
    } 
    lines, err = ReadLines(os.Args[2])
    if err != nil {
        log.Fatal("ReadLines fail", err)
    }
    var words []string
    for _, line := range lines {
        matches := strings.FieldsFunc(line, unicode.IsSpace)
        if len(matches) > 0 {
            words = append(words, matches...)
        }
    }
    fmt.Println("words", len(words))
    for _, word := range(words) {
        if len(word) < minlen {
            continue
        }
        if !st.Contains(word) {
            st.Put(word, 1)
        } else {
            cnt, _ := st.Get(word).(int)
            st.Put(word, cnt + 1)
        }
    }
    fmt.Println("size: ", st.Size())

    // find a key with the highest frequency count
    max := " "
    st.Put(max, 0)
    for _, x := range(st.Keys()) {
        v1, _ := st.Get(x).(int)
        v2, _ := st.Get(max).(int)
        if v1 > v2 {
            max, _ = x.(string)
        }
    }
    fmt.Println("max", max, st.Get(max))
}

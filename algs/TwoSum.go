package algs

import (
    "sort"
)

// O(n*n)
func TwoSum(a []int) int {
    n := len(a)
    cnt := 0
    for i:=0; i < n; i++ {
        for j := i+1; j < n; j++ {
            if a[i] + a[j] == 0 {
                cnt++
            }
        }
    }
    return cnt
}

// O(n*lgn)
func TwoSumFast(a []int) int {
    cnt := 0
    sort.Ints(a)
    n := len(a)
    for i:= 0; i < n; i++ {
        idx := sort.SearchInts(a, -a[i])
        if idx < len(a) && a[idx] == -a[i] && idx > i {
            cnt++
        }
    }
    return cnt
}

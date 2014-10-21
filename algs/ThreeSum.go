package algs

import (
    "sort"
)

// O(N^3)
func ThreeSum(a []int) int {
    n := len(a)
    cnt := 0
    for i := 0; i < n; i++ {
        for j := i+1; j < n; j++ {
            for k := j+1; k < n; k++ {
                if a[i] + a[j] + a[k] == 0 {
                    cnt++
                }
            }
        }
    }
    return cnt
}

// O(N^2*logN)
func ThreeSumFast(a []int) int {
    n := len(a)
    cnt := 0
    sort.Ints(a)
    for i := 0; i < n; i++ {
        for j := i+1; j < n; j++ {
            key := -a[i] - a[j]
            idx := sort.SearchInts(a, key)
            if idx < n && a[idx] == key && idx > j {
                cnt++
            }
        }
    }
    return cnt
}

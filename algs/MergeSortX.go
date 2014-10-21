package algs

import (
    "sort"
)

// merge sort Improvements
// 1. Use insertion sort for small subarrays.
// 2. Test whether array is already in order. 
// 3. Eliminate the copy to the auxiliary array.
func MergeXSort(a []string) {
    var aux []string
    aux = make([]string, len(a))
    copy(aux, a)
    mergexsort(aux, a, 0,  len(a)-1)
}

func mergex(src []string, dst []string, lo int, mid int, hi int) {
    i := lo
    j := mid + 1
    for k := lo; k <= hi; k++ {
        if i > mid {
            dst[k] = src[j]
            j++
        } else if j > hi {
            dst[k] = src[i]
            i++
        } else if src[j] < src[i] {
            dst[k] = src[j]
            j++
        } else {
            dst[k] = src[i]
            i++
        }
    }
}

const CUTOFF = 7

func mergexsort(src []string, dst []string, lo int, hi int) {
    // if (hi <= lo)  return
    if hi <= lo + CUTOFF {
        var a sort.StringSlice
        a = dst[lo:hi+1]
        InsertionSort(a)
        return
    }

    mid := lo + (hi - lo) / 2
    mergexsort(dst, src, lo, mid)
    mergexsort(dst, src, mid+1, hi)

    if src[mid+1] >= src[mid] {
        for i := lo; i <= hi; i++ {
            dst[i] = src[i]
        }
        return
    }

    mergex(src, dst, lo, mid, hi)
}

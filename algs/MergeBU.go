package algs


// Bottom-up mergesort
func MergeBU(a []string) {
    n := len(a)
    aux := make([]string, len(a))
    fmin := func (x, y int) int {
        if x < y {
            return x
        }
        return y
    }
    for sz := 1; sz < n; sz = sz+sz {
        for lo := 0; lo < n - sz; lo = lo + sz + sz {
            merge(a, aux, lo, lo+sz-1, fmin(lo+sz+sz-1, n-1))
        }
    }
}

package algs

func HeapSort(pq []string) {
    n := len(pq) - 1
    for k := n/2; k >= 1; k-- {
        sink(pq, k, n)
    }
    for n > 1 {
        // exch(pq, q, n--)
        pq[1], pq[n] = pq[n], pq[1]
        n--
        sink(pq, 1, n)
    }
}

func sink(pq []string, k, n int) {
    for 2*k <= n {
        j := 2*k
        if j < n && pq[j] < pq[j+1] {
            j++
        }
        if pq[k] > pq[j] {
            break
        }
        pq[k], pq[j] = pq[j], pq[k]
        k = j
    }
}

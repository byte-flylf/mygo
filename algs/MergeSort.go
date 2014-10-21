package algs

// Top-down mergesort.
func MergeSort(a []string) {
	var aux []string
	aux = make([]string, len(a))
	mergesort(a, aux, 0, len(a)-1)
}

// mergesort a[lo..hi] using auxiliary array aux[lo..hi]
func mergesort(a []string, aux []string, lo int, hi int) {
	if hi <= lo {
		return
	}
	mid := lo + (hi-lo)/2
	mergesort(a, aux, lo, mid)
	mergesort(a, aux, mid+1, hi)
	merge(a, aux, lo, mid, hi)
}

//  stably merge a[lo .. mid] with a[mid+1 ..hi] using aux[lo .. hi]
func merge(a []string, aux []string, lo, mid, hi int) {
	for k := lo; k <= hi; k++ {
		aux[k] = a[k] //  if type 'a' is sort.Interface,  don't support this operation!!!!
	}

	// merge back to a[]
	i := lo
	j := mid + 1
	for k := lo; k <= hi; k++ {
		if i > mid {
			a[k] = aux[j]
			j++
		} else if j > hi {
			a[k] = aux[i]
			i++
		} else if aux[j] < aux[i] {
			a[k] = aux[j]
			j++
		} else {
			a[k] = aux[i]
			i++
		}
	}
}

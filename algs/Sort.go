package algs

import (
	"sort"
)

// Proposition. Selection sort uses ~N2/2 compares and N exchanges to sort an array of length N
func SelectionSort(a sort.Interface) {
	n := a.Len()
	for i := 0; i < n; i++ {
		min := i
		for j := i + 1; j < n; j++ {
			if a.Less(j, min) {
				min = j
			}
		}
		a.Swap(i, min)
	}
}

// Proposition. For randomly ordered arrays of length N with with distinct keys,
// insertion sort uses ~N2/4 compares and ~N2/4 exchanges on the average.
// The worst case is ~ N2/2 compares and ~ N2/2 exchanges and the best case is N-1 compares and 0 exchanges.
func InsertionSort(a sort.Interface) {
	n := a.Len()
	for i := 1; i < n; i++ {
		for j := i; j > 0 && a.Less(j, j-1); j-- {
			a.Swap(j, j-1)
		}
	}
}

// Shell sort
func ShellSort(a sort.Interface) {
	n := a.Len()
	h := 1
	// 3x+1 increment sequence:  1, 4, 13, 40, 121, 364, 1093, ...
	for h < n/3 {
		h = 3*h + 1
	}

	for h >= 1 {
		for i := h; i < n; i++ {
			for j := i; j >= h && a.Less(j, j-h); j -= h {
				a.Swap(j, j-h)
			}
		}
		h /= 3
	}
}

// partition the subarray a[lo .. hi] by returning an index j
// so that a[lo .. j-1] <= a[j] <= a[j+1 .. hi]
func partition(a sort.Interface, lo, hi int) int {
	i, j := lo+1, hi
	// v = a[lo]
	for {
		for a.Less(i, lo) {
			if i == hi {
				break
			}
			i++
		}
		for a.Less(lo, j) {
			if j == lo {
				break
			}
			j--
		}
		if i >= j {
			break
		}
		a.Swap(i, j)
		i++
		j--
	}
	// put v = a[j] into position
	a.Swap(lo, j)
	// with a[lo .. j-1] <= a[j] <= a[j+1 .. hi]
	return j
}

func quicksort(a sort.Interface, lo, hi int) {
	if hi <= lo {
		return
	}
	j := partition(a, lo, hi)
	quicksort(a, lo, j-1)
	quicksort(a, j+1, hi)
}

func QuickSort(a sort.Interface) {
	quicksort(a, 0, a.Len()-1)
}

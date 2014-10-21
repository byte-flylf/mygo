package algs

// precondition: array is sorted
func BinarySearch(key int, arr []int) int {
	var lo int = 0
	var hi int = len(arr) - 1
	var mid int
	for lo <= hi {
		mid = lo + (hi-lo)/2
        if key == arr[mid] {
			return mid
		} else if key < arr[mid] {
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	return -1
}

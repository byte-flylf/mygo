package algs

// Connected components
type CC interface {
	Connected(v, w int) bool // are v and w connected ?
	Count() int              // number of connected components
	ID(v int) int            // component identifier for v
}

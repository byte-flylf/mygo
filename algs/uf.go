package algs

// Union-Find API
type UF interface {
    Union(p int, q int)   // instantiate N isolated components 0 through N-1
    Find(p int) int // Return component identifier for component containing p
    Connected(p int, q int) bool   // are elements p and q in the same component
    Count() int         // return number of connected components
}

// implement
// QuickFindUF.go  QuickUnionUF.go  WeightedQuickUnionUF.go

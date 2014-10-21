package algs

// whether there exists a path between two given vertices but to find such a path (if one exists).
type Pather interface {
    HasPathTo(v int) bool   // is there a path from s to v?
    PathTo(v int) []int  // path from s to v, nil if no such path
}

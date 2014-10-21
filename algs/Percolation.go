package algs

import (
	"fmt"
)

/*
* To use this class, instantiate a new Percolation square grid with NxN sides
 * Use the Percolation.open method to open a certain number of grid sites
 * and then test whether or not the grid will percolate with the
 * Percolation.percolates method.
*/
type Percolation struct {
	// reference indices with "i" and "j", where i and j are between 1 and N
	// (1, 1) is the upper left side of the grid
	siteSize int      // The size of one side of the grid (N)
	grid     UF       // Union find data structure representing the grid
	sites    [][]bool // data structure for tracking open/closed sites
	gridSize int      // Size of the Union Find collection
}

// Create N-by-N grid, with all sites blocked.
func NewPercolation(N int) *Percolation {
	p := new(Percolation)
	p.siteSize = N
	p.gridSize = N*N + 2
	p.sites = make([][]bool, N)
	for i := 0; i < N; i++ {
		p.sites[i] = make([]bool, N)
	}
	p.grid = NewWeightedQuickUnionUF(p.gridSize)
	for i := 1; i <= N; i++ {
		p.grid.Union(0, i)                         // virtual top site,  0
		p.grid.Union(p.gridSize-1, p.gridSize-i-1) // virtual bottom site, N*N+1
	}

	return p
}

// Open site (row i, column j) if it is not already open.
func (p *Percolation) Open(i, j int) {
	p.validateArgs(i, j)

	if !p.getSite(i, j) {
		p.setSite(i, j, true)

		// Union to all open sites around us
		if i != 1 && p.getSite(i-1, j) { // We are not at the top
			p.grid.Union(p.getIndexFromArgs(i, j), p.getIndexFromArgs(i-1, j))
		}
		if i != p.siteSize && p.getSite(i+1, j) { // Not at bottom
			p.grid.Union(p.getIndexFromArgs(i, j), p.getIndexFromArgs(i+1, j))
		}
		if j != 1 && p.getSite(i, j-1) { // Not on the left side
			p.grid.Union(p.getIndexFromArgs(i, j), p.getIndexFromArgs(i, j-1))
		}
		if j != p.siteSize && p.getSite(i, j+1) { // Not on the right
			p.grid.Union(p.getIndexFromArgs(i, j), p.getIndexFromArgs(i, j+1))
		}
	}
}

// Returns true if the site (row i, colum j) is open, false otherwise
func (p *Percolation) IsOpen(i, j int) bool {
	p.validateArgs(i, j)
	return p.getSite(i, j)
}

// Returns true if the site (row i, colum j) is open and full
// where "full" means the site is connected to the top row.
// Return false otherwise
func (p *Percolation) IsFull(i, j int) bool {
	p.validateArgs(i, j)
	return p.IsOpen(i, j) && p.grid.Connected(0, p.getIndexFromArgs(i, j))
}

// Returns true if the system percolates from top to bottom, false otherwise.
func (p *Percolation) Percolates() bool {
	return p.grid.Connected(0, p.gridSize-1)
}

func (p *Percolation) validateArgs(i, j int) {
	if i > p.siteSize || i < 1 || j > p.siteSize || j < 1 {
		message := fmt.Sprintf("Indices [%d, %d] are outside bounds [1, %d]", i, j, p.siteSize)
		panic(message)
	}
}

func (p *Percolation) getSite(i, j int) bool {
	return p.sites[i-1][j-1]
}

func (p *Percolation) setSite(i, j int, val bool) {
	p.sites[i-1][j-1] = val
}

func (p *Percolation) getIndexFromArgs(i, j int) int {
	return p.siteSize*(i-1) + j
}

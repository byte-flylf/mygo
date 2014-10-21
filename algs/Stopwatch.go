package algs

import (
    "time"
)

// data type for measuring the running time (wall clock) of a program.
type  Stopwatch struct {
    start time.Time
}

// Create a stopwatch object.
func NewStopwatch() *Stopwatch {
    return &Stopwatch{start: time.Now()}
}

// Return elapsed time (in seconds) since this object was created.
func (self *Stopwatch) ElapsedTime() time.Duration {
    now := time.Now()
    return now.Sub(self.start)
}

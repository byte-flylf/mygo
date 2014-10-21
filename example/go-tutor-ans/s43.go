package main
 
import (
	"fmt"
	"math"
)
 
func Sqrt(x float64) float64 {
	z, oldz, delta := 1.0, 0.0, 1e-10
	
	for math.Abs(z-oldz) > delta {
		oldz, z = z, z - ( (z*z - x) / (2*z) )
	}
	
	return z
}
 
func main() {
	for i := 1.0; i < 10; i++ {
		fmt.Printf("%g: %g %g\n", i, Sqrt(i), math.Sqrt(i))
	}
}

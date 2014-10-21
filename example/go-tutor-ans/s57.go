package main
 
import (
	"fmt"
	"math"
)
 
type ErrNegativeSqrt float64
 
func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %g", float64(e))
}
 
func Sqrt(f float64) (float64, error) {
	if f < 0 {
		return 0, ErrNegativeSqrt(-2)
	}
	
	z, oldz, delta := 1.0, 0.0, 1e-10
	
	for math.Abs(z-oldz) > delta {
		oldz, z = z, z - ( (z*z - f) / (2*z) )
	}
	
	return z, nil
}
 
func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}

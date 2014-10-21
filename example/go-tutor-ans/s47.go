package main
 
import (
	"fmt"
	"math"
)
 
func Cbrt(x complex128) complex128 {
	z, oldz := complex128(1), complex128(0)
	delta := 1e-20
	
	for math.Abs(real(z)-real(oldz)) > delta {
		oldz, z = z, z - (((z*z*z) - x) / (3*(z*z)))
	}
	
	return z
}
 
func main() {
	fmt.Println(Cbrt(2))
	fmt.Println(math.Pow(2, 1/3.0))
}

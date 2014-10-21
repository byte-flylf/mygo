package main
 
import "code.google.com/p/go-tour/pic"
 
func Pic(dx, dy int) [][]uint8 {
	img := make([][]uint8, dy)
	for i := 0; i < dy; i++ {
		img[i] = make([]uint8, dx)
 
		for j := 0; j < dx; j++ {
			// Awesome!!
			img[i][j] = uint8(i^j)
			// Gradient
			//img[i][j] = uint8((i+j)/2)
			// Whoooa
			//img[i][j] = uint8(i*j)
		}
	}
	
	return img
}
 
func main() {
	pic.Show(Pic)
}

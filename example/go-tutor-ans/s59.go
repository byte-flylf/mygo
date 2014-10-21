package main
 
import (
	"image"
	"image/color"
	"code.google.com/p/go-tour/pic"
)
 
type Image struct{
	Rect image.Rectangle
	Pixels [][]color.Color
}
 
func NewImage(w, h int) *Image {
	pixels := make([][]color.Color, h)
	
	for i := 0; i < h; i++ {
		pixels[i] = make([]color.Color, w)
		
		for j := 0; j < w; j++ {
			pixels[i][j] = color.RGBA{uint8(255-i^j), uint8(255-i^j), uint8(((h-i)*(w-j))/(1<<5)), uint8(255)}
		}
	}
	
	return &Image{image.Rect(0, 0, w, h), pixels}
}
 
func (i *Image) ColorModel() color.Model {
	return color.RGBAModel
}
 
func (i *Image) Bounds() image.Rectangle {
	return i.Rect
}
 
func (i *Image) At(x, y int) color.Color {
	return i.Pixels[y][x]
}
 
func main() {
	m := NewImage(256, 256)
	pic.ShowImage(m)
}

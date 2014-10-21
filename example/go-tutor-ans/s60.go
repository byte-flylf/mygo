package main
 
import (
	"io"
	"os"
	"strings"
)
 
type rot13Reader struct {
	r io.Reader
}
 
func (r *rot13Reader) Read(p []byte) (n int, err error) {
	rlen, err := r.r.Read(p)
	for i := 0; i < rlen; i++ {
		if ('A' < p[i] || 'a' < p[i]) && (p[i] <= 'M' || p[i] <= 'm') {
			p[i] = p[i]+13
		} else if ('M' < p[i] || 'm' < p[i]) && (p[i] <= 'Z' || p[i] <= 'z') {
			p[i] = p[i]-13
		}
	}
	return rlen, err
}
 
func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}

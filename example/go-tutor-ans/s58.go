package main
 
import (
	"fmt"
	"bytes"
	"net/http"
)
 
type String string
 
type Struct struct {
	Greeting string
	Punct    string
	Who      string
}
 
func (s String) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	msg := string(s)
	b := bytes.NewBufferString(msg)
	w.Write(b.Bytes())
}
 
func (s *Struct) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("%s%s %s", s.Greeting, s.Punct, s.Who)
	b := bytes.NewBufferString(msg)
	w.Write(b.Bytes())
}
 
func main() {
	http.Handle("/string", String("I'm a frayed knot."))
	http.Handle("/struct", &Struct{"Hello", ":", "Gophers!"})
	http.ListenAndServe("localhost:4000", nil)
}

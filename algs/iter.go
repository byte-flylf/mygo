package algs

// Iterator in Go
func iterSlice(slice []interface{}) (chan interface{}) {
    ch := make(chan interface{})
    go func() {
        for _, val := range slice {
            ch <- val
        }
        close(ch)
    } ()
    return ch
}

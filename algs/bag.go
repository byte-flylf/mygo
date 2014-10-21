package algs

// bag is a collection where removing items is not supported
// its purpose is to provide clients with the ability to collect items and then to iterate through the collected items.
type Bager interface {
	Add(item interface{})
	IsEmpty() bool
	Size() int
    ChannelIterator() (chan interface{})
}

package consumer

import (
	"container/list"
)

var makes int

type queued struct {
	slice []byte
}

func makeRecycler(n int) (get, give chan []byte) {
	get = make(chan []byte)
	give = make(chan []byte)

	go func() {
		q := new(list.List)
		for {
			if q.Len() == 0 {
				q.PushFront(queued{slice: makeBuffer(n)})
			}

			e := q.Front()
			select {
			case b := <-give:
				q.PushFront(queued{slice: b})

			case get <- e.Value.(queued).slice:
				q.Remove(e)
			}
		}

	}()
	return
}

func makeBuffer(n int) []byte {
	makes += 1
	return make([]byte, n)
}

// container/heap/#pkg-examples
package main

import (
    "container/heap"
    "fmt"
)

type Item struct {
    value    string
    priority int
    index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
    return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
    pq[i].index = i
    pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
    n := len(*pq)
    item := x.(*Item)
    item.index = n
    *pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    item := old[n-1]
    item.index = -1
    *pq = old[0 : n-1]
    return item
}

func (pq *PriorityQueue) update(item *Item, value string, priority int) {
    heap.Remove(pq, item.index)
    item.value = value
    item.priority = priority
    heap.Push(pq, item)
}

func main() {
    items := map[string]int{
        "banana": 3, "apple": 2, "pear": 4,
    }

    pq := &PriorityQueue{}
    heap.Init(pq)
    for value, priority := range items {
        item := &Item{value: value, priority: priority}
        heap.Push(pq, item)
    }

    item := &Item{value: "orange", priority: 1}
    heap.Push(pq, item)
    pq.update(item, item.value, 5)
    for pq.Len() > 0 {
        item := heap.Pop(pq).(*Item)
        fmt.Printf("%.2d:%s ", item.priority, item.value)
    }
}

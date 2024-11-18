package utils

import (
	"container/heap"
)

type Element struct {
	value    *Node
	priority int
	index    int
}

type PriorityQueue []*Element

func (pq PriorityQueue) Len() int { return len(pq) }

/* Min heap */
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	elem := x.(*Element)
	elem.index = n
	*pq = append(*pq, elem)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	elem := old[n-1]
	old[n-1] = nil
	elem.index = -1
	*pq = old[0 : n-1]
	return elem
}

func (pq *PriorityQueue) update(elem *Element, value *Node, priority int) {
	elem.value = value
	elem.priority = priority
	heap.Fix(pq, elem.index)
}

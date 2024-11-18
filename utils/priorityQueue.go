package utils

import (
	"container/heap"
)

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

/* Min heap */
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].QueueIndex = i
	pq[j].QueueIndex = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	node := x.(*Node)
	node.QueueIndex = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.QueueIndex = -1
	*pq = old[0 : n-1]
	return node
}

func (pq *PriorityQueue) update(node *Node, priority int) {
	node.Priority = priority
	heap.Fix(pq, node.QueueIndex)
}

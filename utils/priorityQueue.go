/* Based on: https://pkg.go.dev/container/heap#example-package-PriorityQueue */
package utils

import (
	"container/heap"
)

/* Implements build-in heap interface */
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

/* Max heap */
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].QueuePriority > pq[j].QueuePriority
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

func (pq *PriorityQueue) update(node *Node, cost int) {
	node.Cost = cost
	if IsHeuristicOn {
		node.QueuePriority = node.Cost
		heap.Fix(pq, node.QueueIndex)
	}
}

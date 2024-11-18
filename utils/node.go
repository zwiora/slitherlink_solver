package utils

import (
	"math"
)

type Node struct {
	Neighbours   []*Node
	Value        int8
	IsInLoop     bool
	IsVisited    bool
	CanBeRemoved bool
	Priority     int
	QueueIndex   int
}

func (n *Node) GetDegree() int {
	count := 0
	for _, v := range n.Neighbours {
		if v != nil && v.IsInLoop {
			count++
		}
	}
	return count
}

func (n *Node) GetLinesAround(maxNeighbourCount int) int {
	if n.IsInLoop {
		return maxNeighbourCount - n.GetDegree()
	}
	return n.GetDegree()
}

func (n *Node) GetCostOfField(maxNeighbourCount int) int {
	linesAround := n.GetLinesAround(int(maxNeighbourCount))
	return int(math.Abs(float64(linesAround) - float64(n.Value)))
}

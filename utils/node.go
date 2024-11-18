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
	Cost         int
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

func (n *Node) GetLinesAround(maxDegree int) int {
	if n.IsInLoop {
		return maxDegree - n.GetDegree()
	}
	return n.GetDegree()
}

func (n *Node) GetCostOfField(maxDegree int) int {
	linesAround := n.GetLinesAround(int(maxDegree))
	return int(math.Abs(float64(linesAround) - float64(n.Value)))
}

func (n *Node) CalculateNodeCost(g *Graph) int {
	newCost := 0

	if n.Value != -1 {
		newCost += n.GetCostOfField(int(g.MaxDegree))
	}

	for _, v := range n.Neighbours {
		if v != nil && v.Value != -1 {
			newCost += v.GetCostOfField(int(g.MaxDegree))
		}
	}

	n.IsInLoop = false

	if n.Value != -1 {
		newCost -= n.GetCostOfField(int(g.MaxDegree))
	}

	for _, v := range n.Neighbours {
		if v != nil && v.Value != -1 {
			newCost -= v.GetCostOfField(int(g.MaxDegree))
		}
	}

	n.IsInLoop = true

	return newCost
}

func (n *Node) UpdateNodeCost(g *Graph) {
	newCost := n.CalculateNodeCost(g)

	if n.Cost != newCost {
		g.AvailableMoves.update(n, newCost)
	}
}

package utils

import (
	"math"
)

type Node struct {
	Neighbours    []*Node
	Value         int8
	IsInLoop      bool
	IsVisited     bool
	IsForRemoval  bool
	CanBeRemoved  bool
	Cost          int
	QueueIndex    int
	QueuePriority int
}

/* Calculates number of neighbours that are in the loop */
func (n *Node) GetDegree() int {
	count := 0
	for _, v := range n.Neighbours {
		if v != nil && v.IsInLoop {
			count++
		}
	}
	return count
}

/* Calculates number of edges of the loop around the field */
func (n *Node) getLinesAround(maxDegree int) int {
	if n.IsInLoop {
		return maxDegree - n.GetDegree()
	}
	return n.GetDegree()
}

/* Calculates cost of the single field (node) */
func (n *Node) getCostOfField(maxDegree int) int {
	linesAround := n.getLinesAround(int(maxDegree))
	return int(math.Abs(float64(linesAround) - float64(n.Value)))
}

/* Calculates new cost on the node for priority queue */
func (n *Node) calculateNodeCost(g *Graph) int {
	newCost := 0

	if n.Value != -1 {
		newCost += n.getCostOfField(int(g.MaxDegree))
	}

	for _, v := range n.Neighbours {
		if v != nil && v.Value != -1 {
			newCost += v.getCostOfField(int(g.MaxDegree))
		}
	}

	n.IsInLoop = false

	if n.Value != -1 {
		newCost -= n.getCostOfField(int(g.MaxDegree))
	}

	for _, v := range n.Neighbours {
		if v != nil && v.Value != -1 {
			newCost -= v.getCostOfField(int(g.MaxDegree))
		}
	}

	n.IsInLoop = true

	return newCost
}

func (n *Node) SetNodeCost(g *Graph) {
	newCost := n.calculateNodeCost(g)
	n.Cost = newCost

	if IsHeuristicOn {
		n.QueuePriority = newCost
	}
}

/* Calculates and updates cost of the node in the heap */
func (n *Node) UpdateNodeCost(g *Graph) {
	newCost := n.calculateNodeCost(g)
	n.Cost = newCost

	if IsHeuristicOn {
		if n.Cost != n.QueuePriority {
			g.AvailableMoves.update(n, newCost)
		}
	}
}

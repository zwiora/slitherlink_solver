package utils

import (
	"math"
)

type Node struct {
	Neighbours    []*Node
	NextRow       *Node
	Value         int8
	IsInLoop      bool
	IsVisited     bool
	IsDecided     bool
	IsForRemoval  bool
	CanBeRemoved  bool
	Cost          int
	QueueIndex    int
	QueuePriority int
	TemplateGroup *List
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

func (n *Node) IsDeletionBreakingSecondRule() bool {
	oldIsInLoop := n.IsInLoop

	/* Checking if it would have enough edges */
	n.IsInLoop = false
	if n.GetDegree() < int(n.Value) {
		n.IsInLoop = oldIsInLoop
		return true
	}

	/* Checking if neighbour would have enough edges*/
	for _, v := range n.Neighbours {
		if v != nil && !v.IsInLoop {
			if v.GetDegree() < int(v.Value) {
				n.IsInLoop = oldIsInLoop
				return true
			}
		}
	}

	n.IsInLoop = oldIsInLoop

	return false
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
		if n.IsForRemoval {
			n.QueuePriority = 1000
		} else if n.TemplateGroup != nil {
			n.QueuePriority = newCost + 100
		} else {
			n.QueuePriority = newCost
		}
	}
}

/* Calculates and updates cost of the node in the heap */
func (n *Node) UpdateNodeCost(g *Graph) {
	newCost := n.calculateNodeCost(g)
	n.Cost = newCost

	if IsHeuristicOn {
		if n.Cost != n.QueuePriority {
			if n.IsForRemoval {
				g.AvailableMoves.update(n, 10000)
			} else if n.TemplateGroup != nil {
				if HeuristicType == 1 {
					g.AvailableMoves.update(n, newCost+1000)
				} else if HeuristicType == 2 {
					groupSize := 0
					if n.TemplateGroup != nil {
						groupSize += n.TemplateGroup.Length

						if n.TemplateGroup.OppositeList != nil {
							groupSize += n.TemplateGroup.OppositeList.Length
						}
					}
					g.AvailableMoves.update(n, 1000+groupSize)
				} else if HeuristicType == 3 {
					groupSize := 0
					if n.TemplateGroup != nil {
						groupSize += n.TemplateGroup.Length

						if n.TemplateGroup.OppositeList != nil {
							groupSize += n.TemplateGroup.OppositeList.Length
						}
					}
					g.AvailableMoves.update(n, 1000+groupSize*10+newCost)
				} else {
					g.AvailableMoves.update(n, newCost)
				}
			} else {
				g.AvailableMoves.update(n, newCost)
			}
		}
	}
}

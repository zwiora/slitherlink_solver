package utils

import "github.com/golang-collections/collections/queue"

func isTheSameState(n1 *Node, n2 *Node) bool {
	if ((n1 == nil || n1.IsForRemoval) && (n2 == nil || n2.IsForRemoval)) || ((n1 != nil && n1.IsVisited) && (n2 != nil && n2.IsVisited)) {
		return true
	}
	return false
}

func addNeighboursToQueue(n *Node, q *queue.Queue, excludedNodes ...*Node) {
	for i := 0; i < len(n.Neighbours); i++ {
		thisNeighbour := n.Neighbours[i]
		if thisNeighbour != nil {
			isExcluded := false
			for k, v := range excludedNodes {
				if thisNeighbour == v {
					isExcluded = true
					excludedNodes = append(excludedNodes[:k], excludedNodes[k+1:]...)
					break
				}
			}
			if !isExcluded {
				q.Enqueue(n.Neighbours[i])
			}
		}
	}
}

func (n *Node) findZeroTemplates(g *Graph, q *queue.Queue) bool {
	if n.Value == 0 {
		for i := 0; i < len(n.Neighbours); i++ {
			thisNeighbour := n.Neighbours[i]
			if thisNeighbour == nil || thisNeighbour.IsForRemoval || thisNeighbour.IsVisited {
				if thisNeighbour != nil && thisNeighbour.IsVisited {
					n.IsVisited = true
					for j := 0; j < len(n.Neighbours); j++ {
						if j != i && n.Neighbours[j] != nil {
							n.Neighbours[j].IsVisited = true
							addNeighboursToQueue(n.Neighbours[j], q, n)
						}
					}
				} else {
					n.IsForRemoval = true
					for j := 0; j < len(n.Neighbours); j++ {
						if j != i && n.Neighbours[j] != nil {
							n.Neighbours[j].IsForRemoval = true
							addNeighboursToQueue(n.Neighbours[j], q, n)
						}
					}
				}
				return true
			}
		}
	}
	return false
}

/* Returns true, if template found */
func (n *Node) findCornerTemplates(g *Graph, q *queue.Queue) bool {
	for i := 0; i < len(n.Neighbours); i++ {
		thisNeighbour := n.Neighbours[i]
		nextNeighbour := n.Neighbours[(i+1)%int(g.MaxDegree)]
		if isTheSameState(thisNeighbour, nextNeighbour) {
			switch n.Value {
			case 0:
				if thisNeighbour == nil || thisNeighbour.IsForRemoval {
					n.IsForRemoval = true

					oppositeNeighbour := n.Neighbours[(i+2)%int(g.MaxDegree)]
					if oppositeNeighbour != nil {
						oppositeNeighbour.IsForRemoval = true
						addNeighboursToQueue(oppositeNeighbour, q, n)
					}

					oppositeNeighbour = n.Neighbours[(i+3)%int(g.MaxDegree)]
					if oppositeNeighbour != nil {
						oppositeNeighbour.IsForRemoval = true
						addNeighboursToQueue(oppositeNeighbour, q, n)
					}
				} else {
					n.IsVisited = true

					oppositeNeighbour := n.Neighbours[(i+2)%int(g.MaxDegree)]
					if oppositeNeighbour != nil {
						oppositeNeighbour.IsVisited = true
						addNeighboursToQueue(oppositeNeighbour, q, n)
					}

					oppositeNeighbour = n.Neighbours[(i+3)%int(g.MaxDegree)]
					if oppositeNeighbour != nil {
						oppositeNeighbour.IsVisited = true
						addNeighboursToQueue(oppositeNeighbour, q, n)
					}
				}
			case 1:
				if thisNeighbour == nil || thisNeighbour.IsForRemoval {
					n.IsForRemoval = true
				} else {
					n.IsVisited = true
				}
				addNeighboursToQueue(n, q)
			case 2:
				if thisNeighbour == nil || thisNeighbour.IsForRemoval {
					oppositeNeighbour := n.Neighbours[(i+2)%int(g.MaxDegree)]
					if oppositeNeighbour != nil {
						oppositeNeighbour.IsVisited = true
						addNeighboursToQueue(oppositeNeighbour, q, n)
					}

					oppositeNeighbour = n.Neighbours[(i+3)%int(g.MaxDegree)]
					if oppositeNeighbour != nil {
						oppositeNeighbour.IsVisited = true
						addNeighboursToQueue(oppositeNeighbour, q, n)
					}

					if thisNeighbour != nil {
						diagonal := thisNeighbour.Neighbours[(i+1)%int(g.MaxDegree)]
						if diagonal != nil && diagonal.IsVisited {
							n.IsForRemoval = true
						}
					}
				} else {
					oppositeNeighbour := n.Neighbours[(i+2)%int(g.MaxDegree)]
					if oppositeNeighbour != nil {
						oppositeNeighbour.IsForRemoval = true
						addNeighboursToQueue(oppositeNeighbour, q, n)
					}

					oppositeNeighbour = n.Neighbours[(i+3)%int(g.MaxDegree)]
					if oppositeNeighbour != nil {
						oppositeNeighbour.IsForRemoval = true
						addNeighboursToQueue(oppositeNeighbour, q, n)
					}

					/* diagonal might determine this node */
					if thisNeighbour != nil {
						diagonal := thisNeighbour.Neighbours[(i+1)%int(g.MaxDegree)]
						if diagonal == nil || diagonal.IsForRemoval {
							n.IsVisited = true
						}
					}
				}
			case 3:
				if thisNeighbour == nil || thisNeighbour.IsForRemoval {
					n.IsVisited = true
					/* We can determine state of diagonal node */
					if thisNeighbour != nil {
						diagonal := thisNeighbour.Neighbours[(i+1)%int(g.MaxDegree)]
						if diagonal != nil {
							diagonal.IsForRemoval = true
							addNeighboursToQueue(diagonal, q, thisNeighbour, nextNeighbour)
						}
					}
				} else {
					n.IsForRemoval = true
					/* We can determine state of diagonal node */
					if thisNeighbour != nil {
						diagonal := thisNeighbour.Neighbours[(i+1)%int(g.MaxDegree)]
						if diagonal != nil {
							diagonal.IsVisited = true
							addNeighboursToQueue(diagonal, q, thisNeighbour, nextNeighbour)
						}
					}
				}
				addNeighboursToQueue(n, q)
			}

			return true
		}
	}

	return false
}

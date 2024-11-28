package utils

import "github.com/golang-collections/collections/queue"

func isTheSameState(n1 *Node, n2 *Node) (bool, bool) {

	// we can't say if they have the same state
	if (n1 != nil && !n1.IsDecided) || (n2 != nil && !n2.IsDecided) {
		return false, false
	}

	isN1Out := n1 == nil || (n1.IsDecided && n1.IsForRemoval)
	isN2Out := n2 == nil || (n2.IsDecided && n2.IsForRemoval)

	if isN1Out == isN2Out {
		return true, isN1Out
	}

	return false, false

	// if ((n1 == nil || n1.IsForRemoval) && (n2 == nil || n2.IsForRemoval)) || ((n1 != nil && n1.IsVisited) && (n2 != nil && n2.IsVisited)) {
	// 	return true
	// }
	// return false
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

func isNodeDecided(n *Node) bool {
	return n == nil || n.IsDecided
}

func isNodeDecidedOut(n *Node) bool {
	return !(n != nil && n.IsDecided && !n.IsForRemoval)
}

func (n *Node) findZeroTemplates(g *Graph, q *queue.Queue) bool {
	if n.Value == 0 {
		for i := 0; i < len(n.Neighbours); i++ {
			thisNeighbour := n.Neighbours[i]
			if isNodeDecided(thisNeighbour) {
				n.IsDecided = true
				n.IsForRemoval = isNodeDecidedOut(thisNeighbour)
				for j := 0; j < len(n.Neighbours); j++ {
					if j != i && n.Neighbours[j] != nil {
						n.Neighbours[j].IsDecided = true
						n.Neighbours[j].IsForRemoval = n.IsForRemoval
						addNeighboursToQueue(n.Neighbours[j], q, n)
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
		isSame, _ := isTheSameState(thisNeighbour, nextNeighbour)
		if isSame {
			switch n.Value {
			case 0:
				if thisNeighbour == nil || thisNeighbour.IsForRemoval {
					n.IsForRemoval = true
					n.IsDecided = true

					oppositeNeighbour := n.Neighbours[(i+2)%int(g.MaxDegree)]
					if oppositeNeighbour != nil {
						oppositeNeighbour.IsForRemoval = true
						oppositeNeighbour.IsDecided = true
						addNeighboursToQueue(oppositeNeighbour, q, n)
					}

					oppositeNeighbour = n.Neighbours[(i+3)%int(g.MaxDegree)]
					if oppositeNeighbour != nil {
						oppositeNeighbour.IsForRemoval = true
						oppositeNeighbour.IsDecided = true
						addNeighboursToQueue(oppositeNeighbour, q, n)
					}
				} else {
					n.IsForRemoval = false
					n.IsDecided = true

					oppositeNeighbour := n.Neighbours[(i+2)%int(g.MaxDegree)]
					if oppositeNeighbour != nil {
						oppositeNeighbour.IsForRemoval = false
						oppositeNeighbour.IsDecided = true
						addNeighboursToQueue(oppositeNeighbour, q, n)
					}

					oppositeNeighbour = n.Neighbours[(i+3)%int(g.MaxDegree)]
					if oppositeNeighbour != nil {
						oppositeNeighbour.IsForRemoval = false
						oppositeNeighbour.IsDecided = true
						addNeighboursToQueue(oppositeNeighbour, q, n)
					}
				}
			case 1:
				if thisNeighbour == nil || thisNeighbour.IsForRemoval {
					n.IsForRemoval = true
				} else {
					n.IsForRemoval = false
				}
				n.IsDecided = true
				addNeighboursToQueue(n, q)
			case 2:
				if thisNeighbour == nil || thisNeighbour.IsForRemoval {
					oppositeNeighbour := n.Neighbours[(i+2)%int(g.MaxDegree)]
					if oppositeNeighbour != nil {
						oppositeNeighbour.IsForRemoval = false
						oppositeNeighbour.IsDecided = true
						addNeighboursToQueue(oppositeNeighbour, q, n)
					}

					oppositeNeighbour = n.Neighbours[(i+3)%int(g.MaxDegree)]
					if oppositeNeighbour != nil {
						oppositeNeighbour.IsForRemoval = false
						oppositeNeighbour.IsDecided = true
						addNeighboursToQueue(oppositeNeighbour, q, n)
					}

					if thisNeighbour != nil {
						diagonal := thisNeighbour.Neighbours[(i+1)%int(g.MaxDegree)]
						if diagonal != nil && diagonal.IsDecided && !diagonal.IsForRemoval {
							n.IsForRemoval = true
							n.IsDecided = true
						}
					}
				} else {
					oppositeNeighbour := n.Neighbours[(i+2)%int(g.MaxDegree)]
					if oppositeNeighbour != nil {
						oppositeNeighbour.IsForRemoval = true
						oppositeNeighbour.IsDecided = true
						addNeighboursToQueue(oppositeNeighbour, q, n)
					}

					oppositeNeighbour = n.Neighbours[(i+3)%int(g.MaxDegree)]
					if oppositeNeighbour != nil {
						oppositeNeighbour.IsForRemoval = true
						oppositeNeighbour.IsDecided = true
						addNeighboursToQueue(oppositeNeighbour, q, n)
					}

					/* diagonal might determine this node */
					if thisNeighbour != nil {
						diagonal := thisNeighbour.Neighbours[(i+1)%int(g.MaxDegree)]
						if diagonal == nil || diagonal.IsForRemoval {
							n.IsForRemoval = false
							n.IsDecided = true
						}
					}
				}
			case 3:
				if thisNeighbour == nil || thisNeighbour.IsForRemoval {
					n.IsForRemoval = false
					n.IsDecided = true
					/* We can determine state of diagonal node */
					if thisNeighbour != nil {
						diagonal := thisNeighbour.Neighbours[(i+1)%int(g.MaxDegree)]
						if diagonal != nil {
							diagonal.IsForRemoval = true
							diagonal.IsDecided = true
							addNeighboursToQueue(diagonal, q, thisNeighbour, nextNeighbour)
						}
					}
				} else {
					n.IsForRemoval = true
					n.IsDecided = true
					/* We can determine state of diagonal node */
					if thisNeighbour != nil {
						diagonal := thisNeighbour.Neighbours[(i+1)%int(g.MaxDegree)]
						if diagonal != nil {
							diagonal.IsForRemoval = false
							diagonal.IsDecided = true
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

func (n *Node) find31Templates(g *Graph, q *queue.Queue) bool {
	/* value of this node equal to 3 */
	if n.Value == 3 {
		/* is next to the wall */
		for i := 0; i < len(n.Neighbours); i++ {
			thisNeighbour := n.Neighbours[i]
			if isNodeDecided(thisNeighbour) {
				/* is number 1 also next to the wall */
				nextNeigh := n.Neighbours[(i+1)%int(g.MaxDegree)]
				prevNeigh := n.Neighbours[(i-1+int(g.MaxDegree))%int(g.MaxDegree)]
				if (prevNeigh != nil && prevNeigh.Value == 1) || (nextNeigh != nil && prevNeigh.Value == 1) {
					n.IsDecided = true
					n.IsForRemoval = !isNodeDecidedOut(thisNeighbour)
					addNeighboursToQueue(n, q, thisNeighbour)
				}
			}
		}
	}
	return false
}

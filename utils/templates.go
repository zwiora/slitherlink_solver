package utils

import (
	"github.com/golang-collections/collections/queue"
)

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

func addNeighboursToQueueWithExclusions(n *Node, q *queue.Queue, excludedNodes ...*Node) {
	// for i := 0; i < len(n.Neighbours); i++ {
	// 	thisNeighbour := n.Neighbours[i]
	// 	if thisNeighbour != nil && !thisNeighbour.IsDecided {
	// 		isExcluded := false
	// 		for k, v := range excludedNodes {
	// 			if thisNeighbour == v {
	// 				isExcluded = true
	// 				excludedNodes = append(excludedNodes[:k], excludedNodes[k+1:]...)
	// 				break
	// 			}
	// 		}
	// 		if !isExcluded {
	// 			q.Enqueue(n.Neighbours[i])
	// 		}
	// 	}
	// }
}

func addNeighboursToQueue(n *Node, q *queue.Queue) {
	// for i := 0; i < len(n.Neighbours); i++ {
	// 	thisNeighbour := n.Neighbours[i]
	// 	if thisNeighbour != nil && !thisNeighbour.IsDecided {
	// 		q.Enqueue(n.Neighbours[i])
	// 	}
	// }
}

func isNodeDecided(n *Node) bool {
	return n == nil || n.IsDecided
}

func isNodeDecidedOut(n *Node) bool {
	return !(n != nil && n.IsDecided && !n.IsForRemoval)
}

func addNodeToGroup(n *Node, base *Node, q *queue.Queue) {

	if base.IsDecided {
		if n.TemplateGroup == nil {
			n.IsDecided = true
			n.IsForRemoval = base.IsForRemoval
		} else {
			n.TemplateGroup.setValue(base.CanBeRemoved)
		}
	} else {
		l := base.TemplateGroup
		if n.TemplateGroup == nil {
			l.addElement(n)
		} else {
			addLists(l, n.TemplateGroup)
		}
	}

	addNeighboursToQueueWithExclusions(n, q, base)
}

func (n *Node) findZeroTemplates(q *queue.Queue) {
	/* is value 0 */
	if n.Value == 0 {

		if !n.IsDecided && n.TemplateGroup == nil {
			n.TemplateGroup = new(List)
			n.TemplateGroup.addElement(n)
		}

		addNeighboursToQueue(n, q)

		/* check if any neighbour is decided */
		for i := 0; i < len(n.Neighbours); i++ {
			thisNeighbour := n.Neighbours[i]

			if isNodeDecided(thisNeighbour) {
				/* set this node and all neighbours as decided */
				n.TemplateGroup.setValue(isNodeDecidedOut(thisNeighbour))
			} else {
				addNodeToGroup(thisNeighbour, n, q)
			}
		}
	}
}

/* Returns true, if template found */
func (n *Node) findCornerTemplates(g *Graph, q *queue.Queue) bool {

	if n.Value == 2 {
		noDecided := 0
		for i := 0; i < len(n.Neighbours); i++ {
			if n.Neighbours[i] == nil || n.Neighbours[i].IsDecided {
				noDecided++
			}
		}

		if noDecided == 4 {
			return false
		}
	}

	for i := 0; i < len(n.Neighbours); i++ {
		thisNeighbour := n.Neighbours[i]
		nextNeighbour := n.Neighbours[(i+1)%int(g.MaxDegree)]
		isSame, stateOfBoth := isTheSameState(thisNeighbour, nextNeighbour)
		if isSame {
			switch n.Value {
			case 1:
				n.IsForRemoval = stateOfBoth
				n.IsDecided = true
				addNeighboursToQueueWithExclusions(n, q)
			case 2:
				oppositeNeighbour := n.Neighbours[(i+2)%int(g.MaxDegree)]
				if oppositeNeighbour != nil {
					oppositeNeighbour.IsForRemoval = !isNodeDecidedOut(thisNeighbour)
					oppositeNeighbour.IsDecided = true
					addNeighboursToQueueWithExclusions(oppositeNeighbour, q, n)
				}

				oppositeNeighbour = n.Neighbours[(i+3)%int(g.MaxDegree)]
				if oppositeNeighbour != nil {
					oppositeNeighbour.IsForRemoval = !isNodeDecidedOut(thisNeighbour)
					oppositeNeighbour.IsDecided = true
					addNeighboursToQueueWithExclusions(oppositeNeighbour, q, n)
				}
			case 3:
				n.IsForRemoval = !isNodeDecidedOut(thisNeighbour)
				n.IsDecided = true
				/* We can determine state of diagonal node - also possible with other method */
				// if thisNeighbour != nil {
				// 	diagonal := thisNeighbour.Neighbours[(i+1)%int(g.MaxDegree)]
				// 	if diagonal != nil {
				// 		diagonal.IsForRemoval = true
				// 		diagonal.IsDecided = true
				// 		addNeighboursToQueue(diagonal, q, thisNeighbour, nextNeighbour)
				// 	}
				// }

				addNeighboursToQueueWithExclusions(n, q)
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
					addNeighboursToQueueWithExclusions(n, q, thisNeighbour)
				}
			}
		}
	}
	return false
}

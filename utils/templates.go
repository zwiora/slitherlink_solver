package utils

import (
	"github.com/golang-collections/collections/queue"
)

func isTheSameState(n1 *Node, n2 *Node) bool {
	if n1 != nil && n2 != nil && n1.TemplateGroup != nil && n1.TemplateGroup == n2.TemplateGroup {
		return true
	}

	// we can't say if they have the same state
	if (n1 != nil && !n1.IsDecided) || (n2 != nil && !n2.IsDecided) {
		return false
	}

	isN1Out := n1 == nil || (n1.IsDecided && n1.IsForRemoval)
	isN2Out := n2 == nil || (n2.IsDecided && n2.IsForRemoval)

	if isN1Out == isN2Out {
		return true
	}

	return false
}

func isDifferentState(n1 *Node, n2 *Node) bool {
	if n1 != nil && n2 != nil && n1.TemplateGroup != nil && n2.TemplateGroup != nil && n1.TemplateGroup == n2.TemplateGroup.OppositeList {
		return true
	}

	// we can't say if they have different states
	if (n1 != nil && !n1.IsDecided) || (n2 != nil && !n2.IsDecided) {
		return false
	}

	isN1Out := n1 == nil || (n1.IsDecided && n1.IsForRemoval)
	isN2Out := n2 == nil || (n2.IsDecided && n2.IsForRemoval)

	if isN1Out != isN2Out {
		return true
	}

	return false
}

func isNodeDecided(n *Node) bool {
	return n == nil || n.IsDecided
}

func isNodeDecidedOut(n *Node) bool {
	return n == nil || (n.IsDecided && n.IsForRemoval)
}

func nodeState(n *Node) any {
	if n == nil {
		return true
	}

	if n.IsDecided {
		return n.IsForRemoval
	}

	if n.TemplateGroup != nil {
		return n.TemplateGroup
	}

	return nil
}

/* Returns true if any changes have been applied */
func addNodeToGroup(n1 *Node, n2 *Node, g *Graph) bool {
	if (isNodeDecided(n1) && isNodeDecided(n2)) || isTheSameState(n1, n2) {
		return false
	}

	var decided *Node
	var notDecided *Node

	if isNodeDecided(n1) {
		decided = n1
		notDecided = n2
	} else if isNodeDecided(n2) {
		decided = n2
		notDecided = n1
	}

	if notDecided != nil {
		if notDecided.TemplateGroup == nil {
			notDecided.IsDecided = true
			notDecided.IsForRemoval = decided == nil || decided.IsForRemoval
		} else {
			notDecided.TemplateGroup.SetValue(isNodeDecidedOut(decided), nil, g)
		}
	} else /*Neither one is decided */ {
		if n1.TemplateGroup != nil && n2.TemplateGroup != nil {
			addLists(n1.TemplateGroup, n2.TemplateGroup)
		} else if n1.TemplateGroup != nil {
			n1.TemplateGroup.addElement(n2)
		} else {
			n2.TemplateGroup.addElement(n1)
		}
	}

	return true
}

/* Returns true if any changes have been applied */
func addNodeToOppositeGroup(n1 *Node, n2 *Node, g *Graph) bool {
	if isNodeDecided(n1) && isNodeDecided(n2) || isDifferentState(n1, n2) {
		return false
	}

	var decided *Node
	var notDecided *Node

	if isNodeDecided(n1) {
		decided = n1
		notDecided = n2
	} else if isNodeDecided(n2) {
		decided = n2
		notDecided = n1
	}

	if notDecided != nil {

		if notDecided.TemplateGroup == nil {
			notDecided.IsDecided = true
			if decided == nil {
				notDecided.IsForRemoval = false
			} else {
				notDecided.IsForRemoval = !decided.IsForRemoval
			}
		} else {
			notDecided.TemplateGroup.SetValue(!isNodeDecidedOut(decided), nil, g)
		}
	} else /*Neither one is decided */ {

		if n1.TemplateGroup != nil && n2.TemplateGroup != nil {
			addOppositeLists(n1.TemplateGroup, n2.TemplateGroup)
		} else if n1.TemplateGroup != nil {
			n1.TemplateGroup.addOppositeElement(n2)
		} else {
			n2.TemplateGroup.addOppositeElement(n1)
		}
	}

	return true
}

func (n *Node) findZeroTemplates(g *Graph) bool {
	/* is value 0 */
	if n.Value == 0 {

		if !n.IsDecided && n.TemplateGroup == nil {
			n.TemplateGroup = new(List)
			n.TemplateGroup.addElement(n)
		}

		/* check if any neighbour is decided */
		for i := 0; i < len(n.Neighbours); i++ {
			thisNeighbour := n.Neighbours[i]

			if isNodeDecided(thisNeighbour) {
				/* set this node and all neighbours as decided */
				n.TemplateGroup.SetValue(isNodeDecidedOut(thisNeighbour), nil, g)
			} else {
				addNodeToGroup(thisNeighbour, n, g)
			}
		}
		return true
	}
	return false
}

func (n *Node) findNumberTemplates(g *Graph) bool {

	/* ! DZIAŁA TYLKO DLA KWADRATÓW */
	if n.Value != -1 && n.Value != 0 {
		if n.IsDecided || n.TemplateGroup != nil {
			nState := nodeState(n)

			stateList := make(map[any][]int)
			for k, v := range n.Neighbours {
				vState := nodeState(v)
				stateList[vState] = append(stateList[vState], k)
			}

			if len(stateList[nState]) == len(n.Neighbours)-int(n.Value) {
				isChangeMade := false
				for key, slice := range stateList {
					if key != nState {
						for v := range slice {
							if addNodeToOppositeGroup(n, n.Neighbours[slice[v]], g) {
								isChangeMade = true
							}
						}
					}
				}
				return isChangeMade
			}

			return false

		}
		if !n.IsDecided && n.Value != 2 {
			stateList := make(map[any][]int)
			for k, v := range n.Neighbours {
				vState := nodeState(v)
				if vState != nil {
					stateList[vState] = append(stateList[vState], k)
				}
			}

			for key, slice := range stateList {
				if key != nil && len(slice) >= 2 {
					isChangeMade := false
					if n.Value == 1 {
						isChangeMade = addNodeToGroup(n, n.Neighbours[slice[0]], g)
						// time.Sleep(1000 * time.Millisecond)
					} else if n.Value == 3 {
						isChangeMade = addNodeToOppositeGroup(n, n.Neighbours[slice[0]], g)
					}
					return isChangeMade

				}
			}
		}
	}

	return false
}

// /* Returns true, if template found */
// func (n *Node) findCornerTemplates(g *Graph, q *queue.Queue) bool {

// 	// fmt.Println(n)

// 	if n.Value == 2 {
// 		noDecided := 0
// 		for i := 0; i < len(n.Neighbours); i++ {
// 			if n.Neighbours[i] == nil || n.Neighbours[i].IsDecided {
// 				noDecided++
// 			}
// 		}

// 		// fmt.Println(noDecided)

// 		if noDecided == 4 {
// 			// fmt.Println("Wracamy")
// 			return false
// 		}
// 	}

// 	// fmt.Println("Nie wracamy")

// 	for i := 0; i < len(n.Neighbours); i++ {
// 		thisNeighbour := n.Neighbours[i]
// 		nextNeighbour := n.Neighbours[(i+1)%int(g.MaxDegree)]
// 		isSame, stateOfBoth := isTheSameState(thisNeighbour, nextNeighbour)
// 		if isSame && n.Value != 0 {
// 			switch n.Value {
// 			case 1:
// 				n.IsForRemoval = stateOfBoth
// 				n.IsDecided = true
// 			case 2:
// 				oppositeNeighbour := n.Neighbours[(i+2)%int(g.MaxDegree)]
// 				if oppositeNeighbour != nil {
// 					oppositeNeighbour.IsForRemoval = !isNodeDecidedOut(thisNeighbour)
// 					oppositeNeighbour.IsDecided = true
// 				}

// 				oppositeNeighbour = n.Neighbours[(i+3)%int(g.MaxDegree)]
// 				if oppositeNeighbour != nil {
// 					oppositeNeighbour.IsForRemoval = !isNodeDecidedOut(thisNeighbour)
// 					oppositeNeighbour.IsDecided = true
// 				}
// 			case 3:
// 				n.IsForRemoval = !isNodeDecidedOut(thisNeighbour)
// 				n.IsDecided = true
// 				/* We can determine state of diagonal node - also possible with other method */
// 				// if thisNeighbour != nil {
// 				// 	diagonal := thisNeighbour.Neighbours[(i+1)%int(g.MaxDegree)]
// 				// 	if diagonal != nil {
// 				// 		diagonal.IsForRemoval = true
// 				// 		diagonal.IsDecided = true
// 				// 		addNeighboursToQueue(diagonal, q, thisNeighbour, nextNeighbour)
// 				// 	}
// 				// }

// 			}

// 			return true
// 		}
// 	}

// 	return false
// }

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
				if (prevNeigh != nil && prevNeigh.Value == 1 && isNodeDecidedOut(prevNeigh.Neighbours[i]) == isNodeDecidedOut(thisNeighbour)) || (nextNeigh != nil && nextNeigh.Value == 1 && isNodeDecidedOut(nextNeigh.Neighbours[i]) == isNodeDecidedOut(thisNeighbour)) {
					n.IsDecided = true
					n.IsForRemoval = !isNodeDecidedOut(thisNeighbour)
				}
			}
		}
	}
	return false
}

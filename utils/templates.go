package utils

import (
	"fmt"

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

func nodeOppositeState(n *Node) any {
	if n == nil {
		return false
	}

	if n.IsDecided {
		return !n.IsForRemoval
	}

	if n.TemplateGroup != nil && n.TemplateGroup.OppositeList != nil {
		return n.TemplateGroup.OppositeList
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
			if n2.TemplateGroup == nil {
				n2.TemplateGroup = new(List)
				n2.TemplateGroup.addElement(n2)
			}
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
			fmt.Println(n1.TemplateGroup)
			fmt.Println(n2.TemplateGroup)
			addOppositeLists(n1.TemplateGroup, n2.TemplateGroup)
		} else if n1.TemplateGroup != nil {
			n1.TemplateGroup.addOppositeElement(n2)

		} else {
			if n2.TemplateGroup == nil {
				n2.TemplateGroup = new(List)
				n2.TemplateGroup.addElement(n2)
			}
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

	isChangeMade := false

	/* ! DZIAŁA TYLKO DLA KWADRATÓW */
	if n.Value != -1 && n.Value != 0 {
		if n.IsDecided || n.TemplateGroup != nil {
			nState := nodeState(n)

			stateList := make(map[any][]int)
			for k, v := range n.Neighbours {
				vState := nodeState(v)
				stateList[vState] = append(stateList[vState], k)
			}

			/* !!!!!!!!!!!!!!!!!!!!!!!!! */
			if len(stateList[nState]) == len(n.Neighbours)-int(n.Value) {
				for key, slice := range stateList {
					if key != nState {
						for v := range slice {
							if addNodeToOppositeGroup(n, n.Neighbours[slice[v]], g) {
								isChangeMade = true
							}
						}
					}
				}
			} else if n.Value == 1 && len(stateList[nState]) == 2 {
				var firstNode *Node
				var secondNode *Node
				for key, slice := range stateList {
					if key != nState {
						for v := range slice {
							if firstNode == nil {
								firstNode = n.Neighbours[slice[v]]
							} else {
								secondNode = n.Neighbours[slice[v]]
							}
						}
					}
				}

				if addNodeToOppositeGroup(firstNode, secondNode, g) {
					isChangeMade = true
				}
			}

			oppositeState := nodeOppositeState(n)

			if oppositeState != nil && len(stateList[oppositeState]) == int(n.Value) {
				for key, slice := range stateList {
					if key != oppositeState {
						for v := range slice {
							if addNodeToGroup(n, n.Neighbours[slice[v]], g) {
								isChangeMade = true
							}
						}
					}
				}
			} else if n.Value == 3 && len(stateList[oppositeState]) == 2 {
				var firstNode *Node
				var secondNode *Node
				for key, slice := range stateList {
					if key != oppositeState {
						for v := range slice {
							if firstNode == nil {
								firstNode = n.Neighbours[slice[v]]
							} else {
								secondNode = n.Neighbours[slice[v]]
							}
						}
					}
				}

				if addNodeToOppositeGroup(firstNode, secondNode, g) {
					isChangeMade = true
				}
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
					if n.Value == 1 {
						if addNodeToGroup(n, n.Neighbours[slice[0]], g) {
							isChangeMade = true
						}
					} else if n.Value == 3 {
						if addNodeToOppositeGroup(n, n.Neighbours[slice[0]], g) {
							isChangeMade = true
						}
					}
				}
			}

		}
		if n.Value == 2 {
			stateList := make(map[any][]int)
			for k, v := range n.Neighbours {
				vState := nodeState(v)
				stateList[vState] = append(stateList[vState], k)
			}

			for key, slice := range stateList {
				if key != nil && len(slice) >= 2 {
					for k, s := range stateList {
						if k != key {
							for v := range s {
								if addNodeToOppositeGroup(n.Neighbours[slice[0]], n.Neighbours[s[v]], g) {
									isChangeMade = true
								}
							}
						}
					}
					break
				}
			}

			if !isChangeMade {

				if len(stateList[true]) == 1 && len(stateList[false]) == 1 {

					var firstNode *Node
					var secondNode *Node
					for key, slice := range stateList {
						if key != true && key != false {
							for v := range slice {
								if firstNode == nil {
									firstNode = n.Neighbours[slice[v]]
								} else {
									secondNode = n.Neighbours[slice[v]]
								}
							}
						}
					}

					// g.PrintSquaresBoard(true)
					// fmt.Println(n)
					// fmt.Println(firstNode, secondNode)

					if addNodeToOppositeGroup(firstNode, secondNode, g) {
						isChangeMade = true
					}
				}

				for key := range stateList {
					if key != nil && key != true && key != false {
						list := key.(*List)
						if list != nil && list.OppositeList != nil && len(stateList[list.OppositeList]) == 1 {
							var firstNode *Node
							var secondNode *Node
							for k, s := range stateList {
								if k != key && k != list.OppositeList {
									for v := range s {
										if firstNode == nil {
											firstNode = n.Neighbours[s[v]]
										} else {
											secondNode = n.Neighbours[s[v]]
										}
									}
								}
							}

							if addNodeToOppositeGroup(firstNode, secondNode, g) {
								// g.PrintSquaresBoard(true)
								// fmt.Println("Neighbours")
								// for _, v := range n.Neighbours {
								// 	fmt.Println(v)
								// }
								// fmt.Println(firstNode, secondNode)

								isChangeMade = true
							}
						}
					}
				}

			}
		}
	}

	return isChangeMade
}

func (n *Node) findContinousSquareTemplates(g *Graph) bool {
	if !n.IsDecided {
		for i := range n.Neighbours {
			j := (i + 1) % len(n.Neighbours)
			if n.Neighbours[i] != nil && n.Neighbours[j] != nil && n.Neighbours[i].Neighbours[j] != nil {
				if isTheSameState(n.Neighbours[i], n.Neighbours[j]) && isDifferentState(n.Neighbours[i], n.Neighbours[i].Neighbours[j]) {
					return addNodeToGroup(n, n.Neighbours[i], g)
				}
			}
		}
	}

	return false
}

func (n *Node) find33Templates(g *Graph) bool {
	isChangeMade := false
	if n.Value == 3 {
		m := n.Neighbours[0]
		if m != nil && m.Value == 3 && !(n.IsDecided && m.IsDecided) {
			if addNodeToGroup(n, m.Neighbours[0], g) {
				isChangeMade = true
			}
			if addNodeToGroup(m, n.Neighbours[2], g) {
				isChangeMade = true
			}
			if addNodeToOppositeGroup(n, m, g) {
				isChangeMade = true
			}
			if addNodeToGroup(n.Neighbours[3], m.Neighbours[3], g) {
				isChangeMade = true
			}
			if addNodeToGroup(n.Neighbours[1], m.Neighbours[1], g) {
				isChangeMade = true
			}
		}

		m = n.Neighbours[1]
		if m != nil && m.Value == 3 && !(n.IsDecided && m.IsDecided) {
			if addNodeToGroup(n, m.Neighbours[1], g) {
				isChangeMade = true
			}
			if addNodeToGroup(m, n.Neighbours[3], g) {
				isChangeMade = true
			}
			if addNodeToOppositeGroup(n, m, g) {
				isChangeMade = true
			}
			if addNodeToGroup(n.Neighbours[0], m.Neighbours[0], g) {
				isChangeMade = true
			}
			if addNodeToGroup(n.Neighbours[2], m.Neighbours[2], g) {
				isChangeMade = true
			}
		}

	}
	return isChangeMade
}

func (n *Node) find3and3Templates(g *Graph) bool {
	isChangeMade := false
	if n.Value == 3 {
		tmp := n.Neighbours[0]
		if tmp != nil {
			m := tmp.Neighbours[1]
			if m != nil && m.Value == 3 && !(n.IsDecided && m.IsDecided) {
				if addNodeToOppositeGroup(n, n.Neighbours[2], g) {
					isChangeMade = true
				}
				if addNodeToOppositeGroup(n, n.Neighbours[3], g) {
					isChangeMade = true
				}
				if addNodeToOppositeGroup(m, m.Neighbours[0], g) {
					isChangeMade = true
				}
				if addNodeToOppositeGroup(m, m.Neighbours[1], g) {
					isChangeMade = true
				}
			}
		}

		tmp = n.Neighbours[1]
		if tmp != nil {
			m := tmp.Neighbours[2]
			if m != nil && m.Value == 3 && !(n.IsDecided && m.IsDecided) {
				if addNodeToOppositeGroup(n, n.Neighbours[0], g) {
					isChangeMade = true
				}
				if addNodeToOppositeGroup(n, n.Neighbours[3], g) {
					isChangeMade = true
				}
				if addNodeToOppositeGroup(m, m.Neighbours[2], g) {
					isChangeMade = true
				}
				if addNodeToOppositeGroup(m, m.Neighbours[1], g) {
					isChangeMade = true
				}
				// if m.TemplateGroup != nil {
				// 	fmt.Println(m)
				// 	fmt.Println(m.Neighbours[1])
				// 	fmt.Println(m.Neighbours[2])
				// 	fmt.Println(isDifferentState(m, m.Neighbours[1]))
				// 	m.TemplateGroup.print()
				// 	fmt.Println("-----------")
				// }

			}
		}
	}
	return isChangeMade
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

package utils

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
			}
		}
	}
	return isChangeMade
}

func (n *Node) findloopReachingNumberTemplates(g *Graph) bool {
	isChangeMade := false

	if n.Value == 3 {
		for i, v := range n.Neighbours {
			if v != nil {
				w := v.Neighbours[(i+1)%len(n.Neighbours)]
				if isDifferentState(v, w) {
					if n.Neighbours[(i+1)%len(n.Neighbours)] != nil {
						if addNodeToGroup(n.Neighbours[(i+1)%len(n.Neighbours)], w, g) {
							isChangeMade = true
						}
					}
					if addNodeToOppositeGroup(n, n.Neighbours[(i+2)%len(n.Neighbours)], g) {
						isChangeMade = true
					}
					if addNodeToOppositeGroup(n, n.Neighbours[(i+3)%len(n.Neighbours)], g) {
						isChangeMade = true
					}
				}

				w = v.Neighbours[(i-1+len(n.Neighbours))%len(n.Neighbours)]
				if isDifferentState(v, w) {
					if n.Neighbours[(i-1+len(n.Neighbours))%len(n.Neighbours)] != nil {
						if addNodeToGroup(n.Neighbours[(i-1+len(n.Neighbours))%len(n.Neighbours)], w, g) {
							isChangeMade = true
						}
					}
					if addNodeToOppositeGroup(n, n.Neighbours[(i+1)%len(n.Neighbours)], g) {
						isChangeMade = true
					}
					if addNodeToOppositeGroup(n, n.Neighbours[(i+2)%len(n.Neighbours)], g) {
						isChangeMade = true
					}
				}
			}
		}
	} else if n.Value == 1 {
		for i, v := range n.Neighbours {
			if v != nil {
				w := v.Neighbours[(i+1)%len(n.Neighbours)]
				x := n.Neighbours[(i+1)%len(n.Neighbours)]
				if isDifferentState(v, w) && isTheSameState(w, x) {
					if addNodeToGroup(n, n.Neighbours[(i+2)%len(n.Neighbours)], g) {
						isChangeMade = true
					}
					if addNodeToGroup(n, n.Neighbours[(i+3)%len(n.Neighbours)], g) {
						isChangeMade = true
					}
				}

				w = v.Neighbours[(i-1+len(n.Neighbours))%len(n.Neighbours)]
				x = n.Neighbours[(i-1+len(n.Neighbours))%len(n.Neighbours)]
				if isDifferentState(v, w) && isTheSameState(w, x) {
					if addNodeToGroup(n, n.Neighbours[(i+1)%len(n.Neighbours)], g) {
						isChangeMade = true
					}
					if addNodeToGroup(n, n.Neighbours[(i+2)%len(n.Neighbours)], g) {
						isChangeMade = true
					}
				}
			}
		}
	}

	return isChangeMade
}

func (n *Node) find31Templates(g *Graph) bool {
	isChangeMade := false
	/* value of this node equal to 3 */
	if n.Value == 3 {
		/* is next to the wall */
		for i := 0; i < len(n.Neighbours); i++ {
			thisNeighbour := n.Neighbours[i]
			i1 := (i + 1) % int(g.MaxDegree)
			i2 := (i - 1 + int(g.MaxDegree)) % int(g.MaxDegree)
			if n.Neighbours[i1] != nil && n.Neighbours[i1].Value == 1 && (thisNeighbour == nil || isTheSameState(thisNeighbour, thisNeighbour.Neighbours[i1])) {
				if addNodeToOppositeGroup(n, thisNeighbour, g) {
					isChangeMade = true
				}
				break
			}
			if n.Neighbours[i2] != nil && n.Neighbours[i2].Value == 1 && (thisNeighbour == nil || isTheSameState(thisNeighbour, thisNeighbour.Neighbours[i2])) {
				if addNodeToOppositeGroup(n, thisNeighbour, g) {
					isChangeMade = true
				}
				break
			}
		}
	}
	return isChangeMade
}

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

func isOppositeState(n1 *Node, n2 *Node) bool {
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

	if isNodeDecided(n1) && isNodeDecided(n2) || isOppositeState(n1, n2) {
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
		/* Based on the state of n */
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
	if g.Shape == "square" {
		if !n.IsDecided {
			for i := range n.Neighbours {
				j := (i + 1) % len(n.Neighbours)
				if n.Neighbours[i] != nil && n.Neighbours[j] != nil && n.Neighbours[i].Neighbours[j] != nil {
					if isTheSameState(n.Neighbours[i], n.Neighbours[j]) && isOppositeState(n.Neighbours[i], n.Neighbours[i].Neighbours[j]) {
						return addNodeToGroup(n, n.Neighbours[i], g)
					}
				}
			}
		}
	} else if g.Shape == "triangle" {
		if !n.IsDecided {
			for k := range n.Neighbours {
				firstNeighbour := n.Neighbours[k]
				secondNeighbour := n.Neighbours[(k+1)%3]
				if isTheSameState(firstNeighbour, secondNeighbour) {
					tmp := firstNeighbour
					i := k
					for tmp != secondNeighbour && tmp != nil {
						i = (i - 1 + 3) % 3
						tmp = tmp.Neighbours[i]
						if isOppositeState(tmp, firstNeighbour) {
							return addNodeToGroup(n, firstNeighbour, g)
						}
					}

					i = (k + 1) % 3

					for tmp != firstNeighbour && tmp != nil {
						i = (i + 1 + 3) % 3
						tmp = tmp.Neighbours[i]
						if isOppositeState(tmp, secondNeighbour) {
							return addNodeToGroup(n, secondNeighbour, g)
						}
					}
				}

			}
		}
	}

	return false
}

func (n *Node) find33Templates(g *Graph) bool {
	isChangeMade := false
	if g.Shape == "square" {
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
				if addNodeToOppositeGroup(n.Neighbours[3], n.Neighbours[1], g) {
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
				if addNodeToOppositeGroup(n.Neighbours[0], n.Neighbours[2], g) {
					isChangeMade = true
				}
			}

		}
	} else if g.Shape == "honeycomb" {
		if n.Value == 5 {

			for i := 0; i < 3; i++ {
				m := n.Neighbours[i]
				if m != nil && m.Value == 5 && !(n.IsDecided && m.IsDecided) {
					if addNodeToOppositeGroup(n, m, g) {
						isChangeMade = true
					}

					if addNodeToOppositeGroup(n.Neighbours[(i-1+6)%6], n.Neighbours[(i+1)%6], g) {
						isChangeMade = true
					}

					if addNodeToGroup(n, m.Neighbours[i], g) {
						isChangeMade = true
					}
					if addNodeToGroup(n, m.Neighbours[(i-1+6)%6], g) {
						isChangeMade = true
					}
					if addNodeToGroup(n, m.Neighbours[(i+1)%6], g) {
						isChangeMade = true
					}

					if addNodeToGroup(m, n.Neighbours[(i+2)%6], g) {
						isChangeMade = true
					}
					if addNodeToGroup(m, n.Neighbours[(i+3)%6], g) {
						isChangeMade = true
					}
					if addNodeToGroup(m, n.Neighbours[(i+4)%6], g) {
						isChangeMade = true
					}
				}
			}

		}
	} else if g.Shape == "triangle" {
		if n.Value == 2 {
			for k := range n.Neighbours {
				neighbour := n.Neighbours[k]
				if neighbour != nil && neighbour.Value == 2 && !isOppositeState(n, neighbour) {
					if addNodeToOppositeGroup(n, neighbour, g) {
						isChangeMade = true
					}

					i := (k + 1) % 3
					if addNodeToOppositeGroup(n.Neighbours[i], neighbour.Neighbours[i], g) {
						isChangeMade = true
					}

					base := n.Neighbours[i]
					tmp := base
					for tmp != nil {
						i = (i + 1) % 3
						tmp = tmp.Neighbours[i]
						if tmp == neighbour {
							break
						}
						if addNodeToGroup(tmp, base, g) {
							isChangeMade = true
						}
					}

					i = (k + 2) % 3
					base = n.Neighbours[i]
					tmp = base
					for tmp != nil {
						i = (i - 1 + 3) % 3
						tmp = tmp.Neighbours[i]
						if tmp == neighbour {
							break
						}
						if addNodeToGroup(tmp, base, g) {
							isChangeMade = true
						}
					}

					i = (k + 1) % 3
					base = neighbour.Neighbours[i]
					tmp = base
					for tmp != nil {
						i = (i + 1) % 3
						tmp = tmp.Neighbours[i]
						if tmp == n {
							break
						}
						if addNodeToGroup(tmp, base, g) {
							isChangeMade = true
						}
					}

					i = (k + 2) % 3
					base = neighbour.Neighbours[i]
					tmp = base
					for tmp != nil {
						i = (i - 1 + 3) % 3
						tmp = tmp.Neighbours[i]
						if tmp == n {
							break
						}
						if addNodeToGroup(tmp, base, g) {
							isChangeMade = true
						}
					}
				}
			}
		}
	}

	return isChangeMade
}

func (n *Node) find33CornerTemplates(g *Graph) bool {
	isChangeMade := false
	if g.Shape == "square" {
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
	} else if g.Shape == "triangle" {
		if n.Value == 2 {
			for k := 0; k < 3; k++ {

				i := k
				tmp := n.Neighbours[i]
				var m *Node
				for tmp != nil {
					i = (i + 1) % 3
					tmp = tmp.Neighbours[i]

					if tmp == n.Neighbours[(k+2)%3] || tmp == nil {
						break
					}
					if tmp.Value == 2 {
						m = tmp
						if addNodeToOppositeGroup(m, m.Neighbours[(i-1+3)%3], g) {
							isChangeMade = true
						}
						break
					}
				}

				if m == nil {
					i = (k - 1 + 3) % 3
					tmp = n.Neighbours[i]
					for tmp != nil {
						i = (i - 1 + 3) % 3
						tmp = tmp.Neighbours[i]

						if tmp == n.Neighbours[k] || tmp == nil {
							break
						}
						if tmp.Value == 2 {
							m = tmp
							if addNodeToOppositeGroup(m, m.Neighbours[(i+1)%3], g) {
								isChangeMade = true
							}
							break
						}

					}
				}

				if m == nil {
					continue
				}

				if addNodeToOppositeGroup(n, n.Neighbours[(k+1)%3], g) {
					isChangeMade = true
				}

				i = k
				tmp = n.Neighbours[i]
				base := tmp
				isOppositeOn := false
				for tmp != nil {
					i = (i + 1) % 3
					tmp = tmp.Neighbours[i]

					if tmp == n {
						break
					}
					if tmp == m {
						isOppositeOn = true
						continue
					}
					if isOppositeOn {
						if addNodeToOppositeGroup(tmp, base, g) {
							isChangeMade = true
						}
					} else {
						if addNodeToGroup(tmp, base, g) {
							isChangeMade = true
						}
					}
				}

				if !isOppositeState(n.Neighbours[k], n.Neighbours[(k+2)%3]) {
					i = (k - 1 + 3) % 3
					tmp = n.Neighbours[i]
					for tmp != nil {
						i = (i - 1 + 3) % 3
						tmp = tmp.Neighbours[i]

						if tmp == n {
							break
						}
						if tmp == m {
							isOppositeOn = true
							continue
						}
						if isOppositeOn {
							if addNodeToOppositeGroup(tmp, base, g) {
								isChangeMade = true
							}
						} else {
							if addNodeToGroup(tmp, base, g) {
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

func (n *Node) findloopReachingNumberTemplates(g *Graph) bool {
	isChangeMade := false

	if g.Shape == "square" {
		if n.Value == 3 {
			for i, v := range n.Neighbours {
				if v != nil {
					w := v.Neighbours[(i+1)%len(n.Neighbours)]
					if isOppositeState(v, w) {
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
					if isOppositeState(v, w) {
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
					if isOppositeState(v, w) && isTheSameState(w, x) {
						if addNodeToGroup(n, n.Neighbours[(i+2)%len(n.Neighbours)], g) {
							isChangeMade = true
						}
						if addNodeToGroup(n, n.Neighbours[(i+3)%len(n.Neighbours)], g) {
							isChangeMade = true
						}
					}

					w = v.Neighbours[(i-1+len(n.Neighbours))%len(n.Neighbours)]
					x = n.Neighbours[(i-1+len(n.Neighbours))%len(n.Neighbours)]
					if isOppositeState(v, w) && isTheSameState(w, x) {
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
	} else if g.Shape == "honeycomb" {
		if n.Value == 5 {
			for k, v := range n.Neighbours {
				if v != nil && isOppositeState(v, n.Neighbours[(k+1)%6]) {
					for j := 2; j < 6; j++ {
						if addNodeToOppositeGroup(n, n.Neighbours[(k+j)%6], g) {
							isChangeMade = true
						}
					}
				}
			}
		} else if n.Value == 1 {
			for k, v := range n.Neighbours {
				if v != nil && isOppositeState(v, n.Neighbours[(k+1)%6]) {
					for j := 2; j < 6; j++ {
						if addNodeToGroup(n, n.Neighbours[(k+j)%6], g) {
							isChangeMade = true
						}
					}
				}
			}
		}
	} else if g.Shape == "triangle" {
		if n.Value == 2 {

			for k := 0; k < 3; k++ {
				i := k
				tmp := n.Neighbours[i]
				old := tmp
				var m *Node
				for tmp != nil {

					if m != nil {
						if addNodeToOppositeGroup(m, tmp, g) {
							isChangeMade = true
						}
					}

					i = (i + 1) % 3
					tmp = tmp.Neighbours[i]

					if tmp == n {
						break
					}

					if isOppositeState(old, tmp) {
						if addNodeToOppositeGroup(n, n.Neighbours[(k+1)%3], g) {
							isChangeMade = true
						}
						m = old
					}

					old = tmp
				}

				if m != nil {
					i = k
					tmp = n.Neighbours[i]
					for {
						if addNodeToOppositeGroup(old, tmp, g) {
							isChangeMade = true
						}
						if tmp == nil || tmp == m {
							break
						}
						i = (i + 1) % 3
						tmp = tmp.Neighbours[i]
					}
				} else {
					i = (k - 1 + 3) % 3
					tmp = n.Neighbours[i]
					old = tmp
					for tmp != nil {

						if m != nil {
							if addNodeToOppositeGroup(m, tmp, g) {
								isChangeMade = true
							}
						}

						i = (i - 1 + 3) % 3
						tmp = tmp.Neighbours[i]

						if tmp == n {
							break
						}
						if isOppositeState(old, tmp) {
							if addNodeToOppositeGroup(n, n.Neighbours[(k+1)%3], g) {
								isChangeMade = true
							}
							m = old
						}

						old = tmp
					}

					if m != nil {
						i = (k - 1 + 3) % 3
						tmp = n.Neighbours[i]
						for {
							if addNodeToOppositeGroup(old, tmp, g) {
								isChangeMade = true
							}

							if tmp == nil || tmp == m {
								break
							}
							i = (i - 1 + 3) % 3
							tmp = tmp.Neighbours[i]
						}
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

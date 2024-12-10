package solvers

import (
	"container/heap"
	"slitherlink_solver/debug"
	"slitherlink_solver/utils"
	"time"

	"github.com/golang-collections/collections/stack"
)

/*
* Checks if node can be removed from the loop
* Second bool is true if we can't delete part of template because of the second rule
 */
func checkIfCanBeRemoved(n *utils.Node, g *utils.Graph) (bool, bool) {
	debug.Println("Checking if can be removed: ")
	debug.Println(n)

	/* Checking if removal would create loop inside the loop */
	nodeDegree := n.GetDegree()
	if nodeDegree == 0 || nodeDegree == int(g.MaxDegree) {
		debug.Println("\t- SKIP: removal would create loop inside the loop")
		return false, false
	}
	debug.Println("\t- degree ok")

	if n.IsDeletionBreakingSecondRule() {
		debug.Println("\t- SKIP: deletion isn't in the solution")
		if n.IsDecided && n.IsForRemoval {
			return false, true
		}
		return false, false
	}

	/* It can be removed if it's a leaf */
	if n.GetDegree() == 1 {
		debug.Println("\t- it's a leaf")
		return true, false
	}
	debug.Println("\t- not a leaf")

	/* Checking if its neighbour is leaf */
	for _, v := range n.Neighbours {
		if v != nil && v.IsInLoop {
			if v.GetDegree() == 1 {
				debug.Println("\t- SKIP: neighbour is a leaf")
				return false, false
			}
		}
	}
	debug.Println("\t- neighbour isn't a leaf")

	/* Checking if is between two sides of graph*/
	if g.Shape != "triangle" {
		sidesCounter := 0
		for i := 0; i < len(n.Neighbours); i++ {
			thisNeighbour := n.Neighbours[i]
			nextNeighbour := n.Neighbours[(i+1)%int(g.MaxDegree)]
			if ((thisNeighbour == nil || !thisNeighbour.IsInLoop) && (nextNeighbour != nil && nextNeighbour.IsInLoop)) || ((nextNeighbour == nil || !nextNeighbour.IsInLoop) && (thisNeighbour != nil && thisNeighbour.IsInLoop)) {
				sidesCounter++
			}
			if sidesCounter == 3 {
				debug.Println("\t- SKIP: deletion would create two separate graphs")
				return false, false
			}
		}
		debug.Println("\t- deletion wouldn't create two separate graphs")
	}

	/* Checking if is connected via corner*/
	if g.Shape == "square" {
		for k, v := range n.Neighbours {
			if v != nil && v.IsInLoop {
				diagonalNode := v.Neighbours[(k+1)%int(g.MaxDegree)]
				if diagonalNode != nil && !diagonalNode.IsInLoop {
					nextNeighbour := diagonalNode.Neighbours[(k+2)%int(g.MaxDegree)]
					if nextNeighbour.IsInLoop {
						debug.Println("\t- SKIP: deletion would create two graphs with common corner")
						return false, false
					}
				}
			}
		}
	} else if g.Shape == "triangle" {
		for k := range n.Neighbours {
			firstNeighbour := n.Neighbours[k]
			secondNeighbour := n.Neighbours[(k+1)%3]
			if firstNeighbour != nil && firstNeighbour.IsInLoop && secondNeighbour != nil && secondNeighbour.IsInLoop {
				tmp := firstNeighbour
				i := (k - 1 + 3) % 3
				for tmp != secondNeighbour {
					tmp = tmp.Neighbours[i]
					if tmp == nil || !tmp.IsInLoop {
						debug.Println("\t- SKIP: deletion would create two graphs with common corner")
						return false, false
					}
					i = (i - 1 + 3) % 3
				}
			}
		}
	}
	/* No troubles found - node can be deleted */
	return true, false
}

/* Updates heap of nodes that can be removed from the loop - checks if can be removed and its cost*/
func updateAvailableMoves(n *utils.Node, g *utils.Graph) bool {
	stopTesting := false
	canBeRemoved := false
	for i := 0; i < int(g.MaxDegree); i++ {
		/* neighbouring node */
		thisNode := n.Neighbours[(i)%int(g.MaxDegree)]
		if g.Shape != "triangle" {
			if thisNode != nil {

				if thisNode.IsInLoop && !thisNode.IsVisited && !(thisNode.IsDecided && !thisNode.IsForRemoval && thisNode.TemplateGroup == nil) {
					canBeRemoved, stopTesting = checkIfCanBeRemoved(thisNode, g)
					if canBeRemoved {
						if !thisNode.CanBeRemoved {
							thisNode.SetNodeCost(g)
							heap.Push(g.AvailableMoves, thisNode)
							thisNode.CanBeRemoved = true
						} else {
							thisNode.UpdateNodeCost(g)
						}
					} else if !canBeRemoved && thisNode.CanBeRemoved {
						heap.Remove(g.AvailableMoves, thisNode.QueueIndex)
						thisNode.CanBeRemoved = false
					}
				}

				/* neighbour of the neighbour - only updates cost*/
				nextNode := thisNode.Neighbours[(i)%int(g.MaxDegree)]
				if nextNode != nil && nextNode.IsInLoop && nextNode.CanBeRemoved && !nextNode.IsVisited && !(nextNode.IsDecided && !nextNode.IsForRemoval && nextNode.TemplateGroup == nil) {
					nextNode.UpdateNodeCost(g)
				}
				/* another neighbour of the neighbour - in case of honeycomb */
				if g.Shape == "honeycomb" {
					nextNode = thisNode.Neighbours[(i+1)%int(g.MaxDegree)]
					if nextNode != nil && nextNode.IsInLoop && nextNode.CanBeRemoved && !nextNode.IsVisited && !(nextNode.IsDecided && !nextNode.IsForRemoval && nextNode.TemplateGroup == nil) {
						nextNode.UpdateNodeCost(g)
					}

					previousNode := n.Neighbours[(i-1+int(g.MaxDegree))%int(g.MaxDegree)]
					if previousNode == nil {
						nextNode = thisNode.Neighbours[(i-1+int(g.MaxDegree))%int(g.MaxDegree)]
						if nextNode != nil && nextNode.IsInLoop && nextNode.CanBeRemoved && !nextNode.IsVisited && !(nextNode.IsDecided && !nextNode.IsForRemoval && nextNode.TemplateGroup == nil) {
							nextNode.UpdateNodeCost(g)
						}
					}
				}

				/* diagonal node */
				if g.Shape == "square" {
					thisNode := thisNode.Neighbours[(i+1)%int(g.MaxDegree)]
					if thisNode != nil && thisNode.IsInLoop && !thisNode.IsVisited && !(thisNode.IsDecided && !thisNode.IsForRemoval && thisNode.TemplateGroup == nil) {
						canBeRemoved, stopTesting = checkIfCanBeRemoved(thisNode, g)
						if canBeRemoved {
							if !thisNode.CanBeRemoved {
								thisNode.SetNodeCost(g)
								heap.Push(g.AvailableMoves, thisNode)
								thisNode.CanBeRemoved = true
							} else {
								thisNode.UpdateNodeCost(g)
							}
						} else if !canBeRemoved && thisNode.CanBeRemoved {
							heap.Remove(g.AvailableMoves, thisNode.QueueIndex)
							thisNode.CanBeRemoved = false
						}
					}
				}
			}
		} else {
			neighbour := n.Neighbours[i]
			tmp := neighbour
			j := i
			for x := 0; x < 3; x++ {
				if tmp == nil {
					break
				}
				if tmp.IsInLoop && !tmp.IsVisited && !(tmp.IsDecided && !tmp.IsForRemoval && tmp.TemplateGroup == nil) {
					canBeRemoved, stopTesting = checkIfCanBeRemoved(tmp, g)
					if canBeRemoved {
						if !tmp.CanBeRemoved {
							tmp.SetNodeCost(g)
							heap.Push(g.AvailableMoves, tmp)
							tmp.CanBeRemoved = true
						} else {
							tmp.UpdateNodeCost(g)
						}
					} else if !canBeRemoved && tmp.CanBeRemoved {
						heap.Remove(g.AvailableMoves, tmp.QueueIndex)
						tmp.CanBeRemoved = false
					}
				}
				j = (j - 1 + 3) % 3
				tmp = tmp.Neighbours[j]
			}

			neighbour = n.Neighbours[(i+1)%3]
			tmp = neighbour
			j = (i + 1) % 3
			for x := 0; x < 3; x++ {
				if tmp == nil {
					break
				}
				if tmp.IsInLoop && !tmp.IsVisited && !(tmp.IsDecided && !tmp.IsForRemoval && tmp.TemplateGroup == nil) {
					canBeRemoved, stopTesting = checkIfCanBeRemoved(tmp, g)
					if canBeRemoved {
						if !tmp.CanBeRemoved {
							tmp.SetNodeCost(g)
							heap.Push(g.AvailableMoves, tmp)
							tmp.CanBeRemoved = true
						} else {
							tmp.UpdateNodeCost(g)
						}
					} else if !canBeRemoved && tmp.CanBeRemoved {
						heap.Remove(g.AvailableMoves, tmp.QueueIndex)
						tmp.CanBeRemoved = false
					}
				}
				j = (j + 1) % 3
				tmp = tmp.Neighbours[j]
			}
		}
	}

	return stopTesting
}

/* Main solver logic */
func loopSolveRecursion(n *utils.Node, g *utils.Graph, cost int, isSolutionFound *bool, isStateWrong *bool, depth int) {
	utils.NoVisitedStates++
	utils.AvgDepth += float64(depth)
	if depth > utils.MaxDepth {
		utils.MaxDepth = depth
	}
	debug.Println("")
	debug.Println("Node:")
	debug.Println(n)
	debug.Println("Cost:")
	debug.Println(cost)
	debug.Println("Depth:")
	debug.Println(depth)

	n.IsInLoop = false

	if debug.IsDebugMode {
		g.PrintBoard(true)
	}
	debug.Sleep(1000)

	/* Update list with available moves */
	stopTesting := updateAvailableMoves(n, g)
	if stopTesting {
		n.IsInLoop = true
		updateAvailableMoves(n, g)
		return
	}

	g.VisitedNodes.Push(nil)

	/* Select new move */
	for g.AvailableMoves.Len() > 0 {
		debug.Println("Priority queue:")
		for _, v := range *g.AvailableMoves {
			debug.Println(v)
		}
		newElement := heap.Pop(g.AvailableMoves)
		newNode := newElement.(*utils.Node)

		/* Return if it's not possible to remove the node */
		if newNode.IsDecided && !newNode.IsForRemoval {
			newNode.CanBeRemoved = false
			newNode.IsVisited = true
			g.VisitedNodes.Push(newNode)
			continue
		}

		/* Check if solution is found */
		if cost == newNode.Cost {
			newNode.IsInLoop = false
			debug.Println("SOLUTION FOUND")
			*isSolutionFound = true
			return
		}

		/* Delete move from available moves and save it in stack */
		newNode.CanBeRemoved = false
		newNode.IsVisited = true
		g.VisitedNodes.Push(newNode)

		if newNode.TemplateGroup != nil && !newNode.IsDecided {
			/* Set value for the group */
			if !newNode.TemplateGroup.SetValue(true, newNode, g) {
				/* Return if it's not possible to set this value */
				debug.Println("Template not aplicable")

				newNode.TemplateGroup.ClearValue(g)
				if !newNode.TemplateGroup.SetValue(false, newNode, g) {
					*isStateWrong = true
					debug.Println("WRONG STATE")
					debug.Println(depth)
					break
				} else {
					continue
				}
			}
		}

		/* Run recursion with new node */
		loopSolveRecursion(newNode, g, cost-newNode.Cost, isSolutionFound, isStateWrong, depth+1)

		if *isSolutionFound {
			return
		}

		/* If the solution without node wasn't found, then group must be in the loop */
		if newNode.TemplateGroup != nil && newNode.TemplateGroup.SettingNode == newNode {
			newNode.TemplateGroup.ClearValue(g)

			if !newNode.TemplateGroup.SetValue(false, newNode, g) {
				debug.Println("WRONG STATE")
				debug.Println(newNode)
				debug.Println(depth)
				break
			}
		}

		if *isStateWrong {
			debug.Println("BREAK STATE")
			debug.Println(depth)
			*isStateWrong = false
			break
		}
	}

	debug.Println("Clearing changes")

	for {
		thisElement := g.VisitedNodes.Pop()

		if thisElement == nil {
			break
		}

		thisNode := thisElement.(*utils.Node)

		debug.Println("Popped:")
		debug.Println(thisNode)

		/* Clear value of the group */
		if thisNode.IsVisited && thisNode.TemplateGroup != nil && thisNode.TemplateGroup.SettingNode == n {
			thisNode.TemplateGroup.ClearValue(g)
		}

		if thisNode.IsVisited && thisNode.TemplateGroup != nil && thisNode.IsDecided && !thisNode.IsForRemoval {
			thisNode.TemplateGroup.ClearValue(g)
		}
		thisNode.IsVisited = false
		thisNode.CanBeRemoved = true
		heap.Push(g.AvailableMoves, thisNode)
	}

	n.IsInLoop = true
	updateAvailableMoves(n, g)

	debug.Println("RETURN")
}

/* Solver preparation */
func LoopSolve(g *utils.Graph) {
	defer utils.TimeDuration(time.Now())
	debug.Println("START Loop Solver")

	cost := g.CalculateStartCost()
	if utils.IsTemplatesOn {
		g.FindTemplates()
	}

	g.CalculateStartMoves()

	g.VisitedNodes = stack.New()
	isSolutionFound := new(bool)
	isStateWrong := new(bool)

	if debug.IsDebugMode {
		g.PrintBoard(true)
	}
	debug.Println("Cost:")
	debug.Println(cost)

	for g.AvailableMoves.Len() > 0 {
		for _, v := range *g.AvailableMoves {
			debug.Println(v)
		}

		/* Choose new Node */
		newElement := heap.Pop(g.AvailableMoves)
		newNode := newElement.(*utils.Node)

		if newNode.IsDecided && !newNode.IsForRemoval {
			newNode.CanBeRemoved = false
			continue
		}

		/* Solution found */
		if cost == newNode.Cost {
			newNode.IsInLoop = false
			break
		}

		if newNode.TemplateGroup != nil && !newNode.IsDecided {
			if !newNode.TemplateGroup.SetValue(true, newNode, g) {
				debug.Println("Template not aplicable")

				newNode.TemplateGroup.ClearValue(g)
				newNode.TemplateGroup.SetValue(false, newNode, g)
				continue
			}
		}

		newNode.IsVisited = true
		newNode.CanBeRemoved = false

		/* Run recursion with new node */
		loopSolveRecursion(newNode, g, cost-newNode.Cost, isSolutionFound, isStateWrong, 1)

		if *isSolutionFound {
			break
		}

		if newNode.TemplateGroup != nil && newNode.TemplateGroup.SettingNode == newNode {
			debug.Println("Reversing set value")
			newNode.TemplateGroup.ClearValue(g)
			newNode.TemplateGroup.SetValue(false, newNode, g)
		}

		newNode.IsInLoop = true
	}
}

package solvers

import (
	"container/heap"
	"slitherlink_solver/debug"
	"slitherlink_solver/utils"

	"github.com/golang-collections/collections/stack"
)

/*
* Checks if node can be removed from the loop
* Second bool is true if we can't delete part of template because of the second rule
 */
func checkIfCanBeRemoved(n *utils.Node, g *utils.Graph) (bool, bool) {
	debug.Println("Checking if can be removed: ")

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

	// if !n.CanBeRemoved && n.TemplateGroup != nil && n.IsDecided && n.IsForRemoval {
	// 	n.TemplateGroup.Removable++
	// }

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

	/* Checking if is connected via corner*/
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

	// if n.Neighbours[0].Value == 1 && n.Neighbours[1].Value == 0 {
	// 	debug.Sleep(5000)
	// }

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
		if thisNode != nil {

			if thisNode.IsInLoop && !thisNode.IsVisited && !(thisNode.IsDecided && !thisNode.IsForRemoval) {
				canBeRemoved, stopTesting = checkIfCanBeRemoved(thisNode, g)
				if stopTesting {
					return true
				}
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
			if nextNode != nil && nextNode.IsInLoop && nextNode.CanBeRemoved && !nextNode.IsVisited && !(nextNode.IsDecided && !nextNode.IsForRemoval) {
				nextNode.UpdateNodeCost(g)
			}

			/* diagonal node */
			thisNode := thisNode.Neighbours[(i+1)%int(g.MaxDegree)]
			if thisNode != nil && thisNode.IsInLoop && !thisNode.IsVisited && !(thisNode.IsDecided && !thisNode.IsForRemoval) {
				canBeRemoved, stopTesting = checkIfCanBeRemoved(thisNode, g)
				if stopTesting {
					return true
				}
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

	return stopTesting
}

/* Main solver logic */
func loopSolveRecursion(n *utils.Node, g *utils.Graph, cost int, isSolutionFound *bool, depth int) {
	utils.NoVisitedStates++
	utils.AvgDepth += float32(depth)
	if depth > utils.MaxDepth {
		utils.MaxDepth = depth
	}
	debug.Println("")
	debug.Println("Node:")
	debug.Println(n)
	debug.Println("Cost:")
	debug.Println(cost)

	/* Update list with available moves */
	n.IsInLoop = false

	// if n.TemplateGroup != nil && !n.IsDecided {

	// 	if !n.TemplateGroup.SetValue(true, n, g) {
	// 		debug.Println("Template not aplicable")
	// 		if debug.IsDebugMode {
	// 			g.PrintSquaresBoard(true)
	// 		}
	// 		// debug.Sleep(100)

	// 		n.IsInLoop = true
	// 		n.TemplateGroup.ClearValue(g)
	// 		n.TemplateGroup.SetValue(false, n, g)
	// 		return
	// 	}

	// 			n.TemplateGroup.ClearValue(g)
	// 			n.TemplateGroup.SetValue(false, n, g)
	// 			g.PrintSquaresBoard(true)
	// 			// debug.Sleep(9999)

	// 			return
	// 		}
	// 	}
	// 	n.TemplateGroup.Removable--
	// 	n.TemplateGroup.Removed++
	// }

	stopTesting := updateAvailableMoves(n, g)
	if stopTesting {
		// fmt.Println("STOP TESTING")
		// if debug.IsDebugMode {
		// 	g.PrintSquaresBoard(true)
		// }
		// debug.Sleep(1000)

		// n.IsInLoop = true
		// updateAvailableMoves(n, g)
		// return
	}

	// if n.TemplateGroup != nil {
	// if n.TemplateGroup.Removable+n.TemplateGroup.Removed < n.TemplateGroup.Length {
	// debug.Sleep(9999)
	// n.TemplateGroup.ClearValue()
	// n.TemplateGroup.SetValue(false, n.TemplateGroup.SettingNode)
	// n.IsInLoop = true
	// return
	// }
	// }
	g.VisitedNodes.Push(nil)

	if debug.IsDebugMode {
		g.PrintSquaresBoard(true)
	}
	debug.Sleep(50)

	/* Select new move */
	for g.AvailableMoves.Len() > 0 {
		for _, v := range *g.AvailableMoves {
			debug.Println(v)
		}
		var newElement any
		var newNode *utils.Node
		// for {
		newElement = heap.Pop(g.AvailableMoves)
		newNode = newElement.(*utils.Node)

		// if !(newNode.IsDecided && !newNode.IsForRemoval) {

		// newNode.IsVisited = true
		// g.VisitedNodes.Push(newNode)
		// if newNode.CanBeRemoved {
		// newNode.CanBeRemoved = false
		// 	break
		// }
		// }

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

		/* Run recursion with new node */
		loopSolveRecursion(newNode, g, cost-newNode.Cost, isSolutionFound, depth+1)

		if *isSolutionFound {
			return
		}
		// } else {
		// 	newNode.CanBeRemoved = false
		// 	g.VisitedNodes.Push(newNode)
		// }
	}

	/* Clear changes */
	// if n.TemplateGroup != nil {
	// 	n.TemplateGroup.ClearValue(g)
	// }
	for {
		thisElement := g.VisitedNodes.Pop()
		if thisElement == nil {
			break
		}

		thisNode := thisElement.(*utils.Node)
		// if thisNode.IsVisited && thisNode.TemplateGroup != nil && thisNode.IsDecided && !thisNode.IsForRemoval {
		// 	thisNode.TemplateGroup.ClearValue(g)
		// }

		// if thisNode.IsVisited && thisNode.TemplateGroup != nil && thisNode.TemplateGroup.SettingNode == n {
		// 	thisNode.TemplateGroup.ClearValue(g)
		// }
		thisNode.IsVisited = false
		thisNode.CanBeRemoved = true
		heap.Push(g.AvailableMoves, thisNode)
	}

	n.IsInLoop = true
	updateAvailableMoves(n, g)

}

/* Solver preparation */
func LoopSolve(g *utils.Graph) {
	debug.Println("START Loop Solver")

	_, cost := g.CalculateStartCost()
	g.FindTemplates()

	g.CalculateStartMoves()

	g.VisitedNodes = stack.New()
	isSolutionFound := new(bool)

	if debug.IsDebugMode {
		g.PrintSquaresBoard(true)
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

		// if !(newNode.IsDecided && !newNode.IsForRemoval) {
		/* Solution found */
		if cost == newNode.Cost {
			newNode.IsInLoop = false
			break
		}

		newNode.IsVisited = true
		newNode.CanBeRemoved = false

		/* Run recursion with new node */
		loopSolveRecursion(newNode, g, cost-newNode.Cost, isSolutionFound, 1)

		if *isSolutionFound {
			break
		}

		newNode.IsInLoop = true
		// } else {
		// 	newNode.CanBeRemoved = false
		// }
	}

	g.PrintSquaresBoard(false)
}

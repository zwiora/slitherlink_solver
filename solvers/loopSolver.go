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

	if g.Shape != "honeycomb" {
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
		if thisNode != nil {

			if thisNode.IsInLoop && !thisNode.IsVisited {
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
					// g.PrintSquaresBoard(true)
					// fmt.Println(thisNode)
					heap.Remove(g.AvailableMoves, thisNode.QueueIndex)
					thisNode.CanBeRemoved = false
				}
			}

			/* neighbour of the neighbour - only updates cost*/
			nextNode := thisNode.Neighbours[(i)%int(g.MaxDegree)]
			if nextNode != nil && nextNode.IsInLoop && nextNode.CanBeRemoved && !nextNode.IsVisited {
				nextNode.UpdateNodeCost(g)
			}
			/* another neighbour of the neighbour - in case of honeycomb */
			nextNode = thisNode.Neighbours[(i+1)%int(g.MaxDegree)]
			if nextNode != nil && nextNode.IsInLoop && nextNode.CanBeRemoved && !nextNode.IsVisited {
				nextNode.UpdateNodeCost(g)
			}

			previousNode := n.Neighbours[(i-1+int(g.MaxDegree))%int(g.MaxDegree)]
			if previousNode == nil {
				nextNode = thisNode.Neighbours[(i-1+int(g.MaxDegree))%int(g.MaxDegree)]
				if nextNode != nil && nextNode.IsInLoop && nextNode.CanBeRemoved && !nextNode.IsVisited {
					nextNode.UpdateNodeCost(g)
				}
			}

			/* diagonal node */
			if g.Shape != "honeycomb" {
				thisNode := thisNode.Neighbours[(i+1)%int(g.MaxDegree)]
				if thisNode != nil && thisNode.IsInLoop && !thisNode.IsVisited {
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
				newNode.TemplateGroup.SetValue(false, newNode, g)
				continue
			}
		}

		/* Run recursion with new node */
		loopSolveRecursion(newNode, g, cost-newNode.Cost, isSolutionFound, depth+1)

		if *isSolutionFound {
			return
		}

		/* If the solution withut node wasn't found, then group must be in the loop */
		if newNode.TemplateGroup != nil && newNode.TemplateGroup.SettingNode == newNode {
			newNode.TemplateGroup.ClearValue(g)
			newNode.TemplateGroup.SetValue(false, newNode, g)
		}
	}

	for {
		debug.Println("Clearing changes")
		thisElement := g.VisitedNodes.Pop()
		if thisElement == nil {
			break
		}

		thisNode := thisElement.(*utils.Node)

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

}

/* Solver preparation */
func LoopSolve(g *utils.Graph) {
	debug.Println("START Loop Solver")

	cost := g.CalculateStartCost()
	// g.FindTemplates()

	g.CalculateStartMoves()

	g.PrintBoard(true)

	g.VisitedNodes = stack.New()
	isSolutionFound := new(bool)

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
		loopSolveRecursion(newNode, g, cost-newNode.Cost, isSolutionFound, 1)

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

	g.PrintBoard(false)
}

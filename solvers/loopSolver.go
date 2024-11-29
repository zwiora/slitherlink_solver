package solvers

import (
	"container/heap"
	"slitherlink_solver/debug"
	"slitherlink_solver/utils"

	"github.com/golang-collections/collections/stack"
)

/* Checks if node can be removed from the loop */
func checkIfCanBeRemoved(n *utils.Node, g *utils.Graph) bool {
	debug.Println("Checking if can be removed: ")

	/* Checking if removal would create loop inside the loop */
	nodeDegree := n.GetDegree()
	if nodeDegree == 0 || nodeDegree == int(g.MaxDegree) {
		debug.Println("\t- SKIP: removal would create loop inside the loop")
		return false
	}
	debug.Println("\t- degree ok")

	/* Checking if neighbour would have enough edges*/
	for _, v := range n.Neighbours {
		if v != nil && !v.IsInLoop {
			if v.GetDegree() <= int(v.Value) {
				debug.Println("\t- SKIP: deletion isn't in the solution")
				n.IsInLoop = true
				return false
			}
		}
	}

	/* It can be removed if it's a leaf */
	if n.GetDegree() == 1 {
		debug.Println("\t- it's a leaf")
		return true
	}
	debug.Println("\t- not a leaf")

	/* Checking if its neighbour is leaf */
	for _, v := range n.Neighbours {
		if v != nil && v.IsInLoop {
			if v.GetDegree() == 1 {
				debug.Println("\t- SKIP: neighbour is a leaf")
				return false
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
			return false
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
					return false
				}
			}
		}
	}

	/* No troubles found - node can be deleted */
	return true
}

/* Updates heap of nodes that can be removed from the loop - checks if can be removed and its cost*/
func updateAvailableMoves(n *utils.Node, g *utils.Graph) {
	for i := 0; i < int(g.MaxDegree); i++ {
		/* neighbouring node */
		thisNode := n.Neighbours[(i)%int(g.MaxDegree)]
		if thisNode != nil {

			if thisNode.IsInLoop && !thisNode.IsVisited && !(thisNode.IsDecided && !thisNode.IsForRemoval) {
				canBeRemoved := checkIfCanBeRemoved(thisNode, g)
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
				canBeRemoved := checkIfCanBeRemoved(thisNode, g)
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
	g.VisitedNodes.Push(nil)
	updateAvailableMoves(n, g)

	if debug.IsDebugMode {
		g.PrintSquaresBoard(true)
	}
	debug.Sleep(1000)

	/* Select new move */
	for g.AvailableMoves.Len() > 0 {
		for _, v := range *g.AvailableMoves {
			debug.Println(v)
		}
		var newElement any
		var newNode *utils.Node
		for {
			newElement = heap.Pop(g.AvailableMoves)
			newNode = newElement.(*utils.Node)

			newNode.IsVisited = true
			g.VisitedNodes.Push(newNode)

			if newNode.CanBeRemoved {
				newNode.CanBeRemoved = false
				break
			}
		}

		/* Check if solution is found */
		if cost == newNode.Cost {
			newNode.IsInLoop = false
			debug.Println("SOLUTION FOUND")
			*isSolutionFound = true
			return
		}

		/* Delete move from available moves and save it in stack */
		// newNode.CanBeRemoved = false
		// newNode.IsVisited = true
		// g.VisitedNodes.Push(newNode)

		/* Run recursion with new node */
		loopSolveRecursion(newNode, g, cost-newNode.Cost, isSolutionFound, depth+1)

		if *isSolutionFound {
			return
		}
	}

	/* Clear changes */
	for {
		thisElement := g.VisitedNodes.Pop()
		if thisElement == nil {
			break
		}

		thisNode := thisElement.(*utils.Node)
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
	}

	g.PrintSquaresBoard(false)
}

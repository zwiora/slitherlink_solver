package solvers

import (
	"container/heap"
	"slytherlink_solver/debug"
	"slytherlink_solver/utils"

	"github.com/golang-collections/collections/stack"
)

/* Checks if node can be removed from the loop */
func checkIfCanBeRemoved(n *utils.Node, g *utils.Graph) bool {
	debug.Print("Checking if can be removed: ")

	/* Checking if removal would create loop inside the loop */
	nodeDegree := n.GetDegree()
	if nodeDegree == 0 || nodeDegree == int(g.MaxDegree) {
		debug.Print("\t- SKIP: removal would create loop inside the loop")
		return false
	}
	debug.Print("\t- degree ok")

	/* It can be removed if it's a leaf */
	if n.GetDegree() == 1 {
		debug.Print("\t- it's a leaf")
		return true
	}
	debug.Print("\t- not a leaf")

	/* Checking if its neighbour is leaf */
	for _, v := range n.Neighbours {
		if v != nil && v.IsInLoop {
			if v.GetDegree() == 1 {
				debug.Print("\t- SKIP: neighbour is a leaf")
				return false
			}
		}
	}
	debug.Print("\t- neighbour isn't a leaf")

	/* Checking if is between two sides of graph*/
	sidesCounter := 0
	for i := 0; i < len(n.Neighbours); i++ {
		thisNeighbour := n.Neighbours[i]
		nextNeighbour := n.Neighbours[(i+1)%int(g.MaxDegree)]
		if ((thisNeighbour == nil || !thisNeighbour.IsInLoop) && (nextNeighbour != nil && nextNeighbour.IsInLoop)) || ((nextNeighbour == nil || !nextNeighbour.IsInLoop) && (thisNeighbour != nil && thisNeighbour.IsInLoop)) {
			sidesCounter++
		}
		if sidesCounter == 3 {
			debug.Print("\t- SKIP: deletion would create two separate graphs")
			return false
		}
	}
	debug.Print("\t- deletion wouldn't create two separate graphs")

	/* Checking if is connected via corner*/
	for k, v := range n.Neighbours {
		if v != nil && v.IsInLoop {
			diagonalNode := v.Neighbours[(k+1)%int(g.MaxDegree)]
			if diagonalNode != nil && !diagonalNode.IsInLoop {
				nextNeighbour := diagonalNode.Neighbours[(k+2)%int(g.MaxDegree)]
				if nextNeighbour.IsInLoop {

					debug.Print("\t- SKIP: deletion would create two graphs with common corner")
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

			if thisNode.IsInLoop && !thisNode.IsVisited {
				canBeRemoved := checkIfCanBeRemoved(thisNode, g)
				if canBeRemoved {
					if !thisNode.CanBeRemoved {
						thisNode.Cost = thisNode.CalculateNodeCost(g)
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
			if nextNode != nil && nextNode.IsInLoop && nextNode.CanBeRemoved && !nextNode.IsVisited {
				nextNode.UpdateNodeCost(g)
			}

			/* diagonal node */
			thisNode := thisNode.Neighbours[(i+1)%int(g.MaxDegree)]
			if thisNode != nil && thisNode.IsInLoop && !thisNode.IsVisited {
				canBeRemoved := checkIfCanBeRemoved(thisNode, g)
				if canBeRemoved {
					if !thisNode.CanBeRemoved {
						thisNode.Cost = thisNode.CalculateNodeCost(g)
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
func loopSolveRecursion(n *utils.Node, g *utils.Graph, cost int, isSolutionFound *bool) {
	debug.Print("")
	debug.Print("Node:")
	debug.Print(n)
	debug.Print("Cost:")
	debug.Print(cost)

	/* Update list with available moves */
	n.IsInLoop = false
	g.VisitedNodes.Push(nil)
	updateAvailableMoves(n, g)

	debug.PrintBoard(g)
	debug.Sleep(1000)

	/* Select new move */
	for g.AvailableMoves.Len() > 0 {
		// for _, v := range *g.AvailableMoves {
		// 	debug.Print(v)
		// }
		newElement := heap.Pop(g.AvailableMoves)
		newNode := newElement.(*utils.Node)

		/* Check if solution is found */
		if cost == newNode.Cost {
			newNode.IsInLoop = false
			debug.Print("SOLUTION FOUND")
			*isSolutionFound = true
			return
		}

		/* Delete move from available moves and save it in stack */
		newNode.CanBeRemoved = false
		newNode.IsVisited = true
		g.VisitedNodes.Push(newNode)

		/* Run recursion with new node */
		loopSolveRecursion(newNode, g, cost-newNode.Cost, isSolutionFound)

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
func LoopSolve(g *utils.Graph, isCheckingAllSolutions bool) {
	debug.Print("START Loop Solver")

	_, cost := g.CalculateStartCost()
	g.CalculateStartMoves()
	g.VisitedNodes = stack.New()
	isSolutionFound := new(bool)

	debug.PrintBoard(g)
	debug.Print("Cost:")
	debug.Print(cost)

	for g.AvailableMoves.Len() > 0 {

		/* Choose new Node */
		newElement := heap.Pop(g.AvailableMoves)
		newNode := newElement.(*utils.Node)

		if cost == newNode.Cost {
			newNode.IsInLoop = false
			break
		}

		newNode.IsVisited = true
		newNode.CanBeRemoved = false

		/* Run recursion with new node */
		loopSolveRecursion(newNode, g, cost-newNode.Cost, isSolutionFound)

		if *isSolutionFound {
			break
		}

		newNode.IsInLoop = true
	}

	g.PrintSquaresBoard(false)
}

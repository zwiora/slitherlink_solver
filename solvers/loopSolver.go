package solvers

import (
	"fmt"
	"slytherlink_solver/debug"
	"slytherlink_solver/utils"

	"github.com/golang-collections/collections/stack"
)

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

	return true
}

func updateAvailableMoves(n *utils.Node, g *utils.Graph) {
	// update list with available moves
	for i := 0; i < int(g.MaxDegree); i++ {
		/* neighbouring node */
		thisNode := n.Neighbours[(i)%int(g.MaxDegree)]
		if thisNode != nil {

			if thisNode.IsInLoop && !thisNode.IsVisited {
				canBeRemoved := checkIfCanBeRemoved(thisNode, g)
				if canBeRemoved && !thisNode.CanBeRemoved {
					g.AvailableMoves.PushBack(thisNode)
					thisNode.CanBeRemoved = true
				} else if !canBeRemoved && thisNode.CanBeRemoved {
					for e := g.AvailableMoves.Front(); e != nil; e = e.Next() {
						if e.Value == thisNode {
							g.AvailableMoves.Remove(e)
							break
						}
					}
					thisNode.CanBeRemoved = false
				}
			}

			/* diagonal node */
			thisNode := thisNode.Neighbours[(i+1)%int(g.MaxDegree)]
			if thisNode != nil && thisNode.IsInLoop && !thisNode.IsVisited {
				canBeRemoved := checkIfCanBeRemoved(thisNode, g)
				if canBeRemoved && !thisNode.CanBeRemoved {
					g.AvailableMoves.PushBack(thisNode)
					thisNode.CanBeRemoved = true
				} else if !canBeRemoved && thisNode.CanBeRemoved {
					for e := g.AvailableMoves.Front(); e != nil; e = e.Next() {
						if e.Value == thisNode {
							g.AvailableMoves.Remove(e)
							break
						}
					}
					thisNode.CanBeRemoved = false
				}
			}
		}
	}
}

func loopSolveRecursion(n *utils.Node, g *utils.Graph, cost int, isSolutionFound *bool) {
	/* Calculate new cost */
	debug.Print("Old cost:")
	debug.Print(cost)

	newCost := cost

	if n.Value != -1 {
		newCost -= n.GetCostOfField(int(g.MaxDegree))
	}

	for _, v := range n.Neighbours {
		if v != nil && v.Value != -1 {
			newCost -= v.GetCostOfField(int(g.MaxDegree))
		}
	}

	n.IsInLoop = false

	if n.Value != -1 {
		newCost += n.GetCostOfField(int(g.MaxDegree))
	}

	for _, v := range n.Neighbours {
		if v != nil && v.Value != -1 {
			newCost += v.GetCostOfField(int(g.MaxDegree))
		}
	}

	debug.Print("New cost:")
	debug.Print(newCost)

	if newCost == 0 {
		debug.Print("SOLUTION FOUND")
		fmt.Println(newCost)
		// g.PrintSquaresBoard(false)
		*isSolutionFound = true
		return
	}

	g.VisitedNodes.Push(nil)

	updateAvailableMoves(n, g)

	debug.PrintBoard(g)
	debug.Print(g.AvailableMoves.Len())

	debug.Sleep(1000)

	for {
		thisElement := g.AvailableMoves.Front()
		if thisElement == nil {
			break
		}

		thisNode := thisElement.Value.(*utils.Node)

		/* Delete move from options and save in stack */

		thisNode.CanBeRemoved = false
		g.AvailableMoves.Remove(thisElement)
		thisNode.IsVisited = true
		g.VisitedNodes.Push(thisNode)

		loopSolveRecursion(thisNode, g, newCost, isSolutionFound)

		if *isSolutionFound {
			return
		}

		thisNode.IsInLoop = true
	}

	for {
		thisElement := g.VisitedNodes.Pop()
		if thisElement == nil {
			break
		}

		thisNode := thisElement.(*utils.Node)
		thisNode.IsVisited = false
		thisNode.CanBeRemoved = true
		g.AvailableMoves.PushFront(thisNode)
	}

	n.IsInLoop = true
	updateAvailableMoves(n, g)

}

func LoopSolve(g *utils.Graph, isCheckingAllSolutions bool) {
	_, cost := g.CalculateCost()
	g.CalculateStartingMoves()
	// debug.PrintBoard(g)

	g.VisitedNodes = stack.New()
	debug.Print(g.VisitedNodes.Len())

	isSolutionFound := new(bool)

	for {
		thisElement := g.AvailableMoves.Front()
		if thisElement == nil {
			break
		}

		thisNode := thisElement.Value.(*utils.Node)
		thisNode.IsVisited = true
		thisNode.CanBeRemoved = false
		g.AvailableMoves.Remove(thisElement)

		loopSolveRecursion(thisNode, g, cost, isSolutionFound)

		if *isSolutionFound {
			break
		}

		thisNode.IsInLoop = true
	}

	g.PrintSquaresBoard(false)
}

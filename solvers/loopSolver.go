package solvers

import (
	"fmt"
	"slytherlink_solver/debug"
	"slytherlink_solver/utils"
)

func loopSolveRecursion(n *utils.Node, g *utils.Graph, cost int) {
	n.IsVisited = true

	debug.Print("New repetition, testing cell:")
	debug.Print(n)
	debug.Print("Board state:")
	debug.PrintBoard(g)
	// debug.Sleep(1)

	var newNode *utils.Node
	isNewFound := false
	for _, v := range n.Neighbours {
		if v != nil && !v.IsVisited && v.IsInLoop {
			newNode = v
			isNewFound = true
			break
		}
	}

	if isNewFound {
		debug.Print("Checking variant without deletion")
		loopSolveRecursion(newNode, g, cost)
		newNode.IsVisited = false
	}

	nodeDegree := n.GetDegree()

	/* Checking if removal would create loop inside the loop */
	debug.Print("Checking if can be removed: ")
	if nodeDegree != 0 && nodeDegree != int(g.MaxNeighbourCount) {
		debug.Print("\t- degree ok")

		isBridge := false

		/* It's not a bridge if  it's a leaf */
		if n.GetDegree() > 1 {

			debug.Print("\t- not a leaf")

			/* Checking if its neighbour is leaf */
			for _, v := range n.Neighbours {
				if v != nil && v.IsInLoop {
					if v.GetDegree() == 1 {
						isBridge = true
						debug.Print("\t- SKIP: neighbour is a leaf")
						break
					}
				}
			}

			/* Checking if is between two sides of graph*/
			if !isBridge {
				debug.Print("\t- neighbour isn't a leaf")
				sidesCounter := 0
				for i := 0; i < len(n.Neighbours); i++ {
					// fmt.Println(i, " ", sidesCounter)
					thisNeighbour := n.Neighbours[i]
					nextNeighbour := n.Neighbours[(i+1)%int(g.MaxNeighbourCount)]
					if ((thisNeighbour == nil || !thisNeighbour.IsInLoop) && (nextNeighbour != nil && nextNeighbour.IsInLoop)) || ((nextNeighbour == nil || !nextNeighbour.IsInLoop) && (thisNeighbour != nil && thisNeighbour.IsInLoop)) {
						sidesCounter++
					}
					if sidesCounter == 3 {
						isBridge = true
						debug.Print("\t- SKIP: deletion would create two separate graphs")
						break
					}
				}
			}

			/* Checking if is connected via corner*/
			if !isBridge {
				debug.Print("\t- deletion wouldn't create two separate graphs")
				for k, v := range n.Neighbours {
					if v != nil && v.IsInLoop {
						diagonalNode := v.Neighbours[(k+1)%int(g.MaxNeighbourCount)]
						if diagonalNode != nil && !diagonalNode.IsInLoop {
							nextNeighbour := diagonalNode.Neighbours[(k+2)%int(g.MaxNeighbourCount)]
							if nextNeighbour.IsInLoop {
								isBridge = true
								debug.Print("\t- SKIP: deletion would create two graphs with common corner")
								break
							}
						}
					}
				}
			}
		}

		/* Removing node from the loop*/
		if !isBridge {

			debug.Print("\t- deletion possible")

			/* Calculating new cost */
			debug.Print("Old cost:")
			debug.Print(cost)

			newCost := cost

			if n.Value != -1 {
				newCost -= n.GetCostOfField(int(g.MaxNeighbourCount))
			}

			for _, v := range n.Neighbours {
				if v != nil && v.Value != -1 {
					newCost -= v.GetCostOfField(int(g.MaxNeighbourCount))
				}
			}

			n.IsInLoop = false

			if n.Value != -1 {
				newCost += n.GetCostOfField(int(g.MaxNeighbourCount))
			}

			for _, v := range n.Neighbours {
				if v != nil && v.Value != -1 {
					newCost += v.GetCostOfField(int(g.MaxNeighbourCount))
				}
			}

			debug.Print("New cost:")
			debug.Print(newCost)

			if newCost == 0 {
				debug.Print("SOLUTION FOUND")
				g.PrintSquaresBoard(false)
			}

			if isNewFound {
				newNode.IsVisited = false
				debug.Print("Checking variant with deletion")
				loopSolveRecursion(newNode, g, newCost)
				newNode.IsVisited = false
				newNode.IsInLoop = true
			}
		}
	}

	n.IsVisited = false
	n.IsInLoop = true

}

func LoopSolve(g *utils.Graph, isCheckingAllSolutions bool) {

	// isSolutionFound := false
	// g.PrintSquaresBoard()
	_, cost := g.CalculateCost()
	fmt.Println(cost)

	loopSolveRecursion(g.Root, g, cost)

	// lastLineNode := g.Root
	// thisNode := g.Root

	// for n := 0; n < g.SizeY; n++ {
	// 	for m := 0; m < g.SizeX; m++ {
	// 		if thisNode.IsInLoop {
	// 			fmt.Printf("\033[42m")
	// 		}
	// 		if thisNode.Value == -1 {
	// 			fmt.Printf("   \033[49m|")
	// 		} else {
	// 			fmt.Printf(" %d \033[49m|", thisNode.Value)
	// 		}
	// 		thisNode = thisNode.Neighbours[0]
	// 	}
	// 	fmt.Println()
	// 	fmt.Printf("-")

	// 	for m := 0; m < g.SizeX; m++ {
	// 		fmt.Printf("----")
	// 	}
	// 	fmt.Println()
	// 	thisNode = lastLineNode.Neighbours[1]
	// 	lastLineNode = thisNode
	// }

	// for true {

	// 	thisNode := g.Root

	// 	if cost == 0 {
	// 		isSolutionFound = true
	// 	}

	// 	if !isCheckingAllSolutions && isSolutionFound {
	// 		break
	// 	}
	// }
}

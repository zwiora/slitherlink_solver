package solvers

import (
	"slytherlink_solver/utils"
	"time"
)

func loopSolveRecursion(n *utils.Node, g *utils.Graph, cost int) {
	// fmt.Println("New repetition", n)
	nodeDegree := n.GetDegree()

	/* Checking if removal would create loop inside the loop */
	if nodeDegree != 0 && nodeDegree != int(g.MaxNeighbourCount) {
		// fmt.Println("Degree ok")

		isBridge := false

		/* It's not a bridge if  it's a leaf */
		if n.GetDegree() > 1 {

			// fmt.Println("Not a leaf")

			/* Checking if its neighbour is leaf */
			for _, v := range n.Neighbours {
				if v != nil && v.IsInLoop {
					if v.GetDegree() == 1 {
						isBridge = true
						break
					}
				}
			}

			/* Checking if is between two sides of graph*/
			if !isBridge {
				// fmt.Println("Checking sides")
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
						break
					}
				}
			}

			/* Checking if is connected via edge*/
			if !isBridge {
				// fmt.Println("Checking diagonals")
				for k, v := range n.Neighbours {
					if v != nil && v.IsInLoop {
						diagonalNode := v.Neighbours[(k+1)%int(g.MaxNeighbourCount)]
						if diagonalNode != nil && !diagonalNode.IsInLoop {
							nextNeighbour := diagonalNode.Neighbours[(k+2)%int(g.MaxNeighbourCount)]
							if nextNeighbour.IsInLoop {
								isBridge = true
								break
							}
						}
					}
				}
			}
		}

		/* Removing node from the loop*/
		if !isBridge {

			// fmt.Println("Graph ok")

			/* Calculating new cost */
			newCost := cost

			if n.Value != -1 {
				newCost -= n.GetCostOfField(int(g.MaxNeighbourCount))
			}

			for _, v := range n.Neighbours {
				if v != nil && v.Value != -1 {
					newCost -= v.GetCostOfField(int(g.MaxNeighbourCount))
				}
			}

			n.IsVisited = true
			n.IsInLoop = false

			if n.Value != -1 {
				newCost += n.GetCostOfField(int(g.MaxNeighbourCount))
			}

			for _, v := range n.Neighbours {
				if v != nil && v.Value != -1 {
					newCost += v.GetCostOfField(int(g.MaxNeighbourCount))
				}
			}

			// if newCost == 0 {
			g.PrintSquaresBoard()
			// }

			time.Sleep(250 * time.Millisecond)
			// fmt.Println(newCost)

			for i := 0; i < len(n.Neighbours); i++ {
				v := n.Neighbours[i]
				if v != nil && !v.IsVisited && v.IsInLoop {
					// fmt.Println("i: ", i)
					loopSolveRecursion(v, g, newCost)
					// fmt.Println("Powrót")
					v.IsVisited = false
					v.IsInLoop = true
				}
			}
		}

	}

}

func LoopSolve(g *utils.Graph, isCheckingAllSolutions bool) {

	// isSolutionFound := false
	_, cost := g.CalculateCost()

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

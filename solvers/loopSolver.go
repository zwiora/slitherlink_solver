package solvers

import (
	"fmt"
	"slytherlink_solver/utils"
)

func loopSolveRecursion(n *utils.Node, g *utils.Graph, cost int) {
	fmt.Println("New repetition", n)

	/* Checking if removal would create loop inside the loop */
	if n.GetDegree() != int(g.MaxNeighbourCount) {
		fmt.Println("Degree ok")

		/* Checking if removal would break the graph */
		isBridge := false
		for _, v := range n.Neighbours {
			if v != nil && v.IsInLoop {
				if v.GetDegree() == 1 {
					isBridge = true
					break
				}

			}
		}

		if !isBridge || n.GetDegree() == 1 {

			fmt.Println("Graph ok")

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

			g.PrintSquaresBoard()
			fmt.Println(newCost)

			for _, v := range n.Neighbours {
				if v != nil && !v.IsVisited {
					loopSolveRecursion(v, g, newCost)
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

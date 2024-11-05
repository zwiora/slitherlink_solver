package solvers

import (
	"fmt"
	"slytherlink_solver/utils"
)

func LoopSolve(g *utils.Graph, isCheckingAllSolutions bool) {

	// isSolutionFound := false
	f, cost := g.CalculateCost()
	fmt.Println(f, " ", cost)
	// for true {

	// 	thisNode := g.Root

	// 	if !isCheckingAllSolutions && isSolutionFound {
	// 		break
	// 	}
	// }
}

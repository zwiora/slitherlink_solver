package main

import (
	"slytherlink_solver/solvers"
	"slytherlink_solver/utils"
)

func main() {
	g := utils.ConstructBoardFromData("data/test5.sav")
	// g.PrintSquaresBoard()
	// fmt.Println(g.CalculateCost())
	solvers.LoopSolve(g, false)
}

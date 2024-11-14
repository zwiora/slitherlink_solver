package main

import (
	"os"
	"slytherlink_solver/solvers"
	"slytherlink_solver/utils"
)

func main() {
	args := os.Args
	if len(args) > 1 && args[1] == "d" {
		utils.IsDebugMode = true
	}
	g := utils.ConstructBoardFromData("data/test3.sav")
	// g.PrintSquaresBoard()
	// fmt.Println(g.CalculateCost())
	solvers.LoopSolve(g, false)
}

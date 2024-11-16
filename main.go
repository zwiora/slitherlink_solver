package main

import (
	"os"
	"slytherlink_solver/debug"
	"slytherlink_solver/solvers"
	"slytherlink_solver/utils"
)

func main() {
	args := os.Args
	if len(args) > 1 && args[1] == "d" {
		debug.IsDebugMode = true
	}
	g := utils.ConstructBoardFromData("data/test5.sav")
	// g.PrintSquaresBoard()
	// fmt.Println(g.CalculateCost())
	solvers.LoopSolve(g, false)
}

package main

import (
	"os"
	"slytherlink_solver/debug"
	"slytherlink_solver/solvers"
	"slytherlink_solver/utils"
)

func main() {
	args := os.Args
	if len(args) > 2 && args[2] == "d" {
		debug.IsDebugMode = true
	}
	g := utils.ConstructBoardFromData("data/test" + args[1] + ".sav")
	// g.PrintSquaresBoard()
	// fmt.Println(g.CalculateCost())
	solvers.LoopSolve(g, false)
}

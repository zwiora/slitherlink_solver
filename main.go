package main

import (
	"os"
	"slytherlink_solver/debug"
	"slytherlink_solver/solvers"
	"slytherlink_solver/utils"
)

func main() {
	args := os.Args
	if len(args) > 3 && args[3] == "d" {
		debug.IsDebugMode = true
	}
	g := utils.ConstructBoardFromData("data/test" + args[2] + ".sav")

	if args[1] == "on" {
		utils.IsHeuristicOn = true
	}
	solvers.LoopSolve(g)
}

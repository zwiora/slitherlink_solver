package main

import (
	"fmt"
	"os"
	"slitherlink_solver/debug"
	"slitherlink_solver/solvers"
	"slitherlink_solver/utils"
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

	utils.AvgDepth /= float32(utils.NoVisitedStates)
	fmt.Println("Visited states: ", utils.NoVisitedStates)
	fmt.Println("Average depth: ", utils.AvgDepth)
	fmt.Println("Max depth: ", utils.MaxDepth)
}

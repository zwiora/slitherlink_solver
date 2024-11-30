package main

import (
	"fmt"
	"os"
	"slitherlink_solver/debug"
	"slitherlink_solver/solvers"
	"slitherlink_solver/utils"
	"strconv"
)

func main() {
	args := os.Args

	if args[1] == "on" {
		utils.IsHeuristicOn = true
	}

	if len(args) > 3 && args[3] == "d" {
		debug.IsDebugMode = true
	}

	data := utils.ReadMultiplePuzzleStructure("data/multiple.txt")
	boardType := ""
	sizeX := 0
	sizeY := 0
	code := ""
	var err error

	if len(args) > 2 {
		i, err := strconv.Atoi(args[2])
		utils.Check(err)
		boardType = data[i][0]
		sizeX, err = strconv.Atoi(data[i][1])
		utils.Check(err)
		sizeY, err = strconv.Atoi(data[i][2])
		utils.Check(err)
		code = data[i][3]

		g := utils.ConstructBoardFromData(boardType, sizeX, sizeY, code)
		solvers.LoopSolve(g)

		utils.AvgDepth /= float32(utils.NoVisitedStates)
		fmt.Println("Visited states: ", utils.NoVisitedStates)
		fmt.Println("Average depth: ", utils.AvgDepth)
		fmt.Println("Max depth: ", utils.MaxDepth)
		fmt.Println(g.CheckIfSolutionOk())
		g.PrintSquaresBoard(true)
		fmt.Println()
	} else {
		for i := range data {
			fmt.Println(i, ": ", data[i])
			boardType = data[i][0]
			sizeX, err = strconv.Atoi(data[i][1])
			utils.Check(err)
			sizeY, err = strconv.Atoi(data[i][2])
			utils.Check(err)
			code = data[i][3]

			g := utils.ConstructBoardFromData(boardType, sizeX, sizeY, code)

			solvers.LoopSolve(g)

			utils.AvgDepth /= float32(utils.NoVisitedStates)
			fmt.Println("Visited states: ", utils.NoVisitedStates)
			fmt.Println("Average depth: ", utils.AvgDepth)
			fmt.Println("Max depth: ", utils.MaxDepth)
			fmt.Println(g.CheckIfSolutionOk())
			fmt.Println()

			if !g.CheckIfSolutionOk() {
				fmt.Println("BŁĄD")
				break
			}
		}
	}

	// g := utils.ConstructBoardFromData("data/test" + args[2] + ".sav")

}

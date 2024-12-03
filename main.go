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
	dataFile := ""

	if args[1] == "s" {
		dataFile = "multiple"
	} else if args[1] == "h" {
		dataFile = "hexagon"
	} else {
		return
	}

	if args[2] == "on" {
		utils.IsHeuristicOn = true
	}

	if len(args) > 4 && args[4] == "d" {
		debug.IsDebugMode = true
	}

	data := utils.ReadMultiplePuzzleStructure("data/" + dataFile + ".txt")
	boardType := ""
	sizeX := 0
	sizeY := 0
	code := ""
	var err error

	if len(args) > 3 {
		i, err := strconv.Atoi(args[3])
		utils.Check(err)
		fmt.Println(i, ": ", data[i])

		boardType = data[i][0]
		sizeX, err = strconv.Atoi(data[i][1])
		utils.Check(err)
		sizeY, err = strconv.Atoi(data[i][2])
		utils.Check(err)
		code = data[i][3]

		g := utils.ConstructBoardFromData(boardType, sizeX, sizeY, code)
		fmt.Println(g)
		g.PrintBoard(debug.IsDebugMode)

		solvers.LoopSolve(g)

		g.PrintBoard(true)

		return

		utils.AvgDepth /= float32(utils.NoVisitedStates)
		fmt.Println("Visited states: ", utils.NoVisitedStates)
		fmt.Println("Average depth: ", utils.AvgDepth)
		fmt.Println("Max depth: ", utils.MaxDepth)
		fmt.Println(g.CheckIfSolutionOk())
		g.PrintBoard(true)
		fmt.Println()
	} else {
		for i := range data {
			utils.AvgDepth = 0
			utils.MaxDepth = 0
			utils.NoVisitedStates = 0
			fmt.Println(i, ": ", data[i])
			boardType = data[i][0]
			sizeX, err = strconv.Atoi(data[i][1])
			utils.Check(err)
			sizeY, err = strconv.Atoi(data[i][2])
			utils.Check(err)
			code = data[i][3]

			g := utils.ConstructBoardFromData(boardType, sizeX, sizeY, code)
			solvers.LoopSolve(g)
			g.PrintBoard(true)

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

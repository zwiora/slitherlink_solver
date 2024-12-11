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
		dataFile = "square_hard_short"
	} else if args[1] == "h" {
		dataFile = "hexagon"
	} else if args[1] == "t" {
		dataFile = "triangle"
	} else {
		return
	}

	if args[2] == "on" {
		utils.IsHeuristicOn = true

		var err error
		utils.HeuristicType, err = strconv.Atoi(args[3])
		utils.Check(err)

		if utils.HeuristicType > 0 {
			utils.IsTemplatesOn = true
		}
	}

	if len(args) > 5 && args[5] == "d" {
		debug.IsDebugMode = true
	}

	data := utils.ReadMultiplePuzzleStructure("data/" + dataFile + ".txt")
	boardType := ""
	sizeX := 0
	sizeY := 0
	code := ""
	var err error

	sum := 0

	if len(args) > 4 {
		i, err := strconv.Atoi(args[4])
		utils.Check(err)

		boardType = data[i][0]
		sizeX, err = strconv.Atoi(data[i][1])
		utils.Check(err)
		sizeY, err = strconv.Atoi(data[i][2])
		utils.Check(err)
		code = data[i][3]
		g := utils.ConstructBoardFromData(boardType, sizeX, sizeY, code)
		g.PrintBoard(false)
		fmt.Println()

		// return
		solvers.LoopSolve(g)
		fmt.Println()
		g.PrintBoard(true)
		fmt.Println()

		utils.AvgDepth /= float64(utils.NoVisitedStates)
		fmt.Println("Visited states: ", utils.NoVisitedStates)
		fmt.Println("Average depth: ", utils.AvgDepth)
		fmt.Println("Max depth: ", utils.MaxDepth)
		fmt.Println(g.CheckIfSolutionOk())
		fmt.Println()
	} else {
		for i := range data {
			utils.AvgDepth = 0
			utils.MaxDepth = 0
			utils.NoVisitedStates = 0
			boardType = data[i][0]
			sizeX, err = strconv.Atoi(data[i][1])
			utils.Check(err)
			sizeY, err = strconv.Atoi(data[i][2])
			utils.Check(err)
			code = data[i][3]

			g := utils.ConstructBoardFromData(boardType, sizeX, sizeY, code)
			fmt.Print(g.FieldsCount, ";")

			solvers.LoopSolve(g)

			utils.AvgDepth /= float64(utils.NoVisitedStates)
			fmt.Print(";", utils.NoVisitedStates)
			fmt.Print(";", utils.AvgDepth)
			fmt.Print(";", utils.MaxDepth)
			fmt.Println()
			g.CheckIfSolutionOk()

			sum += utils.NoVisitedStates

			if !g.CheckIfSolutionOk() {
				fmt.Println("WRONG SOLUTION")
				break
			}
		}
	}
}

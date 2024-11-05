package main

import (
	"slytherlink_solver/utils"
)

func main() {
	g := utils.ConstructBoardFromData("data/test.sav")
	g.PrintSquaresBoard()
}

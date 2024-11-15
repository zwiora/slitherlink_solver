package debug

import (
	"fmt"
	"time"

	"slytherlink_solver/utils"
)

var IsDebugMode bool

func Print(message any) {
	if IsDebugMode {
		fmt.Println(message)
	}
}

func PrintBoard(g *utils.Graph) {
	if IsDebugMode {
		g.PrintSquaresBoard(true)
	}
}

func Sleep(t int) {
	if IsDebugMode {
		time.Sleep(time.Duration(t) * time.Millisecond)
	}
}

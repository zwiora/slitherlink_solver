package debug

import (
	"fmt"
	"log"
	"time"
	// "slitherlink_solver/utils"
)

var IsDebugMode bool

func Println(message any) {
	if IsDebugMode {
		if message != "" {
			log.Println(message)
		} else {
			fmt.Println()
		}
	}
}

// func PrintBoard(g *utils.Graph) {
// 	if IsDebugMode {
// 		g.PrintSquaresBoard(true)
// 	}
// }

func Sleep(t int) {
	if IsDebugMode {
		time.Sleep(time.Duration(t) * time.Millisecond)
	}
}

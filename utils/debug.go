package utils

import (
	"fmt"
	"time"
)

var IsDebugMode bool

func DebugPrint(message string) {
	if IsDebugMode {
		fmt.Println(message)
	}
}

func DebugPrintBoard(g Graph) {
	if IsDebugMode {
		g.PrintSquaresBoard()
	}
}

func DebugSleep(t int) {
	if IsDebugMode {
		time.Sleep(time.Duration(t) * time.Millisecond)
	}
}

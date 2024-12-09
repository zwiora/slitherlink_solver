package utils

import (
	"fmt"
	"log"
	"time"
)

var IsHeuristicOn bool
var HeuristicType int
var NoVisitedStates int
var AvgDepth float64
var MaxDepth int

func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func TimeDuration(start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("%d", elapsed.Microseconds())
}

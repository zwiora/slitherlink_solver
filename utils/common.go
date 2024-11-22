package utils

import (
	"log"
)

var IsHeuristicOn bool
var NoVisitedStates int
var AvgDepth float32
var MaxDepth int

func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

package utils

import (
	"log"
)

var IsHeuristicOn bool

func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

package utils

import (
	"container/list"
	"fmt"

	"github.com/golang-collections/collections/stack"
)

type Graph struct {
	Root           *Node
	MaxDegree      int8
	maxCost        int
	SizeX          int
	SizeY          int
	shape          string
	AvailableMoves *list.List
	VisitedNodes   *stack.Stack
}

/* Calculates list of available moves at starting position */
func (g *Graph) CalculateStartingMoves() {
	g.AvailableMoves = list.New()
	thisNode := g.Root
	for i := 0; i < int(g.MaxDegree); i++ {
		for {
			if thisNode.Neighbours[i] != nil {
				thisNode = thisNode.Neighbours[i]
				g.AvailableMoves.PushFront(thisNode)
				thisNode.CanBeRemoved = true
			} else {
				break
			}
		}
	}
}

func (g *Graph) PrintSquaresBoard(isDebugMode bool) {
	lastLineNode := g.Root
	thisNode := g.Root

	fmt.Printf("-")

	for m := 0; m < g.SizeX; m++ {
		fmt.Printf("----")
	}
	fmt.Println()
	for n := 0; n < g.SizeY; n++ {
		fmt.Printf("|")
		for m := 0; m < g.SizeX; m++ {
			if thisNode.IsInLoop {
				fmt.Printf("\033[42m")
			}
			if isDebugMode && thisNode.IsVisited {
				fmt.Printf("x")
			} else {
				fmt.Printf(" ")
			}
			if thisNode.Value == -1 {
				fmt.Printf(" ")
			} else {
				fmt.Printf("%d", thisNode.Value)
			}
			if isDebugMode && thisNode.CanBeRemoved {
				fmt.Printf("#")
			} else {
				fmt.Printf(" ")
			}
			fmt.Printf("\033[49m|")
			thisNode = thisNode.Neighbours[0]
		}
		fmt.Println()
		fmt.Printf("-")

		for m := 0; m < g.SizeX; m++ {
			fmt.Printf("----")
		}
		fmt.Println()
		thisNode = lastLineNode.Neighbours[1]
		lastLineNode = thisNode
	}
}

/*
Change isVisited value from true to false in all nodes in the graph. Should be used after traversing the whole graph.
*/
func (g *Graph) ClearIsVisited() {
	thisNode := g.Root
	for {
		thisNode.IsVisited = false

		isNewNode := false
		for _, v := range thisNode.Neighbours {
			if v != nil && v.IsVisited {
				thisNode = v
				isNewNode = true
				break
			}
		}

		if !isNewNode {
			break
		}
	}
}

/*
Calculates sum of all visible values on the board and starting cost assuming there's a loop around whole board
*/
func (g *Graph) CalculateCost() (int, int) {
	fullCost := 0
	startCost := 0
	thisNode := g.Root

	for {
		thisNode.IsVisited = true

		if thisNode.Value >= 0 {
			fullCost += int(thisNode.Value)

			if thisNode.GetDegree() < int(g.MaxDegree) {
				startCost += thisNode.GetCostOfField(int(g.MaxDegree))
			} else {
				startCost += int(thisNode.Value)
			}
		}

		isNewNode := false
		for _, v := range thisNode.Neighbours {
			if v != nil && !v.IsVisited {
				thisNode = v
				isNewNode = true
				break
			}
		}

		if !isNewNode {
			break
		}
	}

	/* Clear isVisited parameters */

	g.ClearIsVisited()

	return fullCost, startCost
}

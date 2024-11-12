package utils

import (
	"fmt"
	"math"
)

type Node struct {
	Neighbours []*Node
	Value      int8
	IsInLoop   bool
	IsVisited  bool
}

type Graph struct {
	Root              *Node
	MaxNeighbourCount int8
	maxCost           int
	SizeX             int
	SizeY             int
	shape             string
}

func (n *Node) GetDegree() int {
	count := 0
	for _, v := range n.Neighbours {
		if v != nil && v.IsInLoop {
			count++
		}
	}

	return count
}

func (n *Node) GetLinesAround(maxNeighbourCount int) int {
	if n.IsInLoop {
		return maxNeighbourCount - n.GetDegree()
	}
	return n.GetDegree()
}

func (n *Node) GetCostOfField(maxNeighbourCount int) int {
	linesAround := n.GetLinesAround(int(maxNeighbourCount))
	return int(math.Abs(float64(linesAround) - float64(n.Value)))
}

func (g *Graph) PrintSquaresBoard() {
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
			if thisNode.IsVisited {
				fmt.Printf("x")
			} else {
				fmt.Printf(" ")
			}
			if thisNode.Value == -1 {
				fmt.Printf("  \033[49m|")
			} else {
				fmt.Printf("%d \033[49m|", thisNode.Value)
			}
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

func (g *Graph) calculatePerimeter() int {
	if g.shape == "square" {
		return g.SizeX*2 + g.SizeY*2 - 4
	}

	return 0
}

/*
Calculates sum of all visible values on the board and starting cost assuming there's a loop around whole board
*/
func (g *Graph) CalculateCost() (int, int) {
	fullCost := 0
	startCost := 0
	thisNode := g.Root
	countVisited := 0
	perimiter := g.calculatePerimeter()

	for {
		thisNode.IsVisited = true
		countVisited++

		if thisNode.Value >= 0 {
			fullCost += int(thisNode.Value)

			if countVisited <= perimiter {
				startCost += thisNode.GetCostOfField(int(g.MaxNeighbourCount))
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

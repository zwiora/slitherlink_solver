package utils

import "fmt"

type Node struct {
	neighbours []*Node
	value      int8
	isInLoop   bool
	isVisited  bool
}

type Graph struct {
	Root              *Node
	maxNeighbourCount int8
	maxCost           int
	sizeX             int
	sizeY             int
	shape             string
}

func (g *Graph) PrintSquaresBoard() {
	lastLineNode := g.Root
	thisNode := g.Root

	fmt.Printf("-")

	for m := 0; m < g.sizeX; m++ {
		fmt.Printf("----")
	}
	fmt.Println()
	for n := 0; n < g.sizeY; n++ {
		fmt.Printf("|")
		for m := 0; m < g.sizeX; m++ {
			if thisNode.isInLoop {
				fmt.Printf("\033[42m")
			}
			if thisNode.value == -1 {
				fmt.Printf("   \033[49m|")
			} else {
				fmt.Printf(" %d \033[49m|", thisNode.value)
			}
			thisNode = thisNode.neighbours[0]
		}
		fmt.Println()
		fmt.Printf("-")

		for m := 0; m < g.sizeX; m++ {
			fmt.Printf("----")
		}
		fmt.Println()
		thisNode = lastLineNode.neighbours[1]
		lastLineNode = thisNode
	}
}

/*
Change isVisited value from true to false in all nodes in the graph. Should be used after traversing the whole graph.
*/
func (g *Graph) clearIsVisited() {
	thisNode := g.Root
	for {
		thisNode.isVisited = false

		isNewNode := false
		for _, v := range thisNode.neighbours {
			if v != nil && v.isVisited {
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
		return g.sizeX*2 + g.sizeY*2 - 4
	}

	return 0
}

/*
Calculates sum of all visible values on the board
*/
func (g *Graph) CalculateCost() (int, int) {
	fullCost := 0
	startCost := 0
	thisNode := g.Root

	for {
		thisNode.isVisited = true

		if thisNode.value > 0 {
			fullCost += int(thisNode.value)
		}

		isNewNode := false
		for _, v := range thisNode.neighbours {
			if v != nil && !v.isVisited {
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

	g.clearIsVisited()

	return fullCost, startCost
}

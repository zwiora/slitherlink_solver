package utils

import "fmt"

type Node struct {
	neighbours []*Node
	value      int8
	deg        int8
	isInLoop   bool
}

type Graph struct {
	root              *Node
	maxNeighbourCount int8
	maxCost           int
	sizeX             int
	sizeY             int
}

func (g *Graph) PrintSquaresBoard() {
	lastLineNode := g.root
	thisNode := g.root

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

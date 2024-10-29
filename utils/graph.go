package utils

import "fmt"

type Node struct {
	neighbours []*Node
	value      int8
	deg        int8
}

type Graph struct {
	root              *Node
	maxNeighbourCount int8
	maxCost           int
	sizeX             int
	sizeY             int
}

func (g *Graph) PrintEmptyBoard() {
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
			if thisNode.value == -1 {
				fmt.Printf("   |")
			} else {
				fmt.Printf(" %d |", thisNode.value)
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

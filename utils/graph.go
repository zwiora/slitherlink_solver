package utils

import (
	"container/heap"
	"fmt"
	"slitherlink_solver/debug"

	"github.com/golang-collections/collections/stack"
)

type Graph struct {
	Root           *Node
	MaxDegree      int8
	maxCost        int
	SizeX          int
	SizeY          int
	FieldsCount    int
	Shape          string
	AvailableMoves *PriorityQueue
	VisitedNodes   *stack.Stack
}

/* Prints full board - type: squares */
func (g *Graph) printSquaresBoard(isDebugMode bool) {
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

			if isDebugMode && thisNode.IsDecided {

				if thisNode.IsForRemoval {
					fmt.Printf("\033[43m")
				} else {
					fmt.Printf("\033[41m")
				}
			}

			if thisNode.Value == -1 {
				fmt.Printf(" ")
			} else {
				fmt.Printf("%d", thisNode.Value)
			}

			if thisNode.IsInLoop {
				fmt.Printf("\033[42m")
			} else {
				fmt.Printf("\033[49m")
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

/* Prints full board - type: honeycomb */
func (g *Graph) printHoneycombBoard(isDebugMode bool) {
	lastLineNode := g.Root
	thisNode := g.Root

	fmt.Printf(" ")

	for m := 0; m < g.SizeX; m++ {
		if m%2 == 0 {
			fmt.Printf("___ ")
		} else {
			fmt.Printf("    ")
		}
	}
	fmt.Println()
	for n := 0; n < g.SizeY; n++ {
		for m := 0; m < g.SizeX; m++ {
			fmt.Print("/")
			if thisNode.IsInLoop {
				fmt.Printf("\033[42m")
			}
			if isDebugMode && thisNode.IsVisited {
				fmt.Printf("x")
			} else {
				fmt.Printf(" ")
			}

			if isDebugMode && thisNode.IsDecided {

				if thisNode.IsForRemoval {
					fmt.Printf("\033[43m")
				} else {
					fmt.Printf("\033[41m")
				}
			}

			if thisNode.Value == -1 {
				fmt.Printf(" ")
			} else {
				fmt.Printf("%d", thisNode.Value)
			}

			if thisNode.IsInLoop {
				fmt.Printf("\033[42m")
			} else {
				fmt.Printf("\033[49m")
			}

			if isDebugMode && thisNode.CanBeRemoved {
				fmt.Printf("#")
			} else {
				fmt.Printf(" ")
			}
			fmt.Printf("\033[49m")
			fmt.Print("\\")

			if (thisNode.Neighbours[0] != nil || (n == g.SizeY-1 && float32(m) >= float32(g.SizeX)/2)) && m != g.SizeX-1 {
				fmt.Print("___")
			}

			if m < (g.SizeX+1)/2 {
				if thisNode.Neighbours[0] != nil && thisNode.Neighbours[0].Neighbours[5] != nil {
					thisNode = thisNode.Neighbours[0].Neighbours[5]
				} else {
					thisNode = lastLineNode.Neighbours[0]
					if n != 0 && m == g.SizeX/2-1 {
						fmt.Print("/")
					}
					fmt.Println()
					fmt.Print("\\___")
				}
			} else {
				if thisNode.Neighbours[5] != nil {
					if thisNode.Neighbours[5].Neighbours[0] != nil {
						thisNode = thisNode.Neighbours[5].Neighbours[0]

					} else if n != g.SizeY-1 || m == (g.SizeX-1) {
						fmt.Print("___/")
					} else {
						fmt.Print("/")
					}

				}
			}
		}

		fmt.Println()

		thisNode = lastLineNode.Neighbours[1]
		lastLineNode = thisNode
	}
	fmt.Print("    ")
	for m := 0; m < g.SizeX/2; m++ {
		fmt.Print("\\___/   ")
	}

	fmt.Println()
}

/* Prints full board - type: triangle */
func (g *Graph) printTriangleBoard(isDebugMode bool) {
	lastLineNode := g.Root
	thisNode := g.Root
	width := g.SizeX / 2

	fmt.Print("  ")
	for m := 0; m < width; m++ {
		fmt.Print("------")
	}
	fmt.Println("-")

	for n := 0; n < g.SizeY/2; n++ {
		fmt.Print(" ")
		lastLineNode = thisNode
		for m := 0; m < width; m++ {
			fmt.Print("/ \\")
			thisNode = thisNode.Neighbours[0]

			if thisNode.IsInLoop {
				fmt.Printf("\033[42m")
			}
			if isDebugMode && thisNode.IsVisited {
				fmt.Printf("x")
			} else {
				fmt.Printf(" ")
			}

			if isDebugMode && thisNode.IsDecided {

				if thisNode.IsForRemoval {
					fmt.Printf("\033[43m")
				} else {
					fmt.Printf("\033[41m")
				}
			}

			if thisNode.Value == -1 {
				fmt.Printf(" ")
			} else {
				fmt.Printf("%d", thisNode.Value)
			}

			if thisNode.IsInLoop {
				fmt.Printf("\033[42m")
			} else {
				fmt.Printf("\033[49m")
			}

			if isDebugMode && thisNode.CanBeRemoved {
				fmt.Printf("#")
			} else {
				fmt.Printf(" ")
			}
			fmt.Printf("\033[49m")

			thisNode = thisNode.Neighbours[1]
		}

		fmt.Print("/ \\")
		fmt.Println()

		thisNode = lastLineNode

		for m := 0; m < width+1; m++ {
			fmt.Print("/")

			if thisNode.IsInLoop {
				fmt.Printf("\033[42m")
			}
			if isDebugMode && thisNode.IsVisited {
				fmt.Printf("x")
			} else {
				fmt.Printf(" ")
			}

			if isDebugMode && thisNode.IsDecided {

				if thisNode.IsForRemoval {
					fmt.Printf("\033[43m")
				} else {
					fmt.Printf("\033[41m")
				}
			}

			if thisNode.Value == -1 {
				fmt.Printf(" ")
			} else {
				fmt.Printf("%d", thisNode.Value)
			}

			if thisNode.IsInLoop {
				fmt.Printf("\033[42m")
			} else {
				fmt.Printf("\033[49m")
			}

			if isDebugMode && thisNode.CanBeRemoved {
				fmt.Printf("#")
			} else {
				fmt.Printf(" ")
			}
			fmt.Print("\033[49m\\ ")

			if thisNode.Neighbours[0] != nil {
				thisNode = thisNode.Neighbours[0]
				thisNode = thisNode.Neighbours[1]
			}
		}

		thisNode = lastLineNode.Neighbours[2]
		lastLineNode = thisNode
		fmt.Println()

		for m := 0; m < width; m++ {
			fmt.Print("------")
		}
		fmt.Println("-----")

		for m := 0; m < width+1; m++ {
			fmt.Print("\\")

			if thisNode.IsInLoop {
				fmt.Printf("\033[42m")
			}
			if isDebugMode && thisNode.IsVisited {
				fmt.Printf("x")
			} else {
				fmt.Printf(" ")
			}

			if isDebugMode && thisNode.IsDecided {

				if thisNode.IsForRemoval {
					fmt.Printf("\033[43m")
				} else {
					fmt.Printf("\033[41m")
				}
			}

			if thisNode.Value == -1 {
				fmt.Printf(" ")
			} else {
				fmt.Printf("%d", thisNode.Value)
			}

			if thisNode.IsInLoop {
				fmt.Printf("\033[42m")
			} else {
				fmt.Printf("\033[49m")
			}

			if isDebugMode && thisNode.CanBeRemoved {
				fmt.Printf("#")
			} else {
				fmt.Printf(" ")
			}
			fmt.Printf("\033[49m")
			fmt.Print("/ ")

			if thisNode.Neighbours[1] != nil {
				thisNode = thisNode.Neighbours[1]
				thisNode = thisNode.Neighbours[0]
			}
		}

		fmt.Println()

		thisNode = lastLineNode

		fmt.Print(" ")

		for m := 0; m < width; m++ {
			fmt.Print("\\ /")

			thisNode = thisNode.Neighbours[1]

			if thisNode.IsInLoop {
				fmt.Printf("\033[42m")
			}
			if isDebugMode && thisNode.IsVisited {
				fmt.Printf("x")
			} else {
				fmt.Printf(" ")
			}

			if isDebugMode && thisNode.IsDecided {

				if thisNode.IsForRemoval {
					fmt.Printf("\033[43m")
				} else {
					fmt.Printf("\033[41m")
				}
			}

			if thisNode.Value == -1 {
				fmt.Printf(" ")
			} else {
				fmt.Printf("%d", thisNode.Value)
			}

			if thisNode.IsInLoop {
				fmt.Printf("\033[42m")
			} else {
				fmt.Printf("\033[49m")
			}

			if isDebugMode && thisNode.CanBeRemoved {
				fmt.Printf("#")
			} else {
				fmt.Printf(" ")
			}
			fmt.Print("\033[49m")

			thisNode = thisNode.Neighbours[0]
		}

		fmt.Println("\\ /")
		fmt.Print("  ")
		for m := 0; m < width; m++ {
			fmt.Print("------")
		}
		fmt.Println("-")

		thisNode = lastLineNode.NextRow
	}

	if g.SizeY%2 == 1 {
		fmt.Print("   ")
		lastLineNode = thisNode
		for m := 0; m < width; m++ {
			fmt.Print("\\")
			if thisNode.IsInLoop {
				fmt.Printf("\033[42m")
			}
			if isDebugMode && thisNode.IsVisited {
				fmt.Printf("x")
			} else {
				fmt.Printf(" ")
			}

			if isDebugMode && thisNode.IsDecided {

				if thisNode.IsForRemoval {
					fmt.Printf("\033[43m")
				} else {
					fmt.Printf("\033[41m")
				}
			}

			if thisNode.Value == -1 {
				fmt.Printf(" ")
			} else {
				fmt.Printf("%d", thisNode.Value)
			}

			if thisNode.IsInLoop {
				fmt.Printf("\033[42m")
			} else {
				fmt.Printf("\033[49m")
			}

			if isDebugMode && thisNode.CanBeRemoved {
				fmt.Printf("#")
			} else {
				fmt.Printf(" ")
			}
			fmt.Printf("\033[49m/ ")

			if thisNode.Neighbours[1] != nil {
				thisNode = thisNode.Neighbours[1]
				thisNode = thisNode.Neighbours[0]
			}
		}

		fmt.Println()
		fmt.Print("    ")
		thisNode = lastLineNode

		for m := 0; m < width-1; m++ {
			thisNode = thisNode.Neighbours[1]

			fmt.Print("\\ /")

			if thisNode.IsInLoop {
				fmt.Printf("\033[42m")
			}
			if isDebugMode && thisNode.IsVisited {
				fmt.Printf("x")
			} else {
				fmt.Printf(" ")
			}

			if isDebugMode && thisNode.IsDecided {

				if thisNode.IsForRemoval {
					fmt.Printf("\033[43m")
				} else {
					fmt.Printf("\033[41m")
				}
			}

			if thisNode.Value == -1 {
				fmt.Printf(" ")
			} else {
				fmt.Printf("%d", thisNode.Value)
			}

			if thisNode.IsInLoop {
				fmt.Printf("\033[42m")
			} else {
				fmt.Printf("\033[49m")
			}

			if isDebugMode && thisNode.CanBeRemoved {
				fmt.Printf("#")
			} else {
				fmt.Printf(" ")
			}
			fmt.Print("\033[49m")

			thisNode = thisNode.Neighbours[0]

		}
		fmt.Println("\\ /")

		fmt.Print("     ")
		for m := 0; m < width-1; m++ {
			fmt.Print("------")
		}
		fmt.Println("-")

	}
}

func (g *Graph) PrintBoard(isDebugMode bool) {
	if g.Shape == "square" {
		g.printSquaresBoard(isDebugMode)
	} else if g.Shape == "honeycomb" {
		g.printHoneycombBoard(isDebugMode)
	} else if g.Shape == "triangle" {
		g.printTriangleBoard(isDebugMode)
	}
}

func (g *Graph) CheckIfSolutionOk() bool {

	if g.Shape == "square" {
		lastLineNode := g.Root
		thisNode := g.Root

		for n := 0; n < g.SizeY; n++ {
			for m := 0; m < g.SizeX; m++ {
				if thisNode.Value != -1 && thisNode.Value != int8(thisNode.getLinesAround(int(g.MaxDegree))) {
					return false
				}
				thisNode = thisNode.Neighbours[0]
			}
			thisNode = lastLineNode.Neighbours[1]
			lastLineNode = thisNode
		}
	} else if g.Shape == "honeycomb" {
		lastLineNode := g.Root
		thisNode := g.Root

		for n := 0; n < g.SizeY; n++ {
			direction := 5
			for m := 0; m < g.SizeX; m++ {
				if thisNode.Value != -1 && thisNode.Value != int8(thisNode.getLinesAround(int(g.MaxDegree))) {
					return false
				}
				if direction == 0 {
					direction = 5
				} else {
					direction = 0
				}
				thisNode = thisNode.Neighbours[direction]
			}
			thisNode = lastLineNode.Neighbours[1]
			lastLineNode = thisNode
		}
	} else if g.Shape == "triangle" {
		thisNode := g.Root
		i := 1
		counter := 0
		for {
			if thisNode.Value != -1 && thisNode.Value != int8(thisNode.getLinesAround(int(g.MaxDegree))) {
				return false
			}

			i = (i + 1) % 2

			if thisNode.Neighbours[i] == nil && counter%2 == 0 {
				thisNode = thisNode.Neighbours[2]
				counter++
				i = 1
			} else if thisNode.Neighbours[i] == nil && thisNode.NextRow != nil {
				thisNode = thisNode.NextRow
				counter++
				if g.SizeY%2 == 1 && counter == g.SizeY-1 {
					i = 0
				} else {
					i = 1
				}
			} else if thisNode.Neighbours[i] != nil {
				thisNode = thisNode.Neighbours[i]
			} else {
				break
			}
		}
	}
	return true

}

/*
Change isVisited value from true to false in all nodes in the graph. Should be used after traversing the whole graph.
*/
func (g *Graph) ClearIsVisited() {
	thisNode := g.Root

	i := 0
	if g.Shape == "triangle" {
		i = 1
	}
	for {
		thisNode.IsVisited = false

		if g.Shape == "honeycomb" {
			if i == 0 || i == 3 {
				i = (i - 1 + 6) % 6
			} else if i == 5 || i == 2 {
				i = (i + 1) % 6
			}
		} else if g.Shape == "triangle" {
			i = (i + 1) % 2
		}

		finished := false
		j := i
		for thisNode.Neighbours[i] == nil || !thisNode.Neighbours[i].IsVisited {
			i = (i + 1) % int(g.MaxDegree)
			if j == i {
				finished = true
				break
			}
		}

		if finished {
			if g.Shape == "triangle" && thisNode.NextRow != nil && thisNode.NextRow.IsVisited {
				thisNode = thisNode.NextRow
				i = 1
				continue
			}
			break
		}

		thisNode = thisNode.Neighbours[i]
	}
}

/*
Calculates sum of all visible values on the board and starting cost assuming there's a loop around the whole board
*/
func (g *Graph) CalculateStartCost() int {

	startCost := 0
	thisNode := g.Root

	i := 0
	if g.Shape == "triangle" {
		i = 1
	}
	for {
		thisNode.IsVisited = true
		if thisNode.Value >= 0 {
			if thisNode.GetDegree() < int(g.MaxDegree) {
				startCost += thisNode.getCostOfField(int(g.MaxDegree))
			} else {
				startCost += int(thisNode.Value)
			}
		}

		if g.Shape == "honeycomb" {
			if i == 0 || i == 3 {
				i = (i - 1 + 6) % 6
			} else if i == 5 || i == 2 {
				i = (i + 1) % 6
			}
		} else if g.Shape == "triangle" {
			i = (i + 1) % 2
		}

		finished := false
		j := i
		for thisNode.Neighbours[i] == nil || thisNode.Neighbours[i].IsVisited {
			i = (i + 1) % int(g.MaxDegree)
			if j == i {
				finished = true
				break
			}
		}

		if finished {
			if g.Shape == "triangle" && thisNode.NextRow != nil && !thisNode.NextRow.IsVisited {
				thisNode = thisNode.NextRow
				i = 1
				continue
			}

			break
		}

		thisNode = thisNode.Neighbours[i]
	}

	g.ClearIsVisited()

	return startCost
}

/* Calculate list of available moves at starting position */
func (g *Graph) CalculateStartMoves() {
	movesArr := []*Node{}
	thisNode := g.Root

	if g.Shape == "square" {
		for i := 0; i < int(g.MaxDegree); i++ {
			for {
				if thisNode.Neighbours[i] == nil {
					break
				}
				thisNode = thisNode.Neighbours[i]

				thisNode.SetNodeCost(g)
				if !(thisNode.IsDecided && !thisNode.IsForRemoval) {
					thisNode.CanBeRemoved = true
					movesArr = append(movesArr, thisNode)
				}

			}
		}
	} else if g.Shape == "honeycomb" {
		i := 5
		for {
			if !(thisNode.IsDecided && !thisNode.IsForRemoval) {
				thisNode.SetNodeCost(g)
				thisNode.CanBeRemoved = true
				movesArr = append(movesArr, thisNode)
			}

			if i == 0 || i == 3 {
				i = (i - 1 + 6) % 6
			} else if i == 5 || i == 2 {
				i = (i + 1) % 6
			}

			for thisNode.Neighbours[i] == nil {
				i = (i + 1) % 6
			}

			thisNode = thisNode.Neighbours[i]

			if thisNode == g.Root {
				break
			}
		}
	} else if g.Shape == "triangle" {
		thisNode := g.Root
		i := 1
		for {

			thisNode.IsVisited = true

			if thisNode.GetDegree() < int(g.MaxDegree) {
				thisNode.SetNodeCost(g)
				thisNode.CanBeRemoved = true
				movesArr = append(movesArr, thisNode)

			}

			i = (i + 1) % 2

			finished := false
			j := i
			for thisNode.Neighbours[i] == nil || thisNode.Neighbours[i].IsVisited {
				i = (i + 1) % int(g.MaxDegree)
				if j == i {
					finished = true
					break
				}
			}

			if finished {
				if thisNode.NextRow != nil && !thisNode.NextRow.IsVisited {
					thisNode = thisNode.NextRow
					i = 1
					continue
				}

				break
			}

			thisNode = thisNode.Neighbours[i]
		}

		g.ClearIsVisited()
	}

	/* Transform array into heap */
	pq := make(PriorityQueue, len(movesArr))
	g.AvailableMoves = &pq

	i := 0
	for _, v := range movesArr {
		(*g.AvailableMoves)[i] = v
		v.QueueIndex = i
		i++
	}

	heap.Init(g.AvailableMoves)
}

/* Should be run after preparation of the solver but before its start */
func (g *Graph) FindTemplates() {
	debug.Println("check templates")

	thisNode := g.Root

	i := 0
	if g.Shape == "triangle" {
		i = 1
	}
	for {
		thisNode.IsVisited = true

		thisNode.findZeroTemplates(g)

		if g.Shape == "honeycomb" {
			if i == 0 || i == 3 {
				i = (i - 1 + 6) % 6
			} else if i == 5 || i == 2 {
				i = (i + 1) % 6
			}
		} else if g.Shape == "triangle" {
			i = (i + 1) % 2
		}

		finished := false
		j := i
		for thisNode.Neighbours[i] == nil || thisNode.Neighbours[i].IsVisited {
			i = (i + 1) % int(g.MaxDegree)
			if j == i {
				finished = true
				break
			}
		}

		if finished {
			if g.Shape == "triangle" && thisNode.NextRow != nil && !thisNode.NextRow.IsVisited {
				thisNode = thisNode.NextRow
				i = 1
				continue
			}
			break
		}

		thisNode = thisNode.Neighbours[i]
	}

	g.ClearIsVisited()

	for {
		newTemplatesFound := 0
		thisNode := g.Root

		i := 0
		if g.Shape == "triangle" {
			i = 1
		}
		for {
			thisNode.IsVisited = true

			// if thisNode.findNumberTemplates(g) {
			// 	newTemplatesFound++
			// }

			if thisNode.find33Templates(g) {
				newTemplatesFound++
			}

			// if thisNode.find33CornerTemplates(g) {
			// 	newTemplatesFound++
			// }

			// if thisNode.find31Templates(g) {
			// 	newTemplatesFound++
			// }

			// if thisNode.findloopReachingNumberTemplates(g) {
			// 	newTemplatesFound++
			// }

			if thisNode.findContinousSquareTemplates(g) {
				newTemplatesFound++
			}

			if g.Shape == "honeycomb" {
				if i == 0 || i == 3 {
					i = (i - 1 + 6) % 6
				} else if i == 5 || i == 2 {
					i = (i + 1) % 6
				}
			} else if g.Shape == "triangle" {
				i = (i + 1) % 2
			}

			finished := false
			j := i
			for thisNode.Neighbours[i] == nil || thisNode.Neighbours[i].IsVisited {
				i = (i + 1) % int(g.MaxDegree)
				if j == i {
					finished = true
					break
				}
			}

			if finished {
				if g.Shape == "triangle" && thisNode.NextRow != nil && !thisNode.NextRow.IsVisited {
					thisNode = thisNode.NextRow
					i = 1
					continue
				}
				break
			}

			thisNode = thisNode.Neighbours[i]
		}

		g.ClearIsVisited()
		g.PrintBoard(true)

		if newTemplatesFound == 0 {
			break
		}
	}

}

package utils

import (
	"container/heap"
	"fmt"
	"slitherlink_solver/debug"

	"github.com/golang-collections/collections/queue"
	"github.com/golang-collections/collections/stack"
)

type Graph struct {
	Root           *Node
	MaxDegree      int8
	maxCost        int
	SizeX          int
	SizeY          int
	shape          string
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
			if n == g.SizeY-1 && m == g.SizeX-1 {
				fmt.Print("___")
			}
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

			if thisNode.Neighbours[0] != nil {
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

					} else if n != g.SizeY-1 {
						fmt.Print("/")
					} else {
						fmt.Print("___/")
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

func (g *Graph) PrintBoard(isDebugMode bool) {
	if g.shape == "square" {
		g.printSquaresBoard(isDebugMode)
	} else if g.shape == "honeycomb" {
		g.printHoneycombBoard(isDebugMode)
	}
}

func (g *Graph) CheckIfSolutionOk() bool {

	if g.shape == "square" {
		lastLineNode := g.Root
		thisNode := g.Root

		for n := 0; n < g.SizeY; n++ {
			for m := 0; m < g.SizeX; m++ {
				if thisNode.Value != -1 && thisNode.Value != int8(thisNode.getLinesAround(int(g.MaxDegree))) {
					// fmt.Println(thisNode)
					return false
				}
				thisNode = thisNode.Neighbours[0]
			}
			thisNode = lastLineNode.Neighbours[1]
			lastLineNode = thisNode
		}

		return true
	}
	return false

}

/*
Change isVisited value from true to false in all nodes in the graph. Should be used after traversing the whole graph.
*/
func (g *Graph) ClearIsVisited() {
	thisNode := g.Root

	i := 0
	for {
		thisNode.IsVisited = false

		g.PrintBoard(true)

		if g.shape == "honeycomb" {
			if i == 0 || i == 3 {
				i = (i - 1 + 6) % 6
			} else if i == 5 || i == 2 {
				i = (i + 1) % 6
			}
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
	for {
		thisNode.IsVisited = true

		g.PrintBoard(true)

		if thisNode.Value >= 0 {
			if thisNode.GetDegree() < int(g.MaxDegree) {
				startCost += thisNode.getCostOfField(int(g.MaxDegree))
			} else {
				startCost += int(thisNode.Value)
			}
		}

		if g.shape == "honeycomb" {
			if i == 0 || i == 3 {
				i = (i - 1 + 6) % 6
			} else if i == 5 || i == 2 {
				i = (i + 1) % 6
			}
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

	if g.shape == "square" {
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
	} else if g.shape == "honeycomb" {
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

/* Constructs the queue of nodes, that will be searched for templates */
func (g *Graph) constructQueueForCheckingTemplates() *queue.Queue {
	nodes := queue.New()
	thisNode := g.Root

	for {
		thisNode.IsVisited = true

		// all templates not near edge use 0 or 3
		if thisNode.Value == 0 || thisNode.Value == 3 {
			nodes.Enqueue(thisNode)
		} else if thisNode.GetDegree() < int(g.MaxDegree) {
			nodes.Enqueue(thisNode)
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

	g.ClearIsVisited()
	return nodes
}

/* Should be run after preparation of the solver but before its start */
func (g *Graph) FindTemplates() {
	nodes := g.constructQueueForCheckingTemplates()
	debug.Println("check templates")

	for nodes.Len() > 0 {
		thisNode := (nodes.Dequeue()).(*Node)

		/* If the final state of the node isn't set */
		if !thisNode.findZeroTemplates(g) {
			thisNode.findNumberTemplates(g)

		}

	}

	for {

		newTemplatesFound := 0
		thisNode := g.Root

		for {
			thisNode.IsVisited = true

			// checking only templates that use state of other nodes
			if thisNode.Value != -1 {
				if thisNode.findNumberTemplates(g) {
					newTemplatesFound++
				}
			}

			if thisNode.find31Templates(g) {
				newTemplatesFound++
				// fmt.Println("SUPER")
				// g.PrintSquaresBoard(true)
				// time.Sleep(time.Millisecond * 1000)
			}

			if thisNode.find33Templates(g) {
				newTemplatesFound++
			}

			if thisNode.find3and3Templates(g) {
				newTemplatesFound++
			}

			if thisNode.findloopReachingNumberTemplates(g) {
				newTemplatesFound++
			}

			if thisNode.findContinousSquareTemplates(g) {
				newTemplatesFound++
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

		g.ClearIsVisited()

		if newTemplatesFound == 0 {
			break
		}
	}

}

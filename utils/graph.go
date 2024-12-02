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

func (g *Graph) CheckIfSolutionOk() bool {

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
Calculates sum of all visible values on the board and starting cost assuming there's a loop around the whole board
*/
func (g *Graph) CalculateStartCost() (int, int) {
	fullCost := 0
	startCost := 0
	thisNode := g.Root

	for {
		thisNode.IsVisited = true

		if thisNode.Value >= 0 {
			fullCost += int(thisNode.Value)

			if thisNode.GetDegree() < int(g.MaxDegree) {
				startCost += thisNode.getCostOfField(int(g.MaxDegree))
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

	g.ClearIsVisited()

	return fullCost, startCost
}

/* Calculate list of available moves at starting position */
func (g *Graph) CalculateStartMoves() {
	movesArr := []*Node{}
	thisNode := g.Root

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
		// if !thisNode.findCornerTemplates(g, nodes) {
		// 	thisNode.find31Templates(g, nodes)
		// }

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
				// if !thisNode.findCornerTemplates(g, nodes) {
				// 	if thisNode.find31Templates(g, nodes) {
				// 		newTemplatesFound++
				// 	}
				// } else {
				// 	newTemplatesFound++
				// }
			}

			if thisNode.find33Templates(g) {
				newTemplatesFound++
			}

			if thisNode.find3and3Templates(g) {
				// fmt.Println("SUPER")
				// g.PrintSquaresBoard(true)
				// time.Sleep(time.Millisecond * 1000)

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

		// fmt.Println("REPEAT")
		// g.PrintSquaresBoard(true)
		// time.Sleep(1000 * time.Millisecond)
		g.ClearIsVisited()

		if newTemplatesFound == 0 {
			break
		}
	}

}

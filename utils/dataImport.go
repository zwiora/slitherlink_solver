package utils

import (
	"os"
	"strconv"
	"strings"
)

/*
Reading file content as a string
*/
func readFile(fileName string) string {
	dat, err := os.ReadFile(fileName)
	Check(err)
	return string(dat)
}

/*
Extracting data about the puzzle from the file:
- type of slitherlink board (by its code)
- size of slitherlink board (width)
- size of slitherlink board (height)
- content of slitherlink board (encoded as string)
*/
func readPuzzleStructure(fileName string) (string, int, int, string) {
	fileContent := readFile(fileName)
	puzzleCode := strings.Split(strings.Split(fileContent, "\n")[3], ":")[2]
	puzzleCodeArr := strings.Split(puzzleCode, "t")
	puzzleType := puzzleCodeArr[1]
	puzzleSizeArr := strings.Split(puzzleCodeArr[0], "x")
	puzzleSizeX, err := strconv.Atoi(puzzleSizeArr[0])
	Check(err)
	puzzleSizeY, err := strconv.Atoi(puzzleSizeArr[1])
	Check(err)
	puzzleContent := strings.Split(strings.Split(fileContent, "\n")[6], ":")[2]
	return puzzleType, puzzleSizeX, puzzleSizeY, puzzleContent
}

func ConstructBoardFromData(fileName string) *Graph {
	puzzleType, puzzleSizeX, puzzleSizeY, puzzleContent := readPuzzleStructure(fileName)

	var board Graph
	var thisNode *Node
	thisNode = &Node{
		Value:    -1,
		IsInLoop: true,
	}
	lastLineNode := thisNode
	board.Root = thisNode
	board.maxCost = 0
	board.SizeX = puzzleSizeX
	board.SizeY = puzzleSizeY
	board.shape = "square"

	if puzzleType == "0de" {
		m := 0
		n := 0

		board.MaxNeighbourCount = 4
		thisNode.Neighbours = make([]*Node, board.MaxNeighbourCount)

		/* Preparing first row */
		for m := 0; m < puzzleSizeX-1; m++ {

			nextNode := &Node{
				Value:    -1,
				IsInLoop: true,
			}
			nextNode.Neighbours = make([]*Node, board.MaxNeighbourCount)

			thisNode.Neighbours[0] = nextNode
			nextNode.Neighbours[2] = thisNode

			thisNode = nextNode

		}

		thisNode = board.Root

		/* Setting content and preparing rest of the nodes */
		for _, character := range strings.Split(puzzleContent, "") {
			characterVal, err := strconv.Atoi(character)

			/* Setting value of the node */
			nodesCounter := 1
			if err == nil {
				thisNode.Value = int8(characterVal)
				board.maxCost += characterVal
			} else {
				nodesCounter = int(character[0]) - int('a') + 1
			}

			for i := 0; i < nodesCounter; i++ {

				if n < puzzleSizeY-1 {

					/* Connect this node with bottom one and vice versa*/
					bottomNode := &Node{
						Value:    -1,
						IsInLoop: true,
					}
					bottomNode.Neighbours = make([]*Node, board.MaxNeighbourCount)
					thisNode.Neighbours[1] = bottomNode
					bottomNode.Neighbours[3] = thisNode

					/* Connect bottom node with its left neighbour and vice versa */
					if m > 0 {
						bottomNode.Neighbours[2] = thisNode.Neighbours[2].Neighbours[1]
						thisNode.Neighbours[2].Neighbours[1].Neighbours[0] = bottomNode
					}
				}

				// for _, v := range thisNode.neighbours {
				// 	if v != nil {
				// 		thisNode.deg++
				// 	}
				// }

				/* Calculating next position */
				m++
				if m >= puzzleSizeX {
					m = 0
					n++
					thisNode = lastLineNode.Neighbours[1]
					lastLineNode = lastLineNode.Neighbours[1]
				} else {
					thisNode = thisNode.Neighbours[0]
				}

				_ = puzzleSizeY
			}
		}
	}

	return &board
}

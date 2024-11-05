package utils

import (
	"fmt"
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
- size of slitherlink board (width and height)
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
		value:    -1,
		isInLoop: true,
	}
	lastLineNode := thisNode
	board.root = thisNode
	board.maxCost = 0
	board.sizeX = puzzleSizeX
	board.sizeY = puzzleSizeY

	if puzzleType == "0de" {
		m := 0
		n := 0

		board.maxNeighbourCount = 4
		thisNode.neighbours = make([]*Node, board.maxNeighbourCount)

		/* Preparing first row */
		for m := 0; m < puzzleSizeX-1; m++ {

			nextNode := &Node{
				value:    -1,
				isInLoop: true,
			}
			nextNode.neighbours = make([]*Node, board.maxNeighbourCount)

			thisNode.neighbours[0] = nextNode
			nextNode.neighbours[2] = thisNode

			thisNode = nextNode

		}

		thisNode = board.root

		/* Setting content and preparing rest of the nodes */
		for _, character := range strings.Split(puzzleContent, "") {
			characterVal, err := strconv.Atoi(character)

			/* Setting value of the node */
			nodesCounter := 1
			if err == nil {
				thisNode.value = int8(characterVal)
				board.maxCost += characterVal
			} else {
				nodesCounter = int(character[0]) - int('a') + 1
			}

			for i := 0; i < nodesCounter; i++ {

				if n < puzzleSizeY-1 {

					/* Connect this node with bottom one and vice versa*/
					bottomNode := &Node{
						value:    -1,
						isInLoop: true,
					}
					bottomNode.neighbours = make([]*Node, board.maxNeighbourCount)
					thisNode.neighbours[1] = bottomNode
					bottomNode.neighbours[3] = thisNode

					/* Connect bottom node with its left neighbour and vice versa */
					if m > 0 {
						bottomNode.neighbours[2] = thisNode.neighbours[2].neighbours[1]
						thisNode.neighbours[2].neighbours[1].neighbours[0] = bottomNode
					}
				}

				for _, v := range thisNode.neighbours {
					if v != nil {
						thisNode.deg++
					}
				}
				fmt.Println(thisNode)

				/* Calculating next position */
				m++
				if m >= puzzleSizeX {
					m = 0
					n++
					fmt.Println()
					thisNode = lastLineNode.neighbours[1]
					lastLineNode = lastLineNode.neighbours[1]
				} else {
					thisNode = thisNode.neighbours[0]
				}

				_ = puzzleSizeY
			}
		}
	}

	return &board
}

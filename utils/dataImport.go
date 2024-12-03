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

/*
Extracting multiple data about the puzzle from the file:
- type of slitherlink board (by its code)
- size of slitherlink board (width)
- size of slitherlink board (height)
- content of slitherlink board (encoded as string)
*/
func ReadMultiplePuzzleStructure(fileName string) [][]string {
	result := [][]string{}
	fileContent := strings.Split(readFile(fileName), "\n")
	boardtype := fileContent[0]
	sizeX := ""
	sizeY := ""
	for i := 1; i < len(fileContent); i++ {
		if strings.Split(fileContent[i], " ")[0] == "s" {
			sizeX = strings.Split(strings.Split(fileContent[i], " ")[1], "x")[0]
			sizeY = strings.Split(strings.Split(fileContent[i], " ")[1], "x")[1]
			continue
		}
		result = append(result, []string{boardtype, sizeX, sizeY, fileContent[i]})
	}
	return result
}

func constructSquaresBoard(board *Graph, puzzleContent string) {
	thisNode := board.Root
	lastLineNode := thisNode
	m := 0
	n := 0

	board.MaxDegree = 4
	thisNode.Neighbours = make([]*Node, board.MaxDegree)

	/* Preparing first row */
	for m := 0; m < board.SizeX-1; m++ {

		nextNode := &Node{
			Value:      -1,
			IsInLoop:   true,
			QueueIndex: -1,
		}
		nextNode.Neighbours = make([]*Node, board.MaxDegree)

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

			if n < board.SizeY-1 {

				/* Connect this node with bottom one and vice versa*/
				bottomNode := &Node{
					Value:    -1,
					IsInLoop: true,
				}
				bottomNode.Neighbours = make([]*Node, board.MaxDegree)
				thisNode.Neighbours[1] = bottomNode
				bottomNode.Neighbours[3] = thisNode

				/* Connect bottom node with its left neighbour and vice versa */
				if m > 0 {
					bottomNode.Neighbours[2] = thisNode.Neighbours[2].Neighbours[1]
					thisNode.Neighbours[2].Neighbours[1].Neighbours[0] = bottomNode
				}
			}

			/* Calculating next position */
			m++
			if m >= board.SizeX {
				m = 0
				n++
				thisNode = lastLineNode.Neighbours[1]
				lastLineNode = lastLineNode.Neighbours[1]
			} else {
				thisNode = thisNode.Neighbours[0]
			}
		}
	}
}

func constructHexBoard(board *Graph, puzzleContent string) {
	thisNode := board.Root
	lastLineNode := thisNode
	m := 0
	n := 0

	board.MaxDegree = 6
	thisNode.Neighbours = make([]*Node, board.MaxDegree)

	/* Preparing first row */
	for m := 0; m < board.SizeX-1; m++ {

		nextNode := &Node{
			Value:      -1,
			IsInLoop:   true,
			QueueIndex: -1,
		}
		nextNode.Neighbours = make([]*Node, board.MaxDegree)

		i := (6 - m%2) % 6
		j := (i + 3) % 6
		thisNode.Neighbours[i] = nextNode
		nextNode.Neighbours[j] = thisNode

		thisNode = nextNode
	}

	thisNode = board.Root

	direction := 0

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

			if n < board.SizeY-1 {

				/* Connect this node with bottom one and vice versa*/
				bottomNode := &Node{
					Value:    -1,
					IsInLoop: true,
				}
				bottomNode.Neighbours = make([]*Node, board.MaxDegree)
				thisNode.Neighbours[1] = bottomNode
				bottomNode.Neighbours[4] = thisNode

				/* Connect bottom node with its top right neighbour (we do this with only half of the nodes) */
				if thisNode.Neighbours[0] != nil {
					bottomNode.Neighbours[5] = thisNode.Neighbours[0]
					thisNode.Neighbours[0].Neighbours[2] = bottomNode
				}

				/* Connect bottom node with its bottom left neighbours (we do this with only half of the nodes) */
				if thisNode.Neighbours[2] != nil {
					newNode := thisNode.Neighbours[2].Neighbours[1]
					if newNode != nil {
						newNode.Neighbours[5] = bottomNode
						bottomNode.Neighbours[2] = newNode
					}
				}

				/* Connect bottom node with its left top neighbour and vice versa */
				if thisNode.Neighbours[2] != nil {
					bottomNode.Neighbours[3] = thisNode.Neighbours[2]
					thisNode.Neighbours[2].Neighbours[0] = bottomNode
				}
			}
			// fmt.Println(thisNode)

			/* Calculating next position */
			m++
			if m >= board.SizeX {
				m = 0
				n++
				thisNode = lastLineNode.Neighbours[1]
				lastLineNode = lastLineNode.Neighbours[1]
				direction = 0
			} else {
				thisNode = thisNode.Neighbours[direction]
				if direction == 0 {
					direction = 5
				} else {
					direction = 0
				}
			}
		}
	}
}

func ConstructBoardFromData(puzzleType string, puzzleSizeX int, puzzleSizeY int, puzzleContent string) *Graph {

	board := new(Graph)
	var thisNode *Node
	thisNode = &Node{
		Value:      -1,
		IsInLoop:   true,
		QueueIndex: -1,
	}
	board.Root = thisNode
	board.maxCost = 0
	board.SizeX = puzzleSizeX
	board.SizeY = puzzleSizeY

	/* Type: squares" */
	if puzzleType == "0de" {
		board.shape = "square"
		constructSquaresBoard(board, puzzleContent)
	} else if puzzleType == "2" {
		board.shape = "honeycomb"
		constructHexBoard(board, puzzleContent)
	}

	return board
}

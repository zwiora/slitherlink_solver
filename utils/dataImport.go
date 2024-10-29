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
- size of slitherlink board
- content of slitherlink board (encoded as string)
*/
func ReadPuzzleStructure(fileName string) (string, int, int, string) {
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

// func ConstructBoardFromData(fileName string) {
// 	// puzzleType, puzzleSize, puzzleContent := ReadPuzzleStructure(fileName)
// 	_, _, puzzleContent := ReadPuzzleStructure(fileName)

// }

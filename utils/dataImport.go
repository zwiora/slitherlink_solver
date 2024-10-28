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
func ReadPuzzleStructure(fileName string) (string, int, string) {
	fileContent := readFile(fileName)
	puzzleCode := strings.Split(strings.Split(fileContent, "\n")[3], ":")[2]
	puzzleType := strings.Split(puzzleCode, "t")[1]
	puzzleSize, err := strconv.Atoi(strings.Split(puzzleCode, "x")[0])
	Check(err)
	puzzleContent := strings.Split(strings.Split(fileContent, "\n")[6], ":")[2]
	return puzzleType, puzzleSize, puzzleContent
}

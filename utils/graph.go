package utils

type Node struct {
	neighbours []*Node
	value      int8
}

type Graph struct {
	root              *Node
	maxNeighbourCount int8
}

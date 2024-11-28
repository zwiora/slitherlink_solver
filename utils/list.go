package utils

type List struct {
	root         *ListElem
	oppositeList *List
}

type ListElem struct {
	value *Node
	next  *ListElem
}

func (l List) addElement(node *Node) {
	newElement := new(ListElem)
	newElement.value = node
	if l.root == nil {
		newElement.next = newElement
	} else {
		newElement.next = l.root
		l.root = newElement
	}
}

// func (l List) markList(isClearing )

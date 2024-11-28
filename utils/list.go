package utils

import (
	"fmt"

	"github.com/golang-collections/collections/queue"
)

type List struct {
	root         *ListElem
	oppositeList *List
	length       int
}

type ListElem struct {
	value *Node
	next  *ListElem
}

func (l *List) addElement(node *Node) {
	l.length++
	newElement := new(ListElem)
	newElement.value = node
	node.TemplateGroup = l
	if l.root == nil {
		newElement.next = newElement
		l.root = newElement
	} else {
		newElement.next = l.root.next
		l.root.next = newElement

	}
}

func (l *List) setValue(isForRemoval bool, q *queue.Queue) {
	if l != nil && !l.root.value.IsDecided {
		thisElement := l.root
		for {
			thisElement.value.IsDecided = true
			thisElement.value.IsForRemoval = isForRemoval
			addNeighboursToQueue(thisElement.value, q)
			thisElement = thisElement.next

			if thisElement == l.root {
				break
			}
		}

		l.oppositeList.setValue(!isForRemoval, q)
	}
}

func (l *List) print() {
	fmt.Println("List: ", l)
	thisElement := l.root
	for {
		fmt.Println(thisElement.value)
		thisElement = thisElement.next

		if thisElement == l.root {
			break
		}
	}
}

func addLists(l1 *List, l2 *List) {

	if l1 == l2 {
		return
	}

	basicList := l1
	additionalList := l2
	if l1.length < l2.length {
		basicList = l2
		additionalList = l1
	}

	thisElement := additionalList.root
	var lastElement *ListElem
	for {
		thisElement.next.value.TemplateGroup = basicList
		thisElement = thisElement.next

		if thisElement.next == additionalList.root {
			lastElement = thisElement
			break
		}
	}

	lastElement.next = basicList.root.next
	basicList.root.next = additionalList.root

	basicList.length += additionalList.length
}

// func (l List) markList(isClearing )

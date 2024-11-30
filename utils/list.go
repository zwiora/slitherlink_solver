package utils

import (
	"fmt"
)

type List struct {
	Root         *ListElem
	OppositeList *List
	Removable    int
	Removed      int
	Length       int
	SettingNode  *Node
}

type ListElem struct {
	Value *Node
	Next  *ListElem
}

func (l *List) addElement(node *Node) {
	l.Length++
	newElement := new(ListElem)
	newElement.Value = node
	node.TemplateGroup = l
	if l.Root == nil {
		newElement.Next = newElement
		l.Root = newElement
	} else {
		newElement.Next = l.Root.Next
		l.Root.Next = newElement

	}
}

func (l *List) SetValue(isForRemoval bool, settingNode *Node, g *Graph) bool {
	if l != nil && !l.Root.Value.IsDecided {

		l.SettingNode = settingNode

		if isForRemoval {
			l.Removable = l.Length
		}

		thisElement := l.Root
		for {
			/* Checking if neighbour would have enough edges*/
			// if settingNode != nil {
			// 	if isForRemoval && thisElement.Value.IsDeletionBreakingSecondRule() {
			// 		/* Deleting this element is against the rules */
			// 		return false
			// 	}
			// }

			thisElement.Value.IsDecided = true
			thisElement.Value.IsForRemoval = isForRemoval
			// if isForRemoval && thisElement.Value.CanBeRemoved {
			// 	thisElement.Value.UpdateNodeCost(g)
			// }
			thisElement = thisElement.Next

			if thisElement == l.Root {
				break
			}
		}

		l.OppositeList.SetValue(!isForRemoval, settingNode, g)
	}

	return true
}

func (l *List) ClearValue(g *Graph) {
	if l != nil && l.Root.Value.IsDecided {
		// fmt.Println("CLEARING")
		// time.Sleep(1000 * time.Millisecond)
		// g.PrintSquaresBoard(true)

		l.Removed = 0
		l.Removable = 0

		thisElement := l.Root
		for {

			if thisElement.Value.IsDecided && thisElement.Value.CanBeRemoved {

				thisElement.Value.IsDecided = false
				thisElement.Value.IsForRemoval = false

				thisElement.Value.UpdateNodeCost(g)
			}

			thisElement = thisElement.Next

			if thisElement == l.Root {
				break
			}
		}

		l.OppositeList.ClearValue(g)
	}
}

func (l *List) print() {
	fmt.Println("List: ", l)
	thisElement := l.Root
	for {
		fmt.Println(thisElement.Value)
		thisElement = thisElement.Next

		if thisElement == l.Root {
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
	if l1.Length < l2.Length {
		basicList = l2
		additionalList = l1
	}

	thisElement := additionalList.Root
	var lastElement *ListElem
	for {
		thisElement.Next.Value.TemplateGroup = basicList
		thisElement = thisElement.Next

		if thisElement.Next == additionalList.Root {
			lastElement = thisElement
			break
		}
	}

	lastElement.Next = basicList.Root.Next
	basicList.Root.Next = additionalList.Root

	basicList.Length += additionalList.Length
}

// func (l List) markList(isClearing )

package utils

import (
	"fmt"
)

type List struct {
	Root         *ListElem
	OppositeList *List
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

func (l *List) addOppositeElement(node *Node) {
	if l.OppositeList == nil {
		l.OppositeList = new(List)
		l.OppositeList.OppositeList = l
	}
	l = l.OppositeList
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
	result := true
	if l != nil && !l.Root.Value.IsDecided {

		l.SettingNode = settingNode

		thisElement := l.Root
		for {
			/* Checking if neighbour would have enough edges*/
			if settingNode != nil {
				if isForRemoval && thisElement.Value.IsDeletionBreakingSecondRule() {
					/* Deleting this element is against the rules */
					return false
				}
			} else {
				thisElement.Value.TemplateGroup = nil
			}

			thisElement.Value.IsDecided = true
			thisElement.Value.IsForRemoval = isForRemoval
			if isForRemoval && thisElement.Value.CanBeRemoved {
				thisElement.Value.UpdateNodeCost(g)
			}
			thisElement = thisElement.Next

			if thisElement == l.Root {
				break
			}
		}
		result = l.OppositeList.SetValue(!isForRemoval, settingNode, g)
	}
	return result
}

func (l *List) ClearValue(g *Graph) {
	if l != nil && l.Root.Value.IsDecided {

		thisElement := l.Root
		for {
			if thisElement.Value.IsDecided {
				thisElement.Value.IsDecided = false
				thisElement.Value.IsForRemoval = false

				if thisElement.Value.CanBeRemoved {
					thisElement.Value.UpdateNodeCost(g)
				}
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

func concatLists(l1 *List, l2 *List) *List {
	if l1 == nil {
		return l2
	} else if l2 == nil {
		return l1
	}
	basicList := l1
	additionalList := l2
	if l1.Length < l2.Length {
		basicList = l2
		additionalList = l1
	}

	thisElement := additionalList.Root
	thisElement.Value.TemplateGroup = basicList
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

	return basicList
}

func addLists(l1 *List, l2 *List) {
	if l1 == l2 {
		return
	}
	opposite := concatLists(l1.OppositeList, l2.OppositeList)
	list := concatLists(l1, l2)
	if opposite != nil {
		list.OppositeList = opposite
		opposite.OppositeList = list
	}
}

func addOppositeLists(l1 *List, l2 *List) {
	if l1 == l2.OppositeList || l1 == l2 {
		return
	}

	op1 := l1.OppositeList
	op2 := l2.OppositeList

	opposite := concatLists(l1, op2)
	list := concatLists(op1, l2)
	list.OppositeList = opposite
	opposite.OppositeList = list
}

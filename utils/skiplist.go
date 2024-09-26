package utils

import (
	"fmt"
	"math"
)

type NodeData struct {
	Key string
	TTL int
}

type Node struct {
	Data  NodeData
	Up    *Node
	Down  *Node
	Left  *Node
	Right *Node
}

type SkipList struct {
	Head          *Node
	Tail          *Node
	NumOfElements int
	Height        int
}

// -1 <-> max
func CreateSkipList() *SkipList {
	HeadNode := &Node{Data: NodeData{"-INF", -1}}
	TailNode := &Node{Data: NodeData{"INF", math.MaxInt32}}

	HeadNode.Right = TailNode
	TailNode.Left = HeadNode

	return &SkipList{Head: HeadNode, Tail: TailNode}
}

// Returns the entry. If nothing is matched, returns the immediate smaller element in the lowest level
func (skipList *SkipList) FindEntry(key string, ttl int) *Node {
	node := NodeData{key, ttl}

	current := skipList.Head

	var path string

	for current != nil && current.Data.Key != "INF" {
		path += current.Data.Key + " "
		if current.Data.Compare(node) == 0 {
			return current
		}

		// Right's data bigger than search node
		if current.Right.Data.Compare(node) == -1 {
			// we are at the lowest level
			if current.Down == nil {
				break
			}
			current = current.Down
		} else {
			current = current.Right
		}
	}

	//fmt.Println(path)
	return current
}

func (skipList *SkipList) Search(key string, ttl int) bool {
	/*
			17
			17	  25				55
			17	  25 31			    55
			17    25 31 38    44    55
		 12	17 20 25 31 38 39 44 50 55

	*/
	node := skipList.FindEntry(key, ttl)
	//fmt.Println(node)
	if node.Data.Key == key {
		return true
	}

	return false

}

func (skipList *SkipList) Insert(key string, ttl int) error {

	prevNode := skipList.FindEntry(key, ttl)

	if prevNode.Data.Key == key {
		return fmt.Errorf("This node is already present !!!")
	}

	//We are at the lowest level
	nextNode := prevNode.Right

	newNode := &Node{Data: NodeData{key, ttl}}

	prevNode.Right = newNode
	newNode.Left = prevNode

	newNode.Right = nextNode
	nextNode.Left = newNode

	skipList.NumOfElements++

	currentLevel := 0

	for {
		r, err := RandomFloat64()

		if err != nil || r <= 0.5 {
			break
		}

		newUpperNode := &Node{Data: NodeData{key, ttl}}

		// Create new level
		if currentLevel >= skipList.Height {
			fmt.Println("CUrrent Level Exceeded")
			skipList.Head.Up = &Node{Data: NodeData{"-INF", -1}}
			skipList.Head.Up.Down = skipList.Head

			skipList.Tail.Up = &Node{Data: NodeData{"INF", math.MaxInt32}}
			skipList.Tail.Up.Down = skipList.Tail

			skipList.Head = skipList.Head.Up
			skipList.Tail = skipList.Tail.Up

			skipList.Head.Right = newUpperNode
			newUpperNode.Left = skipList.Head

			newUpperNode.Right = skipList.Tail
			skipList.Tail.Left = newUpperNode

			newNode.Up = newUpperNode
			newUpperNode.Down = newNode

			skipList.Height++

			break

			// or Create duplicate node in upper level
		} else {
			fmt.Println("Find Upper Prev Elem Entered")
			prevNode = FindUpperLevelPrevElem(prevNode)
			nextNode = prevNode.Right

			prevNode.Right = newUpperNode
			newUpperNode.Left = prevNode

			newUpperNode.Right = nextNode
			nextNode.Left = newUpperNode

			newNode.Up = newUpperNode
			newUpperNode.Down = newNode

			newNode = newNode.Up
		}

		currentLevel++
	}

	return nil
}

func (skipList *SkipList) Delete(key string, ttl int) error {
	node := skipList.FindEntry(key, math.MaxInt32-1)

	if node.Data.Key != key {
		return fmt.Errorf("Can't find this key !!!")
	}

	// We get the uppermost node from FindEntry(..)
	curr := node

	for curr != nil {
		prev := curr.Left
		next := curr.Right
		down := curr.Down

		prev.Right = next
		next.Left = prev

		curr.Right = nil
		curr.Left = nil
		curr.Down = nil
		curr.Up = nil

		curr = down
	}

	return nil

}

func (skipList *SkipList) Update(key string, newTTL int) {

}

func FindUpperLevelPrevElem(prevNode *Node) *Node {
	ptr := prevNode
	// fmt.Print(ptr)
	// fmt.Printf("{Up: %v, Down: %v, Left: %v, Right: %v}\n", ptr.Up, ptr.Down, ptr.Left, ptr.Right)

	// Find first ladder left of prevNode
	for ptr != nil && ptr.Data.Key != "-INF" && ptr.Up == nil {
		ptr = ptr.Left
	}

	// Go up
	if ptr.Up != nil {
		ptr = ptr.Up
	}

	return ptr
}

func (a NodeData) Compare(b NodeData) int {
	// a > b : -1
	// a < b :   1

	if a.Key == "-INF" || b.Key == "INF" {
		return 1
	}

	if a.Key == b.Key {
		return 0
	}

	if a.TTL > b.TTL {
		return -1
	}

	if a.TTL < b.TTL {
		return 1
	}

	if a.Key > b.Key {
		return -1
	}

	return 1
}

func (skipList SkipList) Print() {
	currentHead := skipList.Head
	currentLevel := skipList.Height

	for currentHead != nil {
		curr := currentHead
		fmt.Printf("Level %v ", currentLevel)
		for curr != nil {
			fmt.Printf("> ( %v : %v ) ", curr.Data.Key, curr.Data.TTL)

			//fmt.Printf("> ( %v : {TTL: %v, Up: %v, Down: %v, Left: %v, Right %v }) ", curr.Data.Key, curr.Data.TTL, curr.Up, curr.Down, curr.Left, curr.Right)
			curr = curr.Right
		}
		fmt.Println("")
		currentLevel--
		currentHead = currentHead.Down
	}
}

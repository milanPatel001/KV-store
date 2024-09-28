package utils

import (
	"cmp"
	"fmt"
	"math"
	"sync"
	"time"
)

type NodeData[T cmp.Ordered] struct {
	Key          string
	OrderedValue T
}

type Node[T cmp.Ordered] struct {
	Data  NodeData[T]
	Up    *Node[T]
	Down  *Node[T]
	Left  *Node[T]
	Right *Node[T]
}

type SkipList[T cmp.Ordered] struct {
	Head          *Node[T]
	Tail          *Node[T]
	NumOfElements int
	Height        int
	mu            sync.RWMutex
}

type TTLSkipList struct {
	SkipList[int]
}

// TTL Skiplist
func CreateTTLSkipList() *TTLSkipList {
	HeadNode := &Node[int]{Data: NodeData[int]{"-INF", -1}}
	TailNode := &Node[int]{Data: NodeData[int]{"INF", math.MaxInt32}}

	HeadNode.Right = TailNode
	TailNode.Left = HeadNode

	return &TTLSkipList{SkipList[int]{Head: HeadNode, Tail: TailNode}}
}

func (skipList *TTLSkipList) Search(key string, ttl int) bool {
	return skipList.SkipList.Search(key, ttl)
}

func (skipList *TTLSkipList) Insert(key string, ttl int) error {
	return skipList.SkipList.Insert(key, ttl)
}

func (skipList *TTLSkipList) Delete(key string, ttl int) error {
	return skipList.SkipList.Delete(key, ttl)
}

func (skipList *TTLSkipList) Update(key string, ttl int) error {
	return skipList.SkipList.Update(key, ttl)
}

func (skipList *TTLSkipList) FindUpperLevelPrevElem(prevNode *Node[int]) *Node[int] {
	return skipList.SkipList.FindUpperLevelPrevElem(prevNode)
}

func (skipList *TTLSkipList) DeleteExpiredKeys() []string {
	skipList.mu.Lock()
	defer skipList.mu.Unlock()

	curr := skipList.Head

	// Reach to the Base Level of Head
	for curr.Down != nil {
		curr = curr.Down
	}

	curr = curr.Right

	var deletedKeys []string

	for curr != nil && curr.Data.Key != "INF" {
		if int(time.Now().Unix()) < curr.Data.OrderedValue {
			break
		}

		deletedKeys = append(deletedKeys, curr.Data.Key)

		nextBase := curr.Right

		for curr != nil {
			next := curr.Right
			prev := curr.Left

			prev.Right = next
			next.Left = prev

			curr.Right = nil
			curr.Left = nil
			curr.Down = nil

			curr = curr.Up
		}

		skipList.NumOfElements--

		// NOTE: tower cleanup will be handled by garbage collector

		curr = nextBase
	}

	return deletedKeys
}

/*
***********************
BASE SKIPLIST METHODS
***********************
*/

func (skipList *SkipList[T]) Search(key string, orderedValue T) bool {
	/*
			17
			17	  25				55
			17	  25 31			    55
			17    25 31 38    44    55
		 12	17 20 25 31 38 39 44 50 55

	*/

	skipList.mu.RLock()
	defer skipList.mu.RUnlock()

	node := skipList.FindEntry(key, orderedValue)
	//fmt.Println(node)
	if node.Data.Key == key {
		return true
	}

	return false
}

func (skipList *SkipList[T]) Insert(key string, orderedValue T) error {

	skipList.mu.Lock()
	defer skipList.mu.Unlock()

	prevNode := skipList.FindEntry(key, orderedValue)

	if prevNode.Data.Key == key {
		return fmt.Errorf("This node is already present !!!")
	}

	//We are at the lowest level
	nextNode := prevNode.Right

	newNode := &Node[T]{Data: NodeData[T]{key, orderedValue}}

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

		newUpperNode := &Node[T]{Data: NodeData[T]{key, orderedValue}}

		// Create new level
		if currentLevel >= skipList.Height {

			var minVal T
			maxVal := SetMaxValue[T]()

			skipList.Head.Up = &Node[T]{Data: NodeData[T]{"-INF", minVal}}
			skipList.Head.Up.Down = skipList.Head

			skipList.Tail.Up = &Node[T]{Data: NodeData[T]{"INF", maxVal}}
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
			//fmt.Println("Find Upper Prev Elem Entered")
			prevNode = skipList.FindUpperLevelPrevElem(prevNode)
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

func (skipList *SkipList[T]) Delete(key string, orderedValue T) error {
	skipList.mu.Lock()
	defer skipList.mu.Unlock()

	node := skipList.FindEntry(key, orderedValue)

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

	skipList.NumOfElements--

	return nil

}

func (skipList *SkipList[T]) Update(key string, newTTL int) error {
	skipList.mu.Lock()
	defer skipList.mu.Unlock()

	return nil
}

func (skipList *SkipList[T]) FindUpperLevelPrevElem(prevNode *Node[T]) *Node[T] {
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

// Returns the entry. If nothing is matched, returns the immediate smaller element in the lowest level
func (skipList *SkipList[T]) FindEntry(key string, orderedValue T) *Node[T] {
	node := NodeData[T]{key, orderedValue}

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

func (a NodeData[T]) Compare(b NodeData[T]) int {

	// a > b : -1
	// a < b :   1

	if a.Key == "-INF" || b.Key == "INF" {
		return 1
	}

	if a.Key == b.Key {
		return 0
	}

	if a.OrderedValue > b.OrderedValue {
		return -1
	}

	if a.OrderedValue < b.OrderedValue {
		return 1
	}

	if a.Key > b.Key {
		return -1
	}

	return 1
}

func (skipList *SkipList[T]) Print() {
	skipList.mu.RLock()
	defer skipList.mu.RUnlock()

	currentHead := skipList.Head
	currentLevel := skipList.Height

	for currentHead != nil {
		curr := currentHead
		fmt.Printf("Level %v ", currentLevel)
		for curr != nil {
			fmt.Printf("> ( %v : %v ) ", curr.Data.Key, curr.Data.OrderedValue)

			//fmt.Printf("> ( %v : {TTL: %v, Up: %v, Down: %v, Left: %v, Right %v }) ", curr.Data.Key, curr.Data.TTL, curr.Up, curr.Down, curr.Left, curr.Right)
			curr = curr.Right
		}
		fmt.Println("")
		currentLevel--
		currentHead = currentHead.Down
	}
}

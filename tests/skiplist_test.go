package tests

import (
	"math"
	"prac/utils"
	"testing"
)

type TTLNode = utils.Node[uint32]
type TTLNodeData = utils.NodeData[uint32]

func TestInsertBaseLevel(t *testing.T) {
	skiplist := utils.CreateTTLSkipList(48)

	skiplist.Insert("A", 60)
	skiplist.Insert("B", 20)
	skiplist.Insert("C", 50)
	skiplist.Insert("D", 12)
	skiplist.Insert("E", 20)
	skiplist.Insert("F", 10)

	curr := skiplist.Head
	for curr.Down != nil {
		curr = curr.Down
	}

	curr = curr.Right
	prevVal := curr.Left.Data.OrderedValue
	for curr != nil {
		if prevVal > curr.Data.OrderedValue {
			t.Errorf("Preval : %v is bigger than Curr val: %v. Something wrong with Skiplist Insert.", prevVal, curr.Data.OrderedValue)
			break
		}

		curr = curr.Right
	}

}

func TestExpiryDeletion(t *testing.T) {
	skipList := buildSkipList()

	skipList.DeleteExpiredKeys()

	if skipList.NumOfElements != 1 {
		t.Error("DeleteExpiredKey is not working properly !!!")
	}

	skipList.DeleteExpiredKeys()

	if skipList.NumOfElements != 1 {
		t.Error("DeleteExpiredKey is not working properly !!!")
	}

}

func TestFindUpperPrevElem(t *testing.T) {
	skipList := buildSkipList()

	result := skipList.FindUpperLevelPrevElem(skipList.Head.Right.Down.Down.Right)
	compareNodes(t, skipList.Head.Right.Down, result)

	result = skipList.FindUpperLevelPrevElem(skipList.Head.Right.Down.Down.Right.Right)
	compareNodes(t, skipList.Head.Right.Down.Right, result)

	result = skipList.FindUpperLevelPrevElem(skipList.Head.Down.Down)
	compareNodes(t, skipList.Head.Down, result)

	result = skipList.FindUpperLevelPrevElem(skipList.Head.Down.Down.Right)
	compareNodes(t, skipList.Head.Down, result)

}

func TestSearch(t *testing.T) {
	skipList := buildSkipList()
	//Print(skipList)

	result := skipList.Search("A", 17)

	if result != true {
		t.Error("A was there in skiplist, but couldn't find it.")
	}

	result = skipList.Search("D", 25)

	if result != true {
		t.Error("D was there in skiplist, but couldn't find it.")
	}

	result = skipList.Search("F", 25)

	if result != false {
		t.Error("F wasn't there in skiplist, but somehow got it.")
	}

}

func TestSkipListFindEntry(t *testing.T) {
	skipList := buildSkipList()

	node1 := &TTLNode{Data: TTLNodeData{Key: "A", OrderedValue: 12}}
	node2 := &TTLNode{Data: TTLNodeData{Key: "B", OrderedValue: 17}}
	node3 := &TTLNode{Data: TTLNodeData{Key: "C", OrderedValue: 20}}
	node4 := &TTLNode{Data: TTLNodeData{Key: "D", OrderedValue: 25}}
	node5 := &TTLNode{Data: TTLNodeData{Key: "E", OrderedValue: 1787464175}}

	result := skipList.FindEntry("A", 12)
	compareNodes(t, node1, result)

	result = skipList.FindEntry("B", 17)
	compareNodes(t, node2, result)

	result = skipList.FindEntry("C", 20)
	compareNodes(t, node3, result)

	result = skipList.FindEntry("D", 25)
	compareNodes(t, node4, result)

	result = skipList.FindEntry("E", 1787464175)
	compareNodes(t, node5, result)

	result = skipList.FindEntry("F", 27)
	compareNodes(t, node4, result)

}

func TestComapareMethods(t *testing.T) {
	node1 := TTLNodeData{Key: "F", OrderedValue: 2}
	node2 := TTLNodeData{Key: "C", OrderedValue: 20}
	node3 := TTLNodeData{Key: "J", OrderedValue: 1}
	node4 := TTLNodeData{Key: "F", OrderedValue: math.MaxInt32 - 1}
	node5 := TTLNodeData{Key: "F", OrderedValue: 1}

	result := node1.Compare(node2)

	if result != 1 {
		t.Error("node2 should come after node1")
	}

	if result = node1.Compare(node3); result != -1 {
		t.Error("node1 should come after node3")
	}

	if result = node1.Compare(node4); result != 0 {
		t.Error("node1 should be same as node4")
	}

	if result = node1.Compare(node5); result != 0 {
		t.Error("node1 should be same as node5")
	}
}

func buildSkipList() *utils.TTLSkipList {
	//
	//
	// H2 	 17			 T2
	// H1 	 17    25	 T1
	// H  12 17 20 25 44  T

	skipList := utils.CreateTTLSkipList(48)
	skipList.NumOfElements = 5
	skipList.Height = 2

	nodeH1 := &TTLNode{Data: TTLNodeData{Key: "-INF", OrderedValue: 0}}
	nodeH2 := &TTLNode{Data: TTLNodeData{Key: "-INF", OrderedValue: 0}}

	nodeT1 := &TTLNode{Data: TTLNodeData{Key: "INF", OrderedValue: math.MaxInt32}}
	nodeT2 := &TTLNode{Data: TTLNodeData{Key: "INF", OrderedValue: math.MaxInt32}}

	node1 := &TTLNode{Data: TTLNodeData{Key: "A", OrderedValue: 12}}

	node2 := &TTLNode{Data: TTLNodeData{Key: "B", OrderedValue: 17}}
	node22 := &TTLNode{Data: TTLNodeData{Key: "B", OrderedValue: 17}}
	node23 := &TTLNode{Data: TTLNodeData{Key: "B", OrderedValue: 17}}

	node3 := &TTLNode{Data: TTLNodeData{Key: "C", OrderedValue: 20}}

	node4 := &TTLNode{Data: TTLNodeData{Key: "D", OrderedValue: 25}}
	node42 := &TTLNode{Data: TTLNodeData{Key: "D", OrderedValue: 25}}

	node5 := &TTLNode{Data: TTLNodeData{Key: "E", OrderedValue: 1787464175}} // was 44

	skipList.Head.Up = nodeH1
	nodeH1.Down = skipList.Head
	nodeH1.Up = nodeH2
	nodeH2.Down = nodeH1

	skipList.Tail.Up = nodeT1
	nodeT1.Down = skipList.Tail
	nodeT1.Up = nodeT2
	nodeT2.Down = nodeT1

	skipList.Head.Right = node1
	node1.Left = skipList.Head

	node1.Right = node2
	node2.Left = node1

	// B tower
	node2.Up = node22
	node22.Down = node2
	node22.Up = node23
	node23.Down = node22

	node2.Right = node3
	node3.Left = node2

	node3.Right = node4
	node4.Left = node3

	// C tower
	node4.Up = node42
	node42.Down = node4

	node4.Right = node5
	node5.Left = node4

	node5.Right = skipList.Tail
	skipList.Tail.Left = node5

	// Express Lane 1
	nodeH1.Right = node22
	node22.Left = nodeH1
	node22.Right = node42
	node42.Left = node22
	node42.Right = nodeT1
	nodeT1.Left = node42

	// Express Lane 2
	nodeH2.Right = node23
	node23.Left = nodeH2
	node23.Right = nodeT2
	nodeT2.Left = node23

	skipList.Head = nodeH2
	skipList.Tail = nodeT2

	return skipList
}

func compareNodes(t *testing.T, expected, actual *TTLNode) {
	if expected.Data.Key != actual.Data.Key || expected.Data.OrderedValue != actual.Data.OrderedValue {
		t.Errorf("Expected node with Key: %s and TTL: %d, but got Key: %s and TTL: %d",
			expected.Data.Key, expected.Data.OrderedValue, actual.Data.Key, actual.Data.OrderedValue)
	}
}

func Print(s *utils.TTLSkipList) {
	s.Print()
}

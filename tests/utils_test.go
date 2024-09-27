package tests

import (
	"math"
	"prac/utils"
	"testing"
)

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

	result := utils.FindUpperLevelPrevElem(skipList.Head.Right.Down.Down.Right)
	compareNodes(t, skipList.Head.Right.Down, result)

	result = utils.FindUpperLevelPrevElem(skipList.Head.Right.Down.Down.Right.Right)
	compareNodes(t, skipList.Head.Right.Down.Right, result)

	result = utils.FindUpperLevelPrevElem(skipList.Head.Down.Down)
	compareNodes(t, skipList.Head.Down, result)

	result = utils.FindUpperLevelPrevElem(skipList.Head.Down.Down.Right)
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

	// result = skipList.Search("A", math.MaxInt32-1)

	// if result != true {
	// 	t.Error("A was there in skiplist, but couldn't find it.")
	// }
}

func TestSkipListFindEntry(t *testing.T) {
	skipList := buildSkipList()

	node1 := &utils.Node{Data: utils.NodeData{Key: "A", TTL: 12}}
	node2 := &utils.Node{Data: utils.NodeData{Key: "B", TTL: 17}}
	node3 := &utils.Node{Data: utils.NodeData{Key: "C", TTL: 20}}
	node4 := &utils.Node{Data: utils.NodeData{Key: "D", TTL: 25}}
	node5 := &utils.Node{Data: utils.NodeData{Key: "E", TTL: 1787464175}}

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
	node1 := utils.NodeData{Key: "F", TTL: 2}
	node2 := utils.NodeData{Key: "C", TTL: 20}
	node3 := utils.NodeData{Key: "J", TTL: 1}
	node4 := utils.NodeData{Key: "F", TTL: math.MaxInt32 - 1}
	node5 := utils.NodeData{Key: "F", TTL: 1}

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

func buildSkipList() *utils.SkipList {
	//
	//
	// H2 	 17			 T2
	// H1 	 17    25	 T1
	// H  12 17 20 25 44  T

	skipList := utils.CreateSkipList()
	skipList.NumOfElements = 5
	skipList.Height = 2

	nodeH1 := &utils.Node{Data: utils.NodeData{Key: "-INF", TTL: -1}}
	nodeH2 := &utils.Node{Data: utils.NodeData{Key: "-INF", TTL: -1}}

	nodeT1 := &utils.Node{Data: utils.NodeData{Key: "INF", TTL: math.MaxInt32}}
	nodeT2 := &utils.Node{Data: utils.NodeData{Key: "INF", TTL: math.MaxInt32}}

	node1 := &utils.Node{Data: utils.NodeData{Key: "A", TTL: 12}}

	node2 := &utils.Node{Data: utils.NodeData{Key: "B", TTL: 17}}
	node22 := &utils.Node{Data: utils.NodeData{Key: "B", TTL: 17}}
	node23 := &utils.Node{Data: utils.NodeData{Key: "B", TTL: 17}}

	node3 := &utils.Node{Data: utils.NodeData{Key: "C", TTL: 20}}

	node4 := &utils.Node{Data: utils.NodeData{Key: "D", TTL: 25}}
	node42 := &utils.Node{Data: utils.NodeData{Key: "D", TTL: 25}}

	node5 := &utils.Node{Data: utils.NodeData{Key: "E", TTL: 1787464175}} // was 44

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

func compareNodes(t *testing.T, expected, actual *utils.Node) {
	if expected.Data.Key != actual.Data.Key || expected.Data.TTL != actual.Data.TTL {
		t.Errorf("Expected node with Key: %s and TTL: %d, but got Key: %s and TTL: %d",
			expected.Data.Key, expected.Data.TTL, actual.Data.Key, actual.Data.TTL)
	}
}

func Print(s *utils.SkipList) {
	s.Print()
}

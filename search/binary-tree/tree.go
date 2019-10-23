package binaryTree

import "algorithm-class/inst"

type Tree struct {
	root *Node
	size uint
}

func NewTree() *Tree {
	return new(Tree)
}

func (t *Tree) Size() uint { return t.size }

func (t *Tree) Search(ic *inst.Counter, v int) bool {
	return binarySearch(ic, t.root, v)
}

func (t *Tree) Insert(ic *inst.Counter, v int) {
	t.root = insert(ic, t.root, v)
	t.size++
}

func binarySearch(ic *inst.Counter, n *Node, v int) bool {
	for ic.Do(inst.Compare) && n != nil {
		ic.Do(inst.Compare)
		ic.Do(inst.Indirect)
		if v < n.Value {
			ic.Do(inst.Indirect)
			n = n.Left
			continue
		}
		ic.Do(inst.Indirect)
		ic.Do(inst.Compare)
		if v > n.Value {
			ic.Do(inst.Indirect)
			n = n.Right
			continue
		}
		return true
	}
	return false
}

type Node struct {
	Value       int
	Left, Right *Node
}

func insert(ic *inst.Counter, n *Node, v int) *Node {
	ic.Do(inst.Compare)
	if n == nil {
		return &Node{Value: v}
	}
	ic.Do(inst.Compare)
	if v <= n.Value {
		ic.Use(inst.Indirect, 2)
		n.Left = insert(ic, n.Left, v)
		return n
	}
	ic.Use(inst.Indirect, 2)
	n.Right = insert(ic, n.Right, v)
	return n
}

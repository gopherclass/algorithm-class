package avl

// Operations supported:
//   1. Search
//   2. Insert
//

import (
	"algorithm-class/inst"
)

var Insts = []inst.Kind{}

type Node struct {
	Value       int
	Height      int
	Left, Right *Node
}

// Tree is a implementation of AVL Tree
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
	for ic.Once(inst.Compare) && n != nil {
		ic.Once(inst.Compare)
		ic.Once(inst.Indirect)
		if v < n.Value {
			ic.Once(inst.Indirect)
			n = n.Left
			continue
		}
		ic.Once(inst.Indirect)
		ic.Once(inst.Compare)
		if v > n.Value {
			ic.Once(inst.Indirect)
			n = n.Right
			continue
		}
		return true
	}
	return false
}

func insert(ic *inst.Counter, r *Node, v int) *Node {
	ic.Once(inst.Compare)
	if r == nil {
		return &Node{Value: v, Height: 1}
	}
	ic.Once(inst.Compare)
	ic.Once(inst.Indirect)
	if v <= r.Value {
		ic.Do(inst.Indirect, 2)
		r.Left = insert(ic, r.Left, v)
	} else {
		ic.Do(inst.Indirect, 2)
		r.Right = insert(ic, r.Right, v)
	}
	return resolve(ic, r, v)
}

func resolve(ic *inst.Counter, r *Node, v int) *Node {
	ic.Do(inst.Indirect, 2)
	w := getHeight(ic, r.Right) - getHeight(ic, r.Left)
	ic.Once(inst.Compare)
	if 2 <= w {
		ic.Once(inst.Compare)
		if v <= r.Right.Value {
			return rotrl(ic, r)
		}
		return rotl(ic, r)
	}
	ic.Once(inst.Compare)
	if w <= -2 {
		ic.Once(inst.Compare)
		if v <= r.Left.Value {
			return rotr(ic, r)
		}
		return rotlr(ic, r)
	}
	calcHeight(ic, r)
	return r
}

func max(ic *inst.Counter, x, y int) int {
	ic.Once(inst.Compare)
	if x < y {
		return y
	}
	return x
}

func getHeight(ic *inst.Counter, r *Node) int {
	ic.Once(inst.Compare)
	if r == nil {
		return 0
	}
	ic.Once(inst.Indirect)
	return r.Height
}

func recalcHeight(ic *inst.Counter, r *Node) {
	ic.Do(inst.Indirect, 2)
	calcHeight(ic, r.Left)
	calcHeight(ic, r.Right)
	calcHeight(ic, r)
}

func calcHeight(ic *inst.Counter, r *Node) {
	ic.Once(inst.Compare)
	if r == nil {
		return
	}
	ic.Do(inst.Indirect, 3)
	r.Height = 1 + max(ic, getHeight(ic, r.Left), getHeight(ic, r.Right))
}

func rotr(ic *inst.Counter, x *Node) *Node {
	ic.Do(inst.Indirect, 4)
	r := x.Left
	x.Left = r.Right
	r.Right = x
	recalcHeight(ic, r)
	return r
}

func rotl(ic *inst.Counter, x *Node) *Node {
	ic.Do(inst.Indirect, 4)
	r := x.Right
	x.Right = r.Left
	r.Left = x
	recalcHeight(ic, r)
	return r
}

func rotlr(ic *inst.Counter, x *Node) *Node {
	ic.Do(inst.Indirect, 8)
	r := x.Left.Right
	x.Left.Right = r.Left
	r.Left = x.Left
	x.Left = r.Right
	r.Right = x
	recalcHeight(ic, r)
	return r
}

func rotrl(ic *inst.Counter, x *Node) *Node {
	ic.Do(inst.Indirect, 8)
	r := x.Right.Left
	x.Right.Left = r.Right
	r.Right = x.Right
	x.Right = r.Left
	r.Left = x
	recalcHeight(ic, r)
	return r
}

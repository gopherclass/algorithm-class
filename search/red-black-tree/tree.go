package rb

// Operations supported:
//   1. Search
//   2. Insert
//

import (
	"algorithm-class/inst"
)

type Color bool

func (isRed Color) String() string {
	if isRed {
		return "red"
	}
	return "black"
}

type Node struct {
	Value       int
	Red         Color
	Left, Right *Node
}

// Tree is a implementation of Red-Black Tree
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
	ic.Do(inst.Indirect)
	t.root.Red = false
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

func insert(ic *inst.Counter, r *Node, v int) *Node {
	ic.Do(inst.Compare)
	if r == nil {
		return &Node{Value: v, Red: true}
	}
	ic.Do(inst.Indirect)
	ic.Do(inst.Compare)
	if v <= r.Value {
		ic.Use(inst.Indirect, 2)
		r.Left = insert(ic, r.Left, v)
	} else {
		ic.Use(inst.Indirect, 2)
		r.Right = insert(ic, r.Right, v)
	}
	return resolve(ic, r)
}

func resolve(ic *inst.Counter, r *Node) *Node {
	if ic.Do(inst.Indirect) && isRed(ic, r.Left) &&
		ic.Do(inst.Indirect) && isRed(ic, r.Right) {
		ic.Use(inst.Indirect, 5)
		r.Left.Red = false
		r.Right.Red = false
		r.Red = true
		return r
	}
	if ic.Do(inst.Indirect) && isRed(ic, r.Left) {
		ic.Do(inst.Indirect)
		ic.Do(inst.Indirect)
		if ic.Do(inst.Indirect) && isRed(ic, r.Left.Left) {
			return rotr(ic, r)
		}
		ic.Do(inst.Indirect)
		ic.Do(inst.Indirect)
		if ic.Do(inst.Indirect) && isRed(ic, r.Left.Right) {
			return rotlr(ic, r)
		}
		return r
	}
	if ic.Do(inst.Indirect) && isRed(ic, r.Right) {
		ic.Do(inst.Indirect)
		ic.Do(inst.Indirect)
		if ic.Do(inst.Indirect) && isRed(ic, r.Right.Right) {
			return rotl(ic, r)
		}
		ic.Do(inst.Indirect)
		ic.Do(inst.Indirect)
		if ic.Do(inst.Indirect) && isRed(ic, r.Right.Left) {
			return rotrl(ic, r)
		}
		return r
	}
	return r
}

func isRed(ic *inst.Counter, r *Node) bool {
	ic.Do(inst.Compare)
	ic.Do(inst.Trivial)
	ic.Do(inst.Indirect)
	return r != nil && bool(r.Red)
}

func div(ic *inst.Counter, r *Node) {
	ic.Do(inst.Indirect)
	r.Red = true
	ic.Do(inst.Indirect)
	if r.Left != nil {
		ic.Do(inst.Indirect)
		ic.Do(inst.Indirect)
		r.Left.Red = false
	}
	ic.Do(inst.Indirect)
	if r.Right != nil {
		ic.Do(inst.Indirect)
		ic.Do(inst.Indirect)
		r.Right.Red = false
	}
}

func rotr(ic *inst.Counter, x *Node) *Node {
	ic.Use(inst.Indirect, 4)
	r := x.Left
	x.Left = r.Right
	r.Right = x
	div(ic, r)
	return r
}

func rotl(ic *inst.Counter, x *Node) *Node {
	ic.Use(inst.Indirect, 4)
	r := x.Right
	x.Right = r.Left
	r.Left = x
	div(ic, r)
	return r
}

func rotlr(ic *inst.Counter, x *Node) *Node {
	ic.Use(inst.Indirect, 10)
	r := x.Left.Right
	x.Left.Right = r.Left
	r.Left = x.Left
	x.Left = r.Right
	r.Right = x
	div(ic, r)
	return r
}

func rotrl(ic *inst.Counter, x *Node) *Node {
	ic.Use(inst.Indirect, 10)
	r := x.Right.Left
	x.Right.Left = r.Right
	r.Right = x.Right
	x.Right = r.Left
	r.Left = x
	div(ic, r)
	return r
}

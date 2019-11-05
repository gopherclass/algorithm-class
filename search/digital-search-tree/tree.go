package digital

import (
	"algorithm-class/inst"
)

type Tree struct {
	root *Node
	size uint
	bits uint // 가장 큰 비트 번호. 4개의 비트가 있다면 3이 bits가 된다.
}

func NewTree(bits uint) *Tree {
	return &Tree{bits: bits}
}

func (t *Tree) Size() uint { return t.size }

func (t *Tree) Bits() uint { return t.bits }

func (t *Tree) Insert(ic *inst.Counter, value int) {
	t.root = insert(ic, t.root, value, t.bits)
	t.size++
}

func (t *Tree) Search(ic *inst.Counter, value int) bool {
	return search(ic, t.root, value, t.bits)
}

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

func search(ic *inst.Counter, r *Node, value int, bits uint) bool {
	x := r
	for x != nil {
		ic.Once(inst.Indirect)
		if x.Value == value {
			return true
		}
		ic.Once(inst.Indirect)
		if lookBit(value, bits) {
			ic.Once(inst.Indirect)
			x = x.Right
		} else {
			ic.Once(inst.Indirect)
			x = x.Left
		}
		bits--
	}
	return false
}

func insert(ic *inst.Counter, r *Node, value int, bits uint) *Node {
	rejectCounter(ic)

	if r == nil {
		return &Node{
			Value: value,
		}
	}
	if r.Value == value {
		return r
	}
	if lookBit(value, bits) {
		r.Right = insert(ic, r.Right, value, bits-1)
		return r
	} else {
		r.Left = insert(ic, r.Left, value, bits-1)
		return r
	}
}

func lookBit(n int, i uint) bool {
	return (n>>i)&1 == 1
}

func rejectCounter(ic *inst.Counter) {
	if ic != nil {
		panic(" is not yet implemented")
	}
}

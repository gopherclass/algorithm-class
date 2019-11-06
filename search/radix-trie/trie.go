package radix

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

func (n *Node) isInternal(ic *inst.Counter) bool {
	ic.Once(inst.Indirect)
	return n.Value == 0
}

func search(ic *inst.Counter, r *Node, value int, look uint) bool {
	x := r
	for x != nil && x.isInternal(ic) {
		if lookBit(value, look) {
			ic.Once(inst.Indirect)
			x = x.Right
		} else {
			ic.Once(inst.Indirect)
			x = x.Left
		}
		look--
	}
	if x == nil {
		return false
	}
	ic.Once(inst.Indirect)
	return x.Value == value
}

func insert(ic *inst.Counter, r *Node, value int, look uint) *Node {
	const overlook = ^uint(0)
	rejectCounter(ic)
	rejectSentinel(value)

	if r == nil {
		return &Node{
			Value: value,
		}
	}
	if !r.isInternal(ic) {
		if r.Value == value {
			return r
		}
		if look == overlook {
			return r
		}
		z := *r
		*r = Node{}
		if lookBit(z.Value, look) {
			r.Right = &z
		} else {
			r.Left = &z
		}
	}
	if lookBit(value, look) {
		r.Right = insert(ic, r.Right, value, look-1)
		return r
	} else {
		r.Left = insert(ic, r.Left, value, look-1)
		return r
	}
}

func lookBit(n int, i uint) bool {
	return (n>>i)&1 == 1
}

func rejectCounter(ic *inst.Counter) {
	if ic != nil {
		panic("counter is not yet implemented")
	}
}

func rejectSentinel(value int) {
	if value == 0 {
		panic("sentinel value zero is not allowed to be inserted into patricia tree")
	}
}

package patricia

import (
	"algorithm-class/inst"
	"unsafe"
)

type Tree struct {
	root *Node
	size uint
}

func NewTree() *Tree {
	return new(Tree)
}

func (t *Tree) Insert(ic *inst.Counter, value int) {
	var inserted bool
	t.root, inserted = insert(ic, t.root, value)
	if inserted {
		t.size++
	}
}

func (t *Tree) Search(ic *inst.Counter, value int) bool {
	return search(ic, t.root, value)
}

func (t *Tree) Size() uint {
	return t.size
}

type Node struct {
	Value       int
	Look        uint
	Left, Right *Node
}

func search(ic *inst.Counter, r *Node, value int) bool {
	var x, y *Node
	x = r
	for x != nil && (y == nil || (ic.Do(inst.Indirect, 2) && x.Look < y.Look)) {
		ic.Once(inst.Indirect)
		if lookBit(value, x.Look) {
			ic.Once(inst.Indirect)
			x, y = x.Right, x
		} else {
			ic.Once(inst.Indirect)
			x, y = x.Left, x
		}
	}
	if x == nil {
		return false
	}
	ic.Once(inst.Indirect)
	return x.Value == value
}

func insert(ic *inst.Counter, r *Node, value int) (*Node, bool) {
	rejectCounter(ic)
	rejectSentinel(value)

	var x, y *Node
	x = r
	for x != nil && (y == nil || x.Look < y.Look) {
		look := whereLook(value, x.Value)
		if look <= x.Look {
			if lookBit(value, x.Look) {
				x, y = x.Right, x
			} else {
				x, y = x.Left, x
			}
			continue
		}

		z := &Node{
			Value: value,
			Look:  look,
		}
		if y != nil {
			if y.Left == x {
				y.Left = z
			} else {
				y.Right = z
			}
		}
		if lookBit(x.Value, look) {
			z.Left = z
			z.Right = x
		} else {
			z.Left = x
			z.Right = z
		}
		if y == nil {
			return z, true
		}
		return r, true
	}
	xValue := 0
	if x != nil {
		if x.Value == value {
			return r, false
		}
		xValue = x.Value
	}
	z := &Node{
		Value: value,
		Look:  whereLook(value, xValue),
	}
	one := lookBit(value, z.Look)
	if one {
		z.Right = z
	} else {
		z.Left = z
	}
	if y == nil {
		return z, true
	}
	if y.Left == x {
		if one {
			z.Left = y.Left
		} else {
			z.Right = y.Left
		}
		y.Left = z
	} else {
		if one {
			z.Left = y.Right
		} else {
			z.Right = y.Right
		}
		y.Right = z
	}
	return r, true
}

func lookBit(n int, i uint) bool {
	return (n>>i)&1 == 1
}

func highestBit(n int) uint {
	i := unsafe.Sizeof(n) * 4
	j := i
	for j > 1 {
		j >>= 1
		if (n >> i) > 0 {
			i += j
		} else {
			i -= j
		}
	}
	if (n >> i) > 0 {
		i++
	}
	return uint(i - 1)

}

func whereLook(a int, b int) uint {
	return highestBit(a ^ b)
}

func rejectCounter(ic *inst.Counter) {
	if ic != nil {
		panic(" is not yet implemented")
	}
}

func rejectSentinel(value int) {
	if value == 0 {
		panic("sentinel value zero is not allowed to be inserted into patricia tree")
	}
}

package binaryTree

import (
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

func TestTree(t *testing.T) {
	tree := NewTree()

	rand.Seed(uint64(time.Now().UnixNano()))
	for i := 0; i < 10000; i++ {
		tree.Insert(nil, rand.Int())
	}
	checkTree(t, tree)
}

// func TestCheckExample(t *testing.T) {
// 	root := &Node{
// 		Value: 6,
// 		Left: &Node{
// 			Value: 2,
// 			Left:  &Node{Value: 1},
// 			Right: &Node{
// 				Value: 4,
// 				Left:  &Node{Value: 3},
// 				Right: &Node{Value: 5},
// 			},
// 		},
// 		Right: &Node{
// 			Value: 8,
// 			Left:  &Node{Value: 7},
// 			Right: &Node{Value: 9},
// 		},
// 	}
// 	tree := &Tree{
// 		root: root,
// 		size: 9,
// 	}
// 	checkTree(t, tree)
// }

func checkTree(t *testing.T, tree *Tree) {
	if !check(t, tree) {
		t.Fatal("wrong algorithm")
	}
}

type traverse struct {
	r, p *Node
}

func inorder(buf *[]traverse, r, p *Node) {
	if r.Left != nil {
		inorder(buf, r.Left, r)
	}
	*buf = append(*buf, traverse{r, p})
	if r.Right != nil {
		inorder(buf, r.Right, r)
	}
}

func check(t *testing.T, tree *Tree) bool {
	buf := make([]traverse, 0, tree.Size())
	inorder(&buf, tree.root, tree.root)
	ok := true
	for i, tr := range buf {
		if testing.Verbose() {
			t.Logf("key = %d, parents = %d",
				tr.r.Value, tr.p.Value)
		}
		if ok && i+1 < len(buf) && tr.r.Value > buf[i+1].r.Value {
			ok = false
		}
	}
	return ok
}

package patricia

import (
	"sort"
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

func TestTree(t *testing.T) {
	const N = 500
	rng := newrng()
	tree := NewTree()
	input := rng.Perm(N)
	for i := range input {
		input[i]++
	}
	for _, v := range input {
		tree.Insert(nil, v)
	}
	sort.Ints(input)
	for _, v := range input {
		if !tree.Search(nil, v) {
			t.Fatalf("Search(%d) returns false", v)
		}
	}
	for i := 0; i < N; i++ {
		v := 1 + rng.Intn(2*N)
		expected := binarySearchArray(input, v)
		got := tree.Search(nil, v)
		if expected != got {
			t.Fatalf("binarySearchArray = %t and tree.Search = %t",
				expected, got)
		}
	}
}

func newrng() *rand.Rand {
	seed := uint64(time.Now().UnixNano())
	rng := rand.New(rand.NewSource(seed))
	return rng
}

func binarySearchArray(s []int, v int) bool {
	if len(s) == 0 {
		return false
	}
	i, j := 0, len(s)-1
	for i < j {
		k := (i + j) / 2
		if s[k] == v {
			return true
		}
		if s[k] < v {
			i = k + 1
		} else {
			j = k - 1
		}
	}
	return s[i] == v
}

func TestSampleTree(t *testing.T) {
	if !testing.Verbose() {
		t.SkipNow()
	}
	tree := NewTree()
	for _, v := range []int{1, 19, 5, 18, 3, 26, 9} {
		tree.Insert(nil, v)
	}
	checkTree(t, tree)
}

func checkTree(t *testing.T, tree *Tree) {
	if !check(t, tree) {
		t.Fatal("wrong algorithm")
	}
}

type traverseNode struct {
	r, p *Node
}

func inorder(buf *[]traverseNode, x, y *Node) {
	if x == nil {
		return
	}
	if !isRecursive(x.Left, x) {
		inorder(buf, x.Left, x)
	}
	*buf = append(*buf, traverseNode{x, y})
	if !isRecursive(x.Right, x) {
		inorder(buf, x.Right, x)
	}
}

func isRecursive(x, y *Node) bool {
	if x == nil || y == nil {
		return true
	}
	return x.Look >= y.Look
}

func traverse(tree *Tree) []traverseNode {
	buf := make([]traverseNode, 0, tree.Size())
	inorder(&buf, tree.root, tree.root)
	return buf
}

func check(t *testing.T, tree *Tree) bool {
	if !testing.Verbose() {
		return true
	}
	buf := traverse(tree)
	for _, tr := range buf {
		t.Logf("key = %d, parents = %d",
			tr.r.Value, tr.p.Value)
	}
	return true
}

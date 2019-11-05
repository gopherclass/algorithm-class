package radix

import (
	"sort"
	"strings"
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

func TestTree(t *testing.T) {
	const bits = 5
	const N = 1 << bits
	rng := newrng()
	tree := NewTree(bits)
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
	tree := NewTree(4)
	for _, v := range []int{1, 19, 5, 18, 3, 26, 9} {
		tree.Insert(nil, v)
	}
	logTree(t, tree)
}

type traverseNode struct {
	r, p  *Node
	slope []string
}

func copySlope(s []string) []string {
	t := make([]string, len(s))
	copy(t, s)
	return t
}

func inorder(buf *[]traverseNode, slope *[]string, x, y *Node) {
	if x == nil {
		return
	}
	if x.Left != nil {
		*slope = append(*slope, "left")
		inorder(buf, slope, x.Left, x)
		*slope = (*slope)[:len(*slope)-1]
	}
	if !x.isInternal() {
		*buf = append(*buf, traverseNode{x, y, copySlope(*slope)})
	}
	if x.Right != nil {
		*slope = append(*slope, "right")
		inorder(buf, slope, x.Right, x)
		*slope = (*slope)[:len(*slope)-1]
	}
}

func traverse(tree *Tree) []traverseNode {
	buf := make([]traverseNode, 0, tree.Size())
	slope := make([]string, 0, 10)
	inorder(&buf, &slope, tree.root, tree.root)
	return buf
}

func logTree(t *testing.T, tree *Tree) {
	if !testing.Verbose() {
		return
	}
	buf := traverse(tree)
	for _, tr := range buf {
		t.Logf("%d %s", tr.r.Value, strings.Join(tr.slope, " "))
	}
}

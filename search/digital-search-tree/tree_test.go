package digital

import (
	"sort"
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

func TestTree(t *testing.T) {
	const bits = 4
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
	r, p *Node
}

func preorder(buf *[]traverseNode, x, y *Node) {
	if x == nil {
		return
	}
	*buf = append(*buf, traverseNode{x, y})
	if x.Left != nil {
		preorder(buf, x.Left, x)
	}
	if x.Right != nil {
		preorder(buf, x.Right, x)
	}
}

func traverse(tree *Tree) []traverseNode {
	buf := make([]traverseNode, 0, tree.Size())
	preorder(&buf, tree.root, tree.root)
	// TODO: 어떻게 오름차순으로 노드를 방문할 수 있는거지?
	sort.Sort(byValue(buf))
	return buf
}

func logTree(t *testing.T, tree *Tree) {
	if !testing.Verbose() {
		return
	}
	buf := traverse(tree)
	for _, tr := range buf {
		t.Logf("key = %d, parents = %d",
			tr.r.Value, tr.p.Value)
	}
}

type byValue []traverseNode

func (s byValue) Len() int      { return len(s) }
func (s byValue) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byValue) Less(i, j int) bool {
	return s[i].r.Value <= s[j].r.Value
}

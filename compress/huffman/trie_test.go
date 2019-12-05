package huffman

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	testTree(t, Build("A SIMPLE STRING TO BE ENCODED USING A MINIMAL NUMBER OF BITS"))
}

func TestProblem1(t *testing.T) {
	s := "A SIMPLE STRING TO BE ENCODED USING A MINIMAL NUMBER OF BITS"
	chars := FrequencyTable(s)
	for _, node := range chars {
		if node == nil {
			continue
		}
		t.Log(node)
	}
}

func TestProblem2(t *testing.T) {
	s := "A SIMPLE STRING TO BE ENCODED USING A MINIMAL NUMBER OF BITS"
	tree := Build(s)
	buf := flat(tree.DecodeNode)
	sort.Sort(byId(buf))
	for _, node := range buf {
		t.Log(stringentString(node))
	}
}

func TestProblem3(t *testing.T) {
	s := "A SIMPLE STRING TO BE ENCODED USING A MINIMAL NUMBER OF BITS"
	tree := Build(s)
	for _, node := range tree.Chars {
		if node == nil {
			continue
		}
		s := tree.Encode(node.Char)
		repr := s.String()
		t.Log(string(node.Char), s.Uint64(), len(repr), repr)
		assert.Equal(t, node.Char, tree.Decode(s))
	}
}

func TestQueue(t *testing.T) {
	s := byFrequency{}
	push(&s, &Node{Frequency: 4})
	push(&s, &Node{Frequency: 5})
	push(&s, &Node{Frequency: 3})
	push(&s, &Node{Frequency: 2})
	assert.Equal(t, pop(&s).Frequency, 2)
	assert.Equal(t, pop(&s).Frequency, 3)
	assert.Equal(t, pop(&s).Frequency, 4)
}

func debugPrint(node *Node) {
	adv := newAdvanceTree(node)
	adv.includeNil = true
	for adv.Next() {
		fmt.Println(adv.cur)
	}
}

func stringentString(n *Node) string {
	followId := 0
	if n.Follow != nil {
		followId = n.Follow.Id
		if n.FollowLeft {
			followId *= -1
		}
	}
	char := rune(n.Char)
	if char == 0 {
		char = '∅'
	}
	return fmt.Sprintf("(%d, %q, %d, %d)", n.Id, char, n.Frequency, followId)
}

func flat(node *Node) []*Node {
	var buf []*Node
	adv := newAdvanceTree(node)
	for adv.Next() {
		buf = append(buf, adv.cur...)
	}
	return buf
}

type byId []*Node

func (s byId) Len() int           { return len(s) }
func (s byId) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byId) Less(i, j int) bool { return s[i].Id <= s[j].Id }

func testTree(t *testing.T, tree *Tree) {
	testMin(t, tree)
	testFrequency(t, tree)
}

func testFrequency(t *testing.T, tree *Tree) {
	preorder(t, tree.DecodeNode)
}

func preorder(t *testing.T, node *Node) {
	if node == nil {
		return
	}
	assert.True(t, node.Left != nil || node.Right == nil)
	if node.Left == nil && node.Right == nil {
		return
	}
	expected := node.Left.Frequency + node.Right.Frequency
	got := node.Frequency
	if expected != got {
		t.Fatalf("%d != %d", expected, got)
	}
	preorder(t, node.Left)
	preorder(t, node.Right)
}

func testMin(t *testing.T, tree *Tree) {
	adv := newAdvanceTree(tree.DecodeNode)
	adv.Next()
	lastMin := minNode(adv.cur)
	for adv.Next() {
		min := minNode(adv.cur)
		if min > lastMin {
			t.Fatal("invalid tree")
		}
		lastMin = min
	}
}

func minNode(buf []*Node) int {
	min := buf[0].Frequency
	for _, node := range buf[1:] {
		if node.Frequency < min {
			min = node.Frequency
		}
	}
	return min
}

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

type advanceTree struct {
	cur        []*Node
	next       []*Node
	includeNil bool
}

func newAdvanceTree(cur ...*Node) advanceTree {
	return advanceTree{next: cur}
}

func (w *advanceTree) Next() bool {
	if len(w.cur) == 0 && len(w.next) > 0 {
		w.cur, w.next = w.next, w.cur
		return true
	}
	next := w.next[:0]
	for _, node := range w.cur {
		if node == nil {
			continue
		}
		if w.includeNil || node.Left != nil {
			next = append(next, node.Left)
		}
		if w.includeNil || node.Right != nil {
			next = append(next, node.Right)
		}
	}
	w.cur, w.next = next, w.cur
	return len(next) > 0
}

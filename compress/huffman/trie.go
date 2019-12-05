package huffman

import "fmt"

type Node struct {
	Id          int
	Char        byte
	Frequency   int
	FollowLeft  bool
	Follow      *Node
	Left, Right *Node
}

func (n *Node) String() string {
	if n.Frequency == 0 {
		return "∅"
	}
	return fmt.Sprintf("(%d, %c, %d)", n.Id, n.Char, n.Frequency)
}

func FrequencyTable(s string) []*Node {
	var id int
	chars := make([]*Node, 256)
	for i := 0; i < len(s); i++ {
		ref := &chars[s[i]]
		if *ref == nil {
			*ref = new(Node)
			(*ref).Id = id
			(*ref).Char = s[i]
			id++
		}
		(*ref).Frequency++
	}
	return chars
}

type Tree struct {
	Chars      []*Node
	DecodeNode *Node
}

func Build(s string) *Tree {
	chars := FrequencyTable(s)
	q := makeQueue(chars)
	id := len(q)
	for len(q) > 1 {
		x, y := pop(&q), pop(&q)
		z := &Node{
			Id:        id,
			Frequency: x.Frequency + y.Frequency,
			Left:      x,
			Right:     y,
		}
		x.FollowLeft = true
		x.Follow = z
		y.Follow = z
		push(&q, z)
		id++
	}
	decodeNode := q[0]
	return &Tree{
		Chars:      chars,
		DecodeNode: decodeNode,
	}
}

func makeQueue(chars []*Node) byFrequency {
	var n int
	for _, c := range chars {
		if c != nil && c.Frequency > 0 {
			n++
		}
	}
	q := make(byFrequency, 0, n)
	for _, c := range chars {
		if c != nil && c.Frequency > 0 {
			q = append(q, c)
		}
	}
	heapify(&q)
	return q
}

type BitString struct {
	Low *Bit
}

func NewBitString() *BitString {
	return new(BitString)
}

func (s *BitString) String() string {
	str := ""
	for bit := s.Low; bit != nil; bit = bit.Next {
		if bit.Set {
			str = "1" + str
		} else {
			str = "0" + str
		}
	}
	return str
}

func (s *BitString) Uint64() uint64 {
	w := uint64(1)
	n := uint64(0)
	for bit := s.Low; bit != nil; bit = bit.Next {
		if bit.Set {
			n += w
		}
		w *= 2
	}
	return n
}

func (s *BitString) Empty() bool {
	return s.Low == nil
}

func (s *BitString) Shift(set bool) {
	s.Low = &Bit{
		Set:  set,
		Next: s.Low,
	}
}

// 0 -> false
// 1 -> true
func (s *BitString) Unshift() bool {
	if s.Low == nil {
		panic("empty bit string")
	}
	set := s.Low.Set
	s.Low = s.Low.Next
	return set
}

type Bit struct {
	Set  bool
	Next *Bit
}

func (t *Tree) Encode(c byte) *BitString {
	node := t.Chars[c]
	if node == nil {
		return nil
	}
	s := NewBitString()
	for node.Follow != nil {
		s.Shift(!node.FollowLeft)
		node = node.Follow
	}
	return s
}

func (t *Tree) Decode(s *BitString) byte {
	node := t.DecodeNode
	for !s.Empty() {
		if s.Unshift() {
			node = node.Right
		} else {
			node = node.Left
		}
		if node == nil {
			panic("invalid huffman tree")
		}
	}
	return node.Char
}

type internalTable interface {
	Len() int
	Swap(i, j int)
	Less(i, j int) bool
	Pop() *Node
	Push(*Node)
}

type byFrequency []*Node

func (s byFrequency) Len() int {
	return len(s)
}

func (s byFrequency) Less(i, j int) bool {
	return s[i].Frequency <= s[j].Frequency
}

func (s byFrequency) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *byFrequency) Push(x *Node) {
	*s = append(*s, x)
}

func (s *byFrequency) Pop() *Node {
	n := len(*s) - 1
	x := (*s)[n]
	(*s) = (*s)[:n]
	return x
}

func heapify(s internalTable) {
	for i := s.Len() / 2; i >= 0; i-- {
		drop(s, i)
	}
}

func pop(s internalTable) *Node {
	n := s.Len() - 1
	s.Swap(0, n)
	x := s.Pop()
	drop(s, 0)
	return x
}

func push(s internalTable, x *Node) {
	s.Push(x)
	float(s, s.Len()-1)
}

func drop(s internalTable, i0 int) bool {
	n := s.Len()
	i := i0
	for i < n {
		j := 2*i + 1
		if j >= n {
			break
		}
		if j+1 < n && s.Less(j+1, j) {
			j++
		}
		if s.Less(i, j) {
			break
		}
		s.Swap(i, j)
		i = j
	}
	return i == i0
}

func float(s internalTable, i int) {
	for i > 0 {
		j := (i - 1) / 2
		if s.Less(j, i) {
			break
		}
		s.Swap(i, j)
		i = j
	}
}

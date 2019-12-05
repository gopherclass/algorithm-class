package stringSearch

const Stone = -1

type Pattern struct {
	Chars        []byte
	Next1, Next2 []int
}

type State struct {
	p  *Pattern
	q  *Deque
	s  string
	i  int
	ok bool
}

func NewState(p *Pattern, s string) *State {
	q := NewDeque()
	q.Append(Stone)
	q.Append(0)
	return &State{
		p: p,
		q: q,
		s: s,
	}
}

func (s *State) Next() bool {
	q := s.q
	v := q.Pop()
	if v == Stone {
		if q.IsEmpty() {
			return false
		}
		q.Prepend(Stone)
		s.i++
		return true
	}
	p := s.p
	if isEmptyChar(p.Chars[v]) {
		a := p.Next1[v]
		if a == 0 {
			s.ok = true
			return false
		}
		b := p.Next2[v]
		q.Append(a)
		if a != b {
			q.Append(b)
		}
		return true
	}
	if s.i >= len(s.s) {
		return false
	}
	if s.s[s.i] == p.Chars[v] {
		a := p.Next1[v]
		if a == 0 {
			s.ok = true
			return false
		}
		b := p.Next2[v]
		q.Prepend(b)
		if a != b {
			q.Prepend(a)
		}
		return true
	}
	return true
}

func isEmptyChar(c byte) bool {
	return c == ' ' || c == 0
}

func (s *State) Ok() bool {
	return s.ok
}

func (s *State) Pos() int {
	return s.i
}

type Deque struct {
	Head, Tail *Node
}

func NewDeque() *Deque {
	return new(Deque)
}

type Node struct {
	Value      int
	Prev, Next *Node
}

func (q *Deque) Pop() int {
	if q.Tail == nil {
		panic("empty deque")
	}
	v := q.Tail.Value
	q.Tail = q.Tail.Prev
	if q.Tail != nil {
		q.Tail.Next = nil
	} else {
		q.Head = nil
	}
	return v
}

func (q *Deque) Append(value int) {
	x := new(Node)
	x.Value = value
	x.Prev = q.Tail
	if q.Tail != nil {
		q.Tail.Next = x
	}
	if q.Head == nil {
		q.Head = x
	}
	q.Tail = x
}

func (q *Deque) Prepend(value int) {
	x := new(Node)
	x.Value = value
	x.Next = q.Head
	if q.Head != nil {
		q.Head.Prev = x
	}
	q.Head = x
	if q.Tail == nil {
		q.Tail = x
	}
}

func (q *Deque) IsEmpty() bool {
	return q.Head == nil
}

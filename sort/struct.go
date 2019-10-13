package main

type lesser interface {
	Less(x, y interface{}) bool
}

type lesserInt struct{}

func (lesserInt) Less(x, y interface{}) bool {
	return x.(int) <= y.(int)
}

type lesserRadixInt struct {
	bits  int
	radix int
	pow   []int
}

func (lesserRadixInt) Less(x, y interface{}) bool {
	return x.(int) <= y.(int)
}

func (r lesserRadixInt) Len() int {
	return r.bits
}

func (r lesserRadixInt) Size() int {
	return r.radix
}

func (r *lesserRadixInt) Hash(bit int, x interface{}) int {
	if r.pow == nil {
		r.pow = make([]int, r.bits)
		r.pow[0] = 1
		for i := 1; i < r.bits; i++ {
			r.pow[i] = r.pow[i-1] * r.radix
		}
	}
	return (x.(int) / r.pow[bit]) % r.radix
}

// iterativeHash는 radix sort에 사용되는 알고리즘입니다.
type iterativeHash interface {
	Len() int
	Size() int
	Hash(bit int, x interface{}) int // order = 0, 1, 2, ...
}

type mapping interface {
	Map() int
}

type mappinglesser sequenceless

func (le mappinglesser) Less(x, y interface{}) bool {
	i := x.(mapping).Map()
	j := y.(mapping).Map()
	return sequenceless(le).Less(i, j)
}

type source interface{}

type sequence interface {
	Len() int
	Peek(i int) interface{}
	Swap(i, j int)
	Set(i int, x interface{})
	Slice(i, j int) sequence
}

type vec interface {
	Push(interface{})
	Pop() interface{}
	sequence
}

type lnk interface {
	Len() int
	Pos() int
	Set(interface{})
	Peek() interface{}
	Push(interface{})
	Pop() interface{}
	Next() bool
	Prev() bool
}

func newSequence(c *sortCounter, initialCapacity int) sequence {
	if initialCapacity < 0 {
		initialCapacity = 0
	}
	v := make(asvec, 0, initialCapacity)
	return &measuresequence{
		sequence: &v,
		counter:  c,
	}
}

func newVec(c *sortCounter, initialCapacity int) vec {
	if initialCapacity < 0 {
		initialCapacity = 0
	}
	v := make(asvec, 0, initialCapacity)
	return &measurevec{
		vec:     &v,
		counter: c,
	}
}

func newLink(c *sortCounter) lnk {
	return &measurelnk{
		lnk:     new(aslnk),
		counter: c,
	}
}

func wrapSequence(s sequence, c *sortCounter) sequence {
	return &measuresequence{
		sequence: s,
		counter:  c,
	}
}

func wrapVec(v vec, c *sortCounter) vec {
	return &measurevec{
		vec:     v,
		counter: c,
	}
}

func wrapLink(l lnk, c *sortCounter) lnk {
	return &measurelnk{
		lnk:     l,
		counter: c,
	}
}

type sequenceless struct {
	sequence
	r lesser
}

func (le sequenceless) Less(i, j int) bool {
	x, y := le.Peek(i), le.Peek(j)
	return le.r.Less(x, y)
}

type vecless struct {
	vec
	r lesser
}

func (le vecless) sequenceless() sequence {
	return sequenceless{le.vec, le.r}
}

func (le vecless) Less(i, j int) bool {
	x, y := le.Peek(i), le.Peek(j)
	return le.r.Less(x, y)
}

func lessint(c *sortCounter, x, y int) bool {
	c.Less()
	return x <= y
}

func lesslnk(l lnk, r lesser) bool {
	x := l.Peek()
	if !l.Next() {
		return false
	}
	y := l.Peek()
	return r.Less(x, y)
}

func rewind(l lnk, pos int) {
	if pos < 0 {
		return
	}
	for l.Prev() {
	}
}

type asvec []interface{}

func (s asvec) Len() int                 { return len(s) }
func (s asvec) Swap(i, j int)            { s[i], s[j] = s[j], s[i] }
func (s asvec) Set(i int, x interface{}) { s[i] = x }
func (s asvec) Peek(i int) interface{}   { return s[i] }
func (s asvec) Slice(i, j int) sequence  { return s[i:j] }
func (s *asvec) Push(x interface{})      { *s = append(*s, x) }
func (s *asvec) Pop() interface{} {
	n := len(*s) - 1
	x := (*s)[n]
	*s = (*s)[:n]
	return x
}

type asints []int

func (s asints) Len() int                 { return len(s) }
func (s asints) Swap(i, j int)            { s[i], s[j] = s[j], s[i] }
func (s asints) Set(i int, x interface{}) { s[i] = x.(int) }
func (s asints) Peek(i int) interface{}   { return s[i] }
func (s asints) Slice(i, j int) sequence  { return s[i:j] }
func (s *asints) Push(x interface{})      { *s = append(*s, x.(int)) }
func (s *asints) Pop() interface{} {
	n := len(*s) - 1
	x := (*s)[n]
	*s = (*s)[:n]
	return x
}

type lnkNode struct {
	x          interface{}
	prev, next *lnkNode
}

type aslnk struct {
	len  int
	pos  int
	root *lnkNode
	cur  *lnkNode
}

func (l *aslnk) Len() int { return l.len }
func (l *aslnk) Pos() int { return l.pos }

func (l *aslnk) Peek() interface{} {
	if l.cur == nil {
		return nil
	}
	return l.cur.x
}

func (l *aslnk) Set(x interface{}) {
	if l.cur == nil {
		l.cur = new(lnkNode)
		l.root = l.cur
	}
	l.cur.x = x
}

func (l *aslnk) Push(x interface{}) {
	if l.cur == nil {
		l.cur = new(lnkNode)
		l.cur.x = x
		l.root = l.cur
		l.len++
		l.pos++
		return
	}
	if l.cur.next == nil {
		l.cur.next = &lnkNode{x: x, prev: l.cur}
	}
	l.cur = l.cur.next
	l.cur.x = x
	l.len++
	l.pos++
}

func (l *aslnk) Pop() interface{} {
	if l.cur == nil {
		return nil
	}
	if l.cur.next == nil {
		x := l.cur.x
		l.cur = l.cur.prev
		l.cur.next = nil
		if l.cur == nil {
			l.root = nil
		}
		l.len--
		l.pos--
		return x
	}
	x := l.cur.x
	l.cur = l.cur.prev
	if l.cur == nil {
		l.root = nil
	}
	l.pos--
	return x
}

func (l *aslnk) Next() bool {
	if l.cur == nil || l.cur.next == nil {
		return false
	}
	l.cur = l.cur.next
	l.pos++
	return true
}

func (l *aslnk) Prev() bool {
	if l.cur == nil || l.cur.prev == nil {
		return false
	}
	l.cur = l.cur.prev
	l.pos--
	return true
}

func (l *aslnk) Reset() {
	l.cur = l.root
	l.pos = 0
}

type aslesser struct {
	lesser  lesser
	counter *sortCounter
}

func (r *aslesser) Less(x, y interface{}) bool {
	r.counter.Less()
	return r.lesser.Less(x, y)
}

type measuresequence struct {
	sequence sequence
	counter  *sortCounter
}

func (l *measuresequence) Len() int {
	l.counter.Len()
	return l.sequence.Len()
}

func (l *measuresequence) Swap(i, j int) {
	l.counter.Swap()
	l.sequence.Swap(i, j)
}

func (l *measuresequence) Peek(i int) interface{} {
	l.counter.Peek()
	return l.sequence.Peek(i)
}

func (l *measuresequence) Set(i int, x interface{}) {
	l.counter.Set()
	l.sequence.Set(i, x)
}

func (l *measuresequence) Slice(i, j int) sequence {
	l.counter.Slice()
	return &measuresequence{
		sequence: l.sequence.Slice(i, j),
		counter:  l.counter,
	}
}

type measurevec struct {
	vec     vec
	counter *sortCounter
}

func (l *measurevec) Len() int {
	l.counter.Len()
	return l.vec.Len()
}

func (l *measurevec) Swap(i, j int) {
	l.counter.Swap()
	l.vec.Swap(i, j)
}

func (l *measurevec) Peek(i int) interface{} {
	l.counter.Peek()
	return l.vec.Peek(i)
}

func (l *measurevec) Set(i int, x interface{}) {
	l.counter.Set()
	l.vec.Set(i, x)
}

func (l *measurevec) Slice(i, j int) sequence {
	l.counter.Slice()
	return &measuresequence{
		sequence: l.vec.Slice(i, j),
		counter:  l.counter,
	}
}

func (l *measurevec) Push(x interface{}) {
	l.counter.Push()
	l.vec.Push(x)
}

func (l *measurevec) Pop() interface{} {
	l.counter.Pop()
	return l.vec.Pop()
}

type measurelnk struct {
	lnk     lnk
	counter *sortCounter
}

func (l *measurelnk) Len() int {
	l.counter.Len()
	return l.lnk.Len()
}

func (l *measurelnk) Pos() int {
	l.counter.Pos()
	return l.lnk.Pos()
}

func (l *measurelnk) Set(x interface{}) {
	l.counter.Set()
	l.lnk.Set(x)
}

func (l *measurelnk) Peek() interface{} {
	l.counter.Peek()
	return l.lnk.Peek()
}

func (l *measurelnk) Push(x interface{}) {
	l.counter.Push()
	l.lnk.Push(x)
}

func (l *measurelnk) Pop() interface{} {
	l.counter.Pop()
	return l.lnk.Pop()
}

func (l *measurelnk) Next() bool {
	l.counter.Next()
	return l.lnk.Next()
}

func (l *measurelnk) Prev() bool {
	l.counter.Prev()
	return l.lnk.Prev()
}

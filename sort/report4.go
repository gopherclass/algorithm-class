package main

func init() {
	registerSorter(naturalMergeSort{})
	registerSorter(naturalMergeSortHeap{})
	registerSorter(mergeSort{})
	registerSorter(tournamentSort{})
	registerSorter(heapSort{})
}

func makeruns(s ints) []ints {
	var rs []ints
	i := 0
	n := s.Len()
	for i < n {
		j := i + 1
		for j < n && s.Less(j-1, j) {
			j++
		}
		rs = append(rs, s.Slice(i, j))
		i = j
	}
	return rs
}

type runs struct {
	c *sortCounter
	s []ints
}

func (s runs) Len() int {
	s.c.Len()
	return len(s.s)
}

func (s runs) Swap(i, j int) {
	s.c.Swap()
	s.s[i], s.s[j] = s.s[j], s.s[i]
}

func (s runs) Less(i, j int) bool {
	s.c.Less()
	return s.s[i].Peek(0) <= s.s[j].Peek(0)
}

func (s *runs) Push(x interface{}) {
	s.c.Push()
	s.s = append(s.s, x.(ints))
}

func (s *runs) Pop() interface{} {
	s.c.Pop()
	n := s.Len() - 1
	x := s.s[n]
	s.s = s.s[:n]
	return x
}

func (s runs) Set(i int, r ints) {
	s.c.Set()
	s.s[i] = r
}

func (s runs) Peek(i int) ints {
	s.c.Peek()
	return s.s[i]
}

func merges(into *ints, rs runs) {
	heapinit(rs)
	for rs.Len() > 0 {
		r := rs.Peek(0)
		x, r := r.Peek(0), r.Slice(1, r.Len())
		into.Push(x)
		if r.Len() > 0 {
			rs.Set(0, r)
			heapdrop(rs, 0)
		} else {
			heappop(&rs)
		}
	}
}

func rmerges(c *sortCounter, rs runs, n int) ints {
	src := ints{c, make([]int, 0, n)}
	dst := ints{c, make([]int, 0, n)}
	for rs.Len() > 0 {
		r := rs.Pop().(ints)
		mergeTwo(c, &dst, src, r)
		dst, src = src, dst
		dst = dst.Slice(0, 0)
	}
	return src
}

func mergeTwo(c *sortCounter, into *ints, a, b ints) {
	m, n := a.Len(), b.Len()
	i, j := 0, 0
	for i < m && j < n {
		x := a.Peek(i)
		y := b.Peek(j)
		if c.Less() && x <= y {
			into.Push(x)
			i++
		} else {
			into.Push(y)
			j++
		}
	}
	for i < m {
		into.Push(a.Peek(i))
		i++
	}
	for j < n {
		into.Push(b.Peek(j))
		j++
	}
}

type naturalMergeSortHeap struct{}

func (naturalMergeSortHeap) epithet() string { return "natural-merge-sort-heap" }

func (naturalMergeSortHeap) sort(c *sortCounter, s []int) []int {
	res := ints{c, make([]int, 0, len(s))}
	rs := makeruns(ints{c, s})
	merges(&res, runs{c, rs})
	return res.Ints()
}

type naturalMergeSort struct{}

func (naturalMergeSort) epithet() string { return "natural-merge-sort" }

func (naturalMergeSort) sort(c *sortCounter, s []int) []int {
	rs := makeruns(ints{c, s})
	res := rmerges(c, runs{c, rs}, len(s))
	return res.Ints()
}

type mergeSort struct{}

func (mergeSort) epithet() string { return "merge-sort" }

func (sort mergeSort) sort(c *sortCounter, s []int) []int {
	return sort.do(c, ints{c, s}).Ints()
}

func (sort mergeSort) do(c *sortCounter, s ints) ints {
	n := s.Len()
	if n == 0 {
		return ints{c: c}
	}
	if n == 1 {
		return s
	}
	a := sort.do(c, s.Slice(0, n/2))
	b := sort.do(c, s.Slice(n/2, n))
	res := ints{c, make([]int, 0, a.Len()+b.Len())}
	mergeTwo(c, &res, a, b)
	return res

}

type tournamentSort struct{}

func (tournamentSort) epithet() string { return "tournament-sort" }

func (sort tournamentSort) sort(c *sortCounter, s []int) []int {
	const sentinel = 999123989
	k := msb(len(s)) << 1
	t := ints{c, make([]int, 2*k-1)}
	copydef(t.s[len(t.s)-k:], s, sentinel)

	res := make([]int, 0, len(s))
	for i := 0; i < len(s); i++ {
		sort.tournament(t, k-1)
		res = append(res, t.Peek(0))
		sort.remove(t, sentinel)
	}
	return res
}

func (sort tournamentSort) tournament(s ints, i int) {
	j := s.Len()
	for i > 0 {
		sort.runstage(s, i, j)
		i, j = (i-1)/2, i
	}
}

func (sort tournamentSort) runstage(s ints, i, j int) {
	for i < j {
		u := i
		if !s.Less(i, i+1) {
			u = i + 1
		}
		k := (i - 1) / 2
		s.Set(k, s.Peek(u))
		i += 2
	}
}

func (sort tournamentSort) remove(s ints, sentinel int) {
	v := s.Peek(0)
	i := 0
	for i < s.Len() {
		s.Set(i, sentinel)
		j := 2*i + 1
		if j >= s.Len() {
			break
		}
		if v == s.Peek(j) {
			i = j
		} else {
			i = j + 1
		}
	}
}

func (tournamentSort) servers() []serveY {
	return []serveY{
		serveAccess{},
		serveCompare{},
		serveMicrosecondLapse{},
		serveWeightedSwap{},
	}
}

func msb(n int) int {
	var r uint
	n >>= 1
	for n > 0 {
		n >>= 1
		r++
	}
	return 1 << r
}

func copydef(t, s []int, def int) {
	n := copy(t, s)
	for i := n; i < len(t); i++ {
		t[i] = def
	}
}

type heapSort struct{}

func (heapSort) epithet() string { return "heap-sort" }

func (heapSort) sort(c *sortCounter, s []int) []int {
	t := ints{c, s}
	res := make([]int, 0, len(s))
	heapinit(t)
	for t.Len() > 0 {
		res = append(res, heappop(&t).(int))
	}
	return res
}

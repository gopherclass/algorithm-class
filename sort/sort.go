package main

func init() {
	registerSorter(qsort{})
	registerSorter(naturalMergeSort{})
	registerSorter(naturalMergeSortHeap{})
}

type qsort struct{}

func (qsort) epithet() string { return "qsort" }

func (q qsort) sort(c *sortCounter, s []int) []int {
	c.Len()
	q.qsort(ints{c, s}, 0, len(s)-1)
	return s
}

func (q qsort) qsort(s ints, a, b int) {
	if a >= b {
		return
	}
	i := q.partition(s, a, b)
	q.qsort(s, a, i-1)
	q.qsort(s, i+1, b)
}

func (qsort) partition(s ints, a, b int) int {
	i, j, pv := a, b-1, b
	for {
		for i < j && s.Less(i, pv) {
			i++
		}
		for i < j && !s.Less(j, pv) {
			j--
		}
		if i >= j {
			break
		}
		s.Swap(i, j)
	}
	if !s.Less(i, pv) {
		s.Swap(i, pv)
		return i
	}
	return pv
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
		c.Swap()
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
		if x <= y {
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
	c.Len()
	res := rmerges(c, runs{c, rs}, len(s))
	return res.Ints()
}

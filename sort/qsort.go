package main

func init() {
	registerSorter(qsort{})
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

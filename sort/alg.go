//+build mage

package main

import (
	"fmt"
	"sort"
)

// qsort는 단순한 퀵정렬 알고리즘을 사용하여 정렬합니다.
type qsort struct{}

func (qsort) name() string { return "qsort" }

func (q qsort) sort(s list) {
	q.qsort(s, 0, s.Len()-1)
}
func (q qsort) qsort(s list, a, b int) {
	if a >= b {
		return
	}
	i := q.partition(s, a, b)
	q.qsort(s, a, i-1)
	q.qsort(s, i+1, b)
}

func (qsort) partition(s list, a, b int) int {
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

// iqsort는 퀵정렬 알고리즘과 삽입 정렬 알고리즘을 같이 사용하여 정렬합니다.
type iqsort struct {
	lim int
}

func (iqsort) name() string { return "iqsort" }

func (q iqsort) sort(s list) {
	q.qsort(s, 0, s.Len()-1)
}

func (q *iqsort) qsort(s list, a, b int) {
	if a >= b {
		return
	}
	if b-a <= q.lim {
		isort{}.isort(s, a, b)
		return
	}
	i := q.partition(s, a, b)
	q.qsort(s, a, i-1)
	q.qsort(s, i+1, b)
}

func (*iqsort) partition(s list, a, b int) int {
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

type mqsort struct{}

func (q mqsort) name() string { return "mqsort" }

func (q mqsort) sort(s list) {
	q.qsort(s, 0, s.Len()-1)
}

func (q mqsort) qsort(s list, a, b int) {
	if a >= b {
		return
	}
	i := q.partition(s, a, b)
	q.qsort(s, a, i-1)
	q.qsort(s, i+1, b)
}

func (mqsort) partition(s list, a, b int) int {
	c := (a + b) / 2
	i, j, pv := a+1, b-2, b-1
	mot(s, a, c, b)
	s.Swap(c, pv)
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
	if i < pv && !s.Less(i, pv) {

		s.Swap(i, pv)
		return i
	}
	return pv
}

func mot(s list, a, c, b int) {
	if b-a < 1 {
		if s.Less(b, a) {

			s.Swap(a, b)
		}
		return
	}
	if s.Less(c, a) {

		s.Swap(c, a)
	}
	if s.Less(b, c) {

		s.Swap(b, c)
	}
	if s.Less(c, a) {

		s.Swap(c, a)
	}
}

type miqsort struct {
	lim int
}

func (q miqsort) name() string { return fmt.Sprintf("miqsort(%d)", q.lim) }

func (q miqsort) sort(s list) {
	q.qsort(s, 0, s.Len()-1)
}

func (q *miqsort) qsort(s list, a, b int) {
	if a >= b {
		return
	}
	if b-a <= q.lim {
		isort{}.isort(s, a, b)
		return
	}
	i := q.partition(s, a, b)
	q.qsort(s, a, i-1)
	q.qsort(s, i+1, b)
}

func (*miqsort) partition(s list, a, b int) int {
	c := (a + b) / 2
	mot(s, a, c, b)
	i, j, pv := a+1, b-2, b-1
	s.Swap(c, pv)
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
	if i < pv && !s.Less(i, pv) {

		s.Swap(i, pv)
		return i
	}
	return pv
}

// selection sort
type ssort struct{}

func (ssort) name() string { return "ssort" }

func (ssort) sort(s list) {
	n := s.Len()
	for i := 0; i < n; i++ {
		w := i
		for j := i + 1; j < n; j++ {
			if s.Less(j, w) {
				w = j
			}
		}

		s.Swap(i, w)
	}
}

// bubble sort
type bsort struct{}

func (bsort) name() string { return "bubble sort" }

func (bsort) sort(s list) {
	for i := s.Len() - 1; i >= 0; i-- {
		for j := 0; j < i; j++ {
			if !s.Less(j, j+1) {

				s.Swap(j, j+1)
			}
		}
	}
}

// insertion sort
type isort struct{}

func (isort) name() string { return "insertion sort" }

func (i isort) sort(s list) {
	i.isort(s, 0, s.Len()-1)
}

func (isort) isort(s list, a, b int) {
	for i := a + 1; i <= b; i++ {
		if s.Less(i-1, i) {
			continue
		}
		j := i
		for j >= a+1 && !s.Less(j-1, j) {

			s.Swap(j-1, j)
			j--
		}
	}
}

// shell sort
type shellsort struct{}

func (shellsort) name() string { return "shell sort" }

func (shellsort) sort(s list) {
	h := 1
	for h < s.Len() {
		h = 3*h + 1
	}
	for h > 0 {
		for i := 0; h < s.Len() && i < h; i++ {
			for x := h + i; x < s.Len(); x += h {
				if s.Less(x-h, x) {
					continue
				}
				y := x
				for y >= h && !s.Less(y-h, y) {

					s.Swap(y-h, y)
					y -= h
				}
			}
		}
		h /= 3
	}
}

// cocktail shaker sort
type csort struct{}

func (csort) name() string { return "cocktail shaker sort" }

func (csort) sort(s list) {
	i := 0
	j := s.Len() - 1
	for i < j {
		for k := i + 1; k <= j; k++ {
			if s.Less(k, k-1) {

				s.Swap(k, k-1)
			}
		}
		j--
		for k := j - 1; i <= k; k-- {
			if s.Less(k+1, k) {

				s.Swap(k+1, k)
			}
		}
		i++
	}
}

// exchange sort
type esort struct{}

func (esort) name() string { return "exchange sort" }

func (esort) sort(s list) {
	i := 0
	j := s.Len() - 1
	for i <= j {
		for k := i + 1; k <= j; k++ {
			if s.Less(k, i) {

				s.Swap(k, i)
			}
		}
		i++
	}
}

type stdsort struct{}

func (stdsort) sort(s list) {
	sort.Sort(s)
}

type isqsort interface {
	isqsort()
}

func (qsort) isqsort()   {}
func (iqsort) isqsort()  {}
func (mqsort) isqsort()  {}
func (miqsort) isqsort() {}

type hsort struct{}

func (hsort) name() string { return "heap sort" }

func (hsort) sort(s list) {
	h := heap{s: reversed{s}, n: s.Len()}
	h.init()
	h.sort()
}

type reversed struct {
	list list
}

func (r reversed) Len() int           { return r.list.Len() }
func (r reversed) Swap(i, j int)      { r.list.Swap(i, j) }
func (r reversed) Less(i, j int) bool { return !r.list.Less(i, j) }
func (r reversed) Peek(i int) int     { return r.list.Peek(i) }

type heap struct {
	s list
	n int
}

func (h *heap) init() {
	for i := h.n / 2; i >= 0; i-- {
		heapify(h.s, i, h.n)
	}
}

func (h *heap) pop() {
	h.n--
	h.s.Swap(0, h.n)
	heapify(h.s, 0, h.n)
}

func (h *heap) sort() {
	for h.n >= 1 {
		h.pop()
	}
}

func heapify(s list, x, n int) bool {
	i := x
	for {
		j := 2*i + 1
		if n <= j || j < 0 {
			break
		}
		if j+1 < n && !s.Less(j, j+1) {
			j = j + 1
		}
		if s.Less(i, j) {
			break
		}
		s.Swap(i, j)
		i = j
	}
	return x < i
}


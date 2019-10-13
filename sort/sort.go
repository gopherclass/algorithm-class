package main

import (
	"fmt"
	"runtime"
)

func init() {
	registerSorter(qsort{})
	registerSorter(naturalMergeSort{})
	registerSorter(naturalMergeSortHeap{})
}

type run struct {
	i, j int
}

func (r run) Len() int { return r.j - r.i }
func (r run) Map() int { return r.i }

// T([]run, mappinglesser) -> sequenceless

func makeruns(s sequenceless, c *sortCounter) vec {
	v := newVec(c, 0)
	i, n := 0, s.Len()
	for i < n {
		j := i + 1
		for j < n && s.Less(j-1, j) {
			j++
		}
		v.Push(run{i, j})
		i = j
	}
	return v
}

func mergeRuns(rs vec, s mappinglesser, c *sortCounter) vec {
	sorted := newVec(c, s.Len())
	h := vecless{rs, s}
	heapinit(h)
	for rs.Len() > 0 {
		r := rs.Peek(0).(run)
		sorted.Push(s.Peek(r.i))
		r.i++
		if r.i < r.j {
			rs.Set(0, r)
			heapdrop(h, 0)
		} else {
			heappop(h)
		}
	}
	return sorted
}

func makefullruns(s sequenceless, c *sortCounter) vec {
	rs := newVec(c, 0)
	i, n := 0, s.Len()
	for i < n {
		r := newVec(c, 2)
		r.Push(s.Peek(i))
		j := i + 1
		for j < n && s.Less(j-1, j) {
			r.Push(s.Peek(j))
			j++
		}
		rs.Push(r)
		i = j
	}
	return rs
}

func recursiveMergeRuns(rs vec, le lesser, c *sortCounter) vec {
	n := rs.Len()
	if n == 0 {
		return newVec(c, 0)
	}
	if n <= 1 {
		return rs.Pop().(vec)
	}
	rs.Push(mergeTwoFullRuns(rs, le, c))
	return recursiveMergeRuns(rs, le, c)
}

func mergeTwoFullRuns(rs vec, le lesser, c *sortCounter) vec {
	r0 := rs.Pop().(vec)
	r1 := rs.Pop().(vec)
	m, n := r0.Len(), r1.Len()
	i, j := 0, 0
	v := newVec(c, m+n)
	for c.Less() && i < m && c.Less() && j < n {
		x, y := r0.Peek(i), r1.Peek(j)
		if le.Less(x, y) {
			v.Push(x)
			i++
		} else {
			v.Push(y)
			j++
		}
	}
	for c.Less() && i < m {
		v.Push(r0.Peek(i))
		i++
	}
	for c.Less() && j < n {
		v.Push(r1.Peek(j))
		j++
	}
	return v
}

type naturalMergeSortHeap struct{}

func (naturalMergeSortHeap) epithet() string { return "natural-merge-sort-heap" }

func (naturalMergeSortHeap) sort(s sequence, r lesser, c *sortCounter) source {
	rs := makeruns(sequenceless{s, r}, c)
	return mergeRuns(rs, mappinglesser{s, r}, c)
}

type naturalMergeSort struct{}

func (naturalMergeSort) epithet() string { return "natural-merge-sort" }

func (naturalMergeSort) sort(s sequence, r lesser, c *sortCounter) source {
	rs := makefullruns(sequenceless{s, r}, c)
	sorted := recursiveMergeRuns(rs, r, c)
	if !isSorted(sorted, r) {
		runtime.Breakpoint()
		fmt.Println(naturalMergeSort{}.sort(s, r, c))
	}
	return sorted

}

type qsort struct{}

func (qsort) epithet() string { return "qsort" }

func (q qsort) sort(s sequence, r lesser, c *sortCounter) source {
	q.qsort(sequenceless{s, r}, 0, s.Len()-1)
	return s
}

func (q qsort) qsort(s sequenceless, a, b int) {
	if a >= b {
		return
	}
	i := q.partition(s, a, b)
	q.qsort(s, a, i-1)
	q.qsort(s, i+1, b)
}

func (qsort) partition(s sequenceless, a, b int) int {
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

// func main() {
// 	var a algs
// 	a.alg("selection sort", ssort{}, 600, 3, 300)
// 	a.alg("bubble sort", bsort{}, 600, 3, 300)
// 	a.alg("cocktail shaker sort", csort{}, 600, 3, 300)
// 	a.alg("exchange sort", esort{}, 600, 3, 300)
// 	a.runTests()
// }
//
// func main2() {
// 	var a algs
// 	const iteration = 3
// 	// a.alg("selection sort", ssort{}, 300, iteration, 200)
// 	// a.alg("bubble sort", bsort{}, 300, iteration, 200)
// 	// a.alg("insertion sort", isort{}, 500, iteration, 200)
// 	// a.alg("shell sort", shellsort{}, 500, iteration, 500)
// 	// a.alg("quick sort", qsort{}, 500, iteration, 500)
// 	// a.alg("insertion sort(M=10) + quick sort", iqsort{lim: 10}, 500, iteration, 500)
// 	// a.alg("median of three + quick sort", mqsort{}, 500, iteration, 500)
// 	// a.alg("median of three + insertion(M=10) + quick sort", miqsort{lim: 10}, 500, iteration, 500)
// 	for m := 3; m <= 40; m++ {
// 		a.test(fmt.Sprintf("insertion(M=%d) + quick sort", m), iqsort{lim: m}, 500, 200)
// 	}
// 	for m := 3; m <= 20; m++ {
// 		a.test(fmt.Sprintf("insertion(M=%d) + quick sort", m), iqsort{lim: m}, 500, 2000)
// 	}
// 	a.runTests()
// 	// a.run()
// 	// a.runDraw()
// }

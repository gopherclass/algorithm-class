//+build mage

package main

// aqsort는 Quick Sort 알고리즘을 공격하는 데이터 입력을 찾아낸다. 이
// 알고리즘은 아래 링크에 c로 작성되어 있는 것을 포팅하였다.
//
// M. Douglas McIlroy, A Killer Adversary for Quicksort, Dartmouth College,
// https://www.cs.dartmouth.edu/~doug/mdmspe.pdf
//
// M. Douglas McIlroy, https://www.cs.dartmouth.edu/~doug/aqsort.c
type aqsort struct {
	gas       int
	nsolid    int
	candidate int
	ptr       []int
	poison    []int
}

func (a *aqsort) aqsort(sorter sorter, n int) []int {
	a.gas = n - 1
	a.nsolid = 0
	a.candidate = 0
	a.poison = make([]int, n)
	a.ptr = make([]int, n)
	for i := range a.poison {
		a.poison[i] = a.gas
		a.ptr[i] = i
	}
	sorter.sort(a)
	return a.poison
}

func (a *aqsort) Len() int {
	return len(a.poison)
}

func (a *aqsort) Swap(i, j int) {
	a.ptr[i], a.ptr[j] = a.ptr[j], a.ptr[i]
}

func (a *aqsort) Less(i, j int) bool {
	x, y := a.ptr[i], a.ptr[j]
	if a.poison[x] == a.gas && a.poison[y] == a.gas {
		if x == a.candidate {
			a.freeze(x)
		} else {
			a.freeze(y)
		}
	}
	if a.poison[x] == a.gas {
		a.candidate = x
	} else if a.poison[y] == a.gas {
		a.candidate = y
	}
	return a.poison[x] <= a.poison[y]
}

func (a *aqsort) Peek(i int) int {
	panic("unimplemented Peek")
}

func (a *aqsort) freeze(i int) {
	a.poison[i] = a.nsolid
	a.nsolid++
}

func antiqsort(sorter sorter, n int) []int {
	var a aqsort
	return a.aqsort(sorter, n)

}

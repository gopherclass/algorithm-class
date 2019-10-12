package main

import (
	"fmt"
	"io"
	"os"
)

type testResult struct {
	pass      bool
	stat      sortStat
	errinput  []int
	erroutput []int
}

func testSort(sorter sorter) testResult {
	const maxiteration = 200
	var result testResult
	var rawinput = make([]int, 0, maxiteration)
	for i := uint(0); i <= maxiteration; i++ {
		size := i
		input := fuzzInput(size)(i)
		rawinput = append(rawinput[:0], input...)
		sorted, r, counter := measureSort(sorter, input)
		accCounter(&result.stat, counter)
		if !isSorted(sorted, r) {
			result.errinput = rawinput
			result.erroutput = input
			return result
		}
	}
	result.pass = true
	return result
}

func isSorted(src source, r lesser) bool {
	switch src := src.(type) {
	case sequence:
		return isSortedSequence(src, r)
	case lnk:
		return isSortedLink(src, r)
	default:
		panic("invalid type")
	}
}

func isSortedSequence(s sequence, r lesser) bool {
	n := s.Len()
	if n == 0 {
		return true
	}
	l := sequenceless{s, r}
	for i := 1; i < n; i++ {
		if !l.Less(i-1, i) {
			return false
		}
	}
	return true
}

func isSortedLink(l lnk, r lesser) bool {
	rewind(l, -1)
	n := l.Len()
	for i := 1; i < n; i++ {
		if !lesslnk(l, r) {
			return false
		}
	}
	return true
}

func runTest(sorter sorter) bool {
	r := testSort(sorter)
	if !r.pass {
		testFail(sorter, showcase{
			stat:      r.stat,
			errinput:  r.errinput,
			erroutput: r.erroutput,
		})
		return false
	}
	testPass(sorter, showcase{
		stat: r.stat,
	})
	return true
}

func testFail(sorter sorter, showcase showcase) {
	showTest(os.Stderr, "Fail", sorter, showcase)
}

func testPass(sorter sorter, showcase showcase) {
	showTest(os.Stdout, "OK", sorter, showcase)
}

type showcase struct {
	stat      sortStat
	errinput  []int
	erroutput []int
}

func showTest(w io.Writer, verb string, sorter sorter, showcase showcase) {
	fmt.Fprintf(w, "%s %s, len = %.2f, compare = %.2f, swap = %.2f, peek = %.2f, time = %s",
		verb,
		sorter.epithet(),
		showcase.stat.averageLen,
		showcase.stat.averageLess,
		showcase.stat.averageSwap,
		showcase.stat.averagePeek,
		showcase.stat.averageLapse.String(),
	)
	if showcase.errinput != nil {
		fmt.Fprintf(w, ", input = %#v", showcase.errinput)
	}
	if showcase.erroutput != nil {
		fmt.Fprintf(w, ", got = %#v", showcase.erroutput)
	}
	fmt.Fprintln(w)
}

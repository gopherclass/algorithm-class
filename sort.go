//+build mage

package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"golang.org/x/exp/rand"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const imageWidth = 13 * vg.Centimeter
const imageHeight = 4.63 * vg.Centimeter

type slice interface {
	Len() int
	Swap(i, j int)
	Less(i, j int) bool // s[j] <= s[j]
}

type list interface {
	Len() int
	Swap(i, j int)
	Less(i, j int) bool
	Peek(i int) int
}

type sorter interface {
	sort(ic *instCounter, s slice)
}

type asc []int

func (s asc) Len() int           { return len(s) }
func (s asc) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s asc) Less(i, j int) bool { return s[i] <= s[j] }

type byComparison struct {
	list list
}

type instCounter struct {
	insts map[string]uint
}

func newCounter() *instCounter {
	return &instCounter{
		insts: make(map[string]uint),
	}
}

func (i *instCounter) inc(kind string) bool {
	if i == nil {
		return true
	}
	i.insts[kind]++
	return true
}

func (i *instCounter) inst(kind string) uint {
	return i.insts[kind]
}

func (i *instCounter) reset() {
	if i == nil {
		return
	}
	for key := range i.insts {
		delete(i.insts, key)
	}
}

type qsort struct{}

func (q qsort) sort(ic *instCounter, s slice) {
	q.qsort(ic, s, 0, s.Len()-1)
}
func (q qsort) qsort(ic *instCounter, s slice, a, b int) {
	if a >= b {
		return
	}
	i := q.partition(ic, s, a, b)
	q.qsort(ic, s, a, i-1)
	q.qsort(ic, s, i+1, b)
}

func (qsort) partition(ic *instCounter, s slice, a, b int) int {
	i, j, pv := a, b-1, b
	for {
		for i < j && ic.inc("compare") && s.Less(i, pv) {
			i++
		}
		for i < j && ic.inc("compare") && !s.Less(j, pv) {
			j--
		}
		if i >= j {
			break
		}
		ic.inc("swap")
		s.Swap(i, j)
	}
	if ic.inc("compare") && !s.Less(i, pv) {
		ic.inc("swap")
		s.Swap(i, pv)
		return i
	}
	return pv
}

// 중간값 분할 qsort
type iqsort struct {
	lim int
}

func (q iqsort) sort(ic *instCounter, s slice) {
	q.qsort(ic, s, 0, s.Len()-1)
}

func (q *iqsort) qsort(ic *instCounter, s slice, a, b int) {
	if a >= b {
		return
	}
	if b-a <= q.lim {
		isort{}.isort(ic, s, a, b)
		return
	}
	i := q.partition(ic, s, a, b)
	q.qsort(ic, s, a, i-1)
	q.qsort(ic, s, i+1, b)
}

func (*iqsort) partition(ic *instCounter, s slice, a, b int) int {
	i, j, pv := a, b-1, b
	for {
		for i < j && ic.inc("compare") && s.Less(i, pv) {
			i++
		}
		for i < j && ic.inc("compare") && !s.Less(j, pv) {
			j--
		}
		if i >= j {
			break
		}
		ic.inc("swap")
		s.Swap(i, j)
	}
	if ic.inc("compare") && !s.Less(i, pv) {
		ic.inc("swap")
		s.Swap(i, pv)
		return i
	}
	return pv
}

type mqsort struct{}

func (q mqsort) sort(ic *instCounter, s slice) {
	q.qsort(ic, s, 0, s.Len()-1)
}

func (q mqsort) qsort(ic *instCounter, s slice, a, b int) {
	if a >= b {
		return
	}
	i := q.partition(ic, s, a, b)
	q.qsort(ic, s, a, i-1)
	q.qsort(ic, s, i+1, b)
}

func (mqsort) partition(ic *instCounter, s slice, a, b int) int {
	c := (a + b) / 2
	i, j, pv := a+1, b-2, b-1
	mot(ic, s, a, c, b)
	s.Swap(c, pv)
	for {
		for i < j && ic.inc("compare") && s.Less(i, pv) {
			i++
		}
		for i < j && ic.inc("compare") && !s.Less(j, pv) {
			j--
		}
		if i >= j {
			break
		}
		ic.inc("swap")
		s.Swap(i, j)
	}
	if i < pv && ic.inc("compare") && !s.Less(i, pv) {
		ic.inc("swap")
		s.Swap(i, pv)
		return i
	}
	return pv
}

func mot(ic *instCounter, s slice, a, c, b int) {
	if b-a < 1 {
		if ic.inc("compare") && s.Less(b, a) {
			ic.inc("swap")
			s.Swap(a, b)
		}
		return
	}
	if ic.inc("compare") && s.Less(c, a) {
		ic.inc("swap")
		s.Swap(c, a)
	}
	if ic.inc("compare") && s.Less(b, c) {
		ic.inc("swap")
		s.Swap(b, c)
	}
	if ic.inc("compare") && s.Less(c, a) {
		ic.inc("swap")
		s.Swap(c, a)
	}
}

type miqsort struct {
	lim int
}

func (q miqsort) sort(ic *instCounter, s slice) {
	q.qsort(ic, s, 0, s.Len()-1)
}

func (q *miqsort) qsort(ic *instCounter, s slice, a, b int) {
	if a >= b {
		return
	}
	if b-a <= q.lim {
		isort{}.isort(ic, s, a, b)
		return
	}
	i := q.partition(ic, s, a, b)
	q.qsort(ic, s, a, i-1)
	q.qsort(ic, s, i+1, b)
}

func (*miqsort) partition(ic *instCounter, s slice, a, b int) int {
	c := (a + b) / 2
	mot(ic, s, a, c, b)
	i, j, pv := a+1, b-2, b-1
	s.Swap(c, pv)
	for {
		for i < j && ic.inc("compare") && s.Less(i, pv) {
			i++
		}
		for i < j && ic.inc("compare") && !s.Less(j, pv) {
			j--
		}
		if i >= j {
			break
		}
		ic.inc("swap")
		s.Swap(i, j)
	}
	if i < pv && ic.inc("compare") && !s.Less(i, pv) {
		ic.inc("swap")
		s.Swap(i, pv)
		return i
	}
	return pv
}

func (qsort) isqsort()   {}
func (iqsort) isqsort()  {}
func (mqsort) isqsort()  {}
func (miqsort) isqsort() {}

// selection sort
type ssort struct{}

func (ssort) sort(ic *instCounter, s slice) {
	for i := 0; i < s.Len(); i++ {
		w := i
		for j := i + 1; j < s.Len(); j++ {
			if ic.inc("compare") && s.Less(j, w) {
				w = j
			}
		}
		ic.inc("swap")
		s.Swap(i, w)
	}
}

// bubble sort
type bsort struct{}

func (bsort) sort(ic *instCounter, s slice) {
	for i := s.Len() - 1; i >= 0; i-- {
		for j := 0; j < i; j++ {
			if ic.inc("compare") && !s.Less(j, j+1) {
				ic.inc("swap")
				s.Swap(j, j+1)
			}
		}
	}
}

// insertion sort
type isort struct{}

func (i isort) sort(ic *instCounter, s slice) {
	i.isort(ic, s, 0, s.Len()-1)
}

func (isort) isort(ic *instCounter, s slice, a, b int) {
	for i := a + 1; i <= b; i++ {
		if ic.inc("compare") && s.Less(i-1, i) {
			continue
		}
		j := i
		for j >= a+1 && ic.inc("compare") && !s.Less(j-1, j) {
			ic.inc("swap")
			s.Swap(j-1, j)
			j--
		}
	}
}

// shell sort
type shellsort struct{}

func (shellsort) sort(ic *instCounter, s slice) {
	h := 1
	for h < s.Len() {
		h = 3*h + 1
	}
	for h > 0 {
		for i := 0; h < s.Len() && i < h; i++ {
			for x := h + i; x < s.Len(); x += h {
				if ic.inc("compare") && s.Less(x-h, x) {
					continue
				}
				y := x
				for y >= h && ic.inc("compare") && !s.Less(y-h, y) {
					ic.inc("swap")
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

func (csort) sort(ic *instCounter, s slice) {
	i := 0
	j := s.Len() - 1
	for i < j {
		for k := i + 1; k <= j; k++ {
			if ic.inc("compare") && s.Less(k, k-1) {
				ic.inc("swap")
				s.Swap(k, k-1)
			}
		}
		j--
		for k := j - 1; i <= k; k-- {
			if ic.inc("compare") && s.Less(k+1, k) {
				ic.inc("swap")
				s.Swap(k+1, k)
			}
		}
		i++
	}
}

// exchange sort
type esort struct{}

func (esort) sort(ic *instCounter, s slice) {
	i := 0
	j := s.Len() - 1
	for i <= j {
		for k := i + 1; k <= j; k++ {
			if ic.inc("compare") && s.Less(k, i) {
				ic.inc("swap")
				s.Swap(k, i)
			}
		}
		i++
	}
}

type stdqsort struct{}

func (stdqsort) sort(ic *instCounter, s slice) {
	sort.Sort(&stdasc{ic: ic, s: s})
}

type stdasc struct {
	ic *instCounter
	s  slice
}

func (a *stdasc) Len() int { return a.s.Len() }
func (a *stdasc) Swap(i, j int) {
	a.ic.inc("swap")
	a.s.Swap(i, j)
}

func (a *stdasc) Less(i, j int) bool {
	a.ic.inc("compare")
	return a.s.Less(i, j)
}

type label []string

func (buf label) append(s string) label {
	return append(buf, s)
}

func (buf *label) String() string {
	return strings.Join([]string(*buf), " - ")
}

type ifqsort interface {
	isqsort()
}

type sizedInputFunc func(int) inputFunc

func benchmark(label label, sorter sorter, maxsize int, iteration uint) {
	swappedSorted := func(swapFactor float64) sizedInputFunc {
		return func(size int) inputFunc {
			swap := int(float64(size) * swapFactor)
			return almostSortedInput(size, swap)
		}
	}
	type inputType struct {
		name      string
		makeinput sizedInputFunc
	}
	var inputs = []inputType{
		{"fuzz input", fuzzInput},
		{"sorted input", sortedInput},
		{"reversed input", reversedInput},
		// {"almost sorted input (0.75 swapped)", swappedSorted(0.75)},
		{"almost sorted input (0.1 swapped)", swappedSorted(0.10)},
		// {"almost sorted input (0.25 swapped)", swappedSorted(0.25)},
		// {"almost sorted input (0.125 swapped)", swappedSorted(0.125)},
	}
	if _, ok := sorter.(ifqsort); ok {
		killqsort := func(size int) inputFunc {
			killer := antiqsort(nil, sorter, size)
			return constInput(killer)
		}
		inputs = append(inputs, inputType{"killing input", killqsort})
	}
	for _, input := range inputs {
		benchmarkInput(
			label.append(input.name),
			sorter,
			maxsize,
			iteration,
			input.makeinput,
		)
	}
}

func benchmarkInput(label label, sorter sorter, maxsize int, iteration uint, makeinput sizedInputFunc) {
	samples := benchmarkSort(sorter, maxsize, iteration, makeinput)
	start := func(serve serveAxis) {
		pl, err := plotSamples(label, samples, serve)
		if err != nil {
			log.Printf("plotSamples(): %s %+v", label.String(), err)
			return
		}
		err = plotSave(label.append(serve.label()), pl)
		if err != nil {
			log.Printf("plotSave() %s %+v", label.String(), err)
		}
	}
	start(serveCompare{})
	start(serveSwap{})
	start(serveLapse{})
}

type serveCompare struct{}

func (serveCompare) label() string             { return "compare" }
func (serveCompare) serve(s *sortStat) float64 { return mean(s.compare, s.iteration) }

type serveSwap struct{}

func (serveSwap) label() string             { return "swap" }
func (serveSwap) serve(s *sortStat) float64 { return mean(s.swap, s.iteration) }

type serveLapse struct{}

func (serveLapse) label() string { return "time" }
func (serveLapse) serve(s *sortStat) float64 {
	return float64(convMicroseconds(meanDuration(s.lapse, s.iteration)))
}

func benchmarkSort(sorter sorter, maxsize int, iteration uint, makeinput sizedInputFunc) []sortStat {
	samples := make([]sortStat, 0, maxsize+1)
	for size := 0; size <= maxsize; size++ {
		inputFunc := makeinput(size)
		var stat sortStat
		stat.inputs = size
		for i := uint(1); i <= iteration; i++ {
			orig := inputFunc(i)
			counter := measureSort(sorter, asc(orig))
			stat.accumulate(counter)
		}
		samples = append(samples, stat)
	}
	return samples
}

type serveAxis interface {
	label() string
	serve(*sortStat) float64
}

func plotSamples(label label, samples []sortStat, serve serveAxis) (*plot.Plot, error) {
	pl, err := plot.New()
	if err != nil {
		return nil, err
	}
	pl.Title.Text = label.String()
	adjustFontSize(&pl.Title.Font, pl.Title.Text, 0.6*vg.Centimeter)
	pl.X.Label.Text = "Size"
	pl.X.Label.Font.Size = 0.6 * vg.Centimeter
	pl.Y.Label.Text = serve.label()
	pl.Y.Label.Font.Size = 0.6 * vg.Centimeter
	// TODO: comma tick?
	// pl.Y.Tick.Marker = ?

	plotFunc(pl, 1, "y = x", identifyFunc)
	plotFunc(pl, 2, "y = x²", quadraticFunc)
	plotFunc(pl, 3, "y = xlog(x)", xlogxFunc)

	line, _, err := plotter.NewLinePoints(serveSamples(samples, serve))
	if err != nil {
		return nil, err
	}
	line.Width = 0.12 * vg.Centimeter
	line.Color = plotutil.Color(0)
	pl.Add(line)
	pl.Legend.Add(serve.label(), line)

	return pl, nil
}

func adjustFontSize(font *vg.Font, s string, size vg.Length) {
	font.Size = size
	for {
		w := font.Width(s)
		if w <= imageWidth {
			return
		}
		font.Size -= 0.5
	}
}

func serveSamples(samples []sortStat, serve serveAxis) plotter.XYs {
	xys := make(plotter.XYs, len(samples))
	for i := range samples {
		stat := &samples[i]
		xys[i] = plotter.XY{
			X: float64(stat.inputs),
			Y: serve.serve(stat),
		}
	}
	return xys
}

func plotSave(label label, pl *plot.Plot) error {
	dir := label[0]
	os.Mkdir(dir, os.ModePerm)
	name := fmt.Sprintf("%s.jpg", label.String())
	path := filepath.Join(dir, name)
	return pl.Save(imageWidth, imageHeight, path)
}

func plotFunc(pl *plot.Plot, style int, legend string, fn func(float64) float64) {
	f := plotter.NewFunction(fn)
	// f.Dashes = plotutil.Dashes(style)
	f.Color = plotutil.Color(style)
	f.Width = 0.1 * vg.Centimeter
	pl.Add(f)
	pl.Legend.Add(legend, f)
}

func identifyFunc(x float64) float64 { return x }

func quadraticFunc(x float64) float64 { return x * x }

func xlogxFunc(x float64) float64 {
	if x <= 1 {
		return 0
	}
	return x * math.Log2(x)
}

type algs struct {
	tests []func() bool
	draws []func()
}

func (a *algs) alg(name string, sorter sorter, maxsize int, iteration uint, first int) {
	a.tests = append(a.tests, func() bool {
		return testSort(name, sorter, 100, fuzzInput(maxsize)) &&
			testSort(name, sorter, 100, fuzzInput(first)) &&
			testSort(name, sorter, 100, fuzzInput(first*2)) &&
			testSort(name, sorter, 100, fuzzInput(first*3))

	})
	a.draws = append(a.draws, func() {
		startTime := time.Now()
		benchmark(label{name}, sorter, maxsize, iteration)
		lapse := time.Since(startTime)
		fmt.Printf("OK %s: %s\n", name, lapse.String())
	})
}

func (a *algs) test(name string, sorter sorter, maxsize int, iteration uint) {
	a.tests = append(a.tests, func() bool {
		return testSort(name, sorter, iteration, fuzzInput(maxsize))
	})
}

func (a *algs) run() {
	if a.runTests() {
		a.runDraw()
	}
}

func (a *algs) runTests() bool {
	for _, test := range a.tests {
		ok := test()
		if !ok {
			return false
		}
	}
	return true
}

func (a *algs) runDraw() {
	for _, draw := range a.draws {
		draw()
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var a algs
	a.alg("selection sort", ssort{}, 600, 3, 300)
	a.alg("bubble sort", bsort{}, 600, 3, 300)
	a.alg("cocktail shaker sort", csort{}, 600, 3, 300)
	a.alg("exchange sort", esort{}, 600, 3, 300)
	a.runTests()
}

func main2() {
	var a algs
	const iteration = 3
	// a.alg("selection sort", ssort{}, 300, iteration, 200)
	// a.alg("bubble sort", bsort{}, 300, iteration, 200)
	// a.alg("insertion sort", isort{}, 500, iteration, 200)
	// a.alg("shell sort", shellsort{}, 500, iteration, 500)
	// a.alg("quick sort", qsort{}, 500, iteration, 500)
	// a.alg("insertion sort(M=10) + quick sort", iqsort{lim: 10}, 500, iteration, 500)
	// a.alg("median of three + quick sort", mqsort{}, 500, iteration, 500)
	// a.alg("median of three + insertion(M=10) + quick sort", miqsort{lim: 10}, 500, iteration, 500)
	for m := 3; m <= 40; m++ {
		a.test(fmt.Sprintf("insertion(M=%d) + quick sort", m), iqsort{lim: m}, 500, 200)
	}
	for m := 3; m <= 20; m++ {
		a.test(fmt.Sprintf("insertion(M=%d) + quick sort", m), iqsort{lim: m}, 500, 2000)
	}
	a.runTests()
	// a.run()
	// a.runDraw()
}

// const SEED = 123712381313881319783872
var SEED = time.Now().UnixNano()

var rngSource = rand.New(rand.NewSource(uint64(SEED)))

type inputFunc func(size uint) []int

func testSort(name string, sorter sorter, iteration uint, inputFunc inputFunc) bool {
	var stat sortStat
	var counters []sortCounter
	var orig []int
	for i := uint(0); i <= iteration; i++ {
		data := inputFunc(i)
		orig = append(orig[:0], data...)
		counter := measureSort(sorter, asc(data))
		if !isSorted(data) {
			testFail(name, &stat, counters, orig, data)
			return true
		}
		counters = append(counters, counter)
		stat.accumulate(counter)
	}
	testPass(name, &stat, counters)
	return true
}

func measureSort(sorter sorter, s slice) sortCounter {
	ic := newCounter()
	timer := newTimer()
	sorter.sort(ic, s)
	lapse := timer.stop()
	return sortCounter{
		inputs:  uint(s.Len()),
		compare: ic.inst("compare"),
		swap:    ic.inst("swap"),
		lapse:   lapse,
	}
}

type sortCounter struct {
	inputs  uint
	compare uint
	swap    uint
	lapse   time.Duration
}

type sortStat struct {
	compare   uint
	swap      uint
	lapse     time.Duration
	iteration uint
	inputs    int
}

func (s *sortStat) accumulate(c sortCounter) {
	s.compare += c.compare
	s.swap += c.swap
	s.lapse += c.lapse
	s.iteration++
}

func testFail(name string, stat *sortStat, counters []sortCounter, s, f []int) {
	showStat(os.Stderr, "Fail", name, stat, s, f)
}

func testPass(name string, stat *sortStat, counters []sortCounter) {
	showStat(os.Stdout, "OK", name, stat, nil, nil)
}

func mean(tot, n uint) float64 {
	if n == 0 {
		return 0.0
	}
	return float64(tot) / float64(n)
}

func meanDuration(d time.Duration, n uint) time.Duration {
	if n == 0 {
		return 0
	}
	return d / time.Duration(n)
}

func showStat(w io.Writer, verb, name string, stat *sortStat, orig, got []int) {
	fmt.Fprintf(w, "%s %s, compare = %.2f, swap = %.2f, time = %s",
		verb,
		name,
		mean(stat.compare, stat.iteration),
		mean(stat.swap, stat.iteration),
		meanDuration(stat.lapse, stat.iteration).String(),
	)
	if orig != nil {
		fmt.Fprintf(w, ", original = %#v", orig)
	}
	if got != nil {
		fmt.Fprintf(w, ", got = %#v", got)
	}
	fmt.Fprintln(w)
}

func convNanoseconds(d time.Duration) int64  { return int64(d) }
func convMicroseconds(d time.Duration) int64 { return int64(d) / 1e3 }
func convMilliseconds(d time.Duration) int64 { return int64(d) / 1e6 }

func isSorted(s []int) bool {
	if len(s) == 0 {
		return true
	}
	for i := range s[1:] {
		if s[i] > s[i+1] {
			return false
		}
	}
	return true
}

func exportCounters(path string, counters []sortCounter) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(f, 8, 0, 1, ' ', 0)
	fmt.Fprintln(w, "input size\tcompare\tswap\tlapse")
	for _, c := range counters {
		fmt.Fprintf(w, "%d\t%d\t%d\t%d\n",
			c.inputs,
			c.compare,
			c.swap,
			convNanoseconds(c.lapse),
		)
	}
	err = w.Flush()
	if err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func constInput(s []int) inputFunc {
	t := make([]int, len(s))
	return func(iteration uint) []int {
		copy(t, s)
		return t
	}
}

func fuzzInput(size int) inputFunc {
	s := make([]int, size)
	return func(iteration uint) []int {
		for i := range s {
			s[i] = rngSource.Intn(int(size))
		}
		return s
	}
}

func sortedInput(size int) inputFunc {
	s := make([]int, size)
	return func(iteration uint) []int {
		for i := range s {
			s[i] = i
		}
		return s
	}
}

func reversedInput(size int) inputFunc {
	s := make([]int, size)
	return func(iteration uint) []int {
		n := size - 1
		for i := range s {
			s[i] = n
			n--
		}
		return s
	}
}

func almostSortedInput(size, swap int) inputFunc {
	s := make([]int, size)
	return func(iteration uint) []int {
		for i := range s {
			s[i] = i
		}
		for i := 0; i < swap; i++ {
			i, j := rngSource.Intn(size), rngSource.Intn(size)
			s[i], s[j] = s[j], s[i]
		}
		return s
	}
}

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

func (a *aqsort) aqsort(ic *instCounter, sorter sorter, n int) []int {
	a.gas = n - 1
	a.nsolid = 0
	a.candidate = 0
	a.poison = make([]int, n)
	a.ptr = make([]int, n)
	for i := range a.poison {
		a.poison[i] = a.gas
		a.ptr[i] = i
	}
	sorter.sort(ic, a)
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

func (a *aqsort) freeze(i int) {
	a.poison[i] = a.nsolid
	a.nsolid++
}

func antiqsort(ic *instCounter, sorter sorter, n int) []int {
	var a aqsort
	return a.aqsort(ic, sorter, n)

}

func processTime() (time.Duration, error) {
	process, err := syscall.GetCurrentProcess()
	if err != nil {
		return 0, err
	}
	var r syscall.Rusage
	err = syscall.GetProcessTimes(process,
		&r.CreationTime,
		&r.ExitTime,
		&r.KernelTime,
		&r.UserTime,
	)
	if err != nil {
		return 0, err
	}
	convInt64 := func(time syscall.Filetime) int64 {
		return int64(time.HighDateTime)<<32 + int64(time.LowDateTime)
	}
	nsec := convInt64(r.KernelTime) + convInt64(r.UserTime)
	nsec *= 100
	return time.Duration(nsec), nil
}

func mustProcessTime() time.Duration {
	t, err := processTime()
	if err != nil {
		panic(err)
	}
	return t
}

type timer struct {
	start time.Duration
}

func newTimer() *timer {
	return &timer{start: mustProcessTime()}
}

func (t *timer) stop() time.Duration {
	now := mustProcessTime()
	return now - t.start
}

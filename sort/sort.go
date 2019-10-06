//+build mage

package main

import (
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"golang.org/x/exp/rand"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

var rngSeed = time.Now().UnixNano()

var rngSource = rand.New(rand.NewSource(uint64(rngSeed)))

// list는 배열을 나타냅니다. partial order개 정의되어 있어야 합니다.
type list interface {
	Len() int
	Swap(i, j int)
	Less(i, j int) bool
	Peek(i int) int
}

// sorter는 정렬 알고리즘을 나타냅니다.
type sorter interface {
	name() string
	sort(list list)
}

type aslist []int

func (s aslist) Len() int           { return len(s) }
func (s aslist) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s aslist) Less(i, j int) bool { return s[i] <= s[j] }
func (s aslist) Peek(i int) int     { return s[i] }

type measurelist struct {
	list  list
	nlen  int
	nswap int
	nless int
	npeek int
}

func (l *measurelist) Len() int {
	l.nlen++
	return l.list.Len()
}

func (l *measurelist) Swap(i, j int) {
	l.nswap++
	l.list.Swap(i, j)
}

func (l *measurelist) Less(i, j int) bool {
	l.nless++
	return l.list.Less(i, j)
}

func (l *measurelist) Peek(i int) int {
	l.npeek++
	return l.list.Peek(i)
}

type label []string

func (buf label) append(s string) label {
	return append(buf, s)
}

func (buf *label) String() string {
	return strings.Join([]string(*buf), " - ")
}

type benchmarkResult struct {
	records []benchmarkRecord
}

type benchmarkRecord struct {
	sortName      string
	inputName     string // TODO: extends
	sizedCounters [][]sortCounter
}

func benchmark(sorter sorter, maxsize, iteration uint) benchmarkResult {
	swappedSorted := func(swapFactor float64) sizedInputFunc {
		return func(size uint) inputFunc {
			swap := uint(float64(size) * swapFactor)
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
	if _, ok := sorter.(isqsort); ok {
		killqsort := func(size uint) inputFunc {
			killer := antiqsort(sorter, int(size))
			return constInput(killer)
		}
		inputs = append(inputs, inputType{"killing input", killqsort})
	}
	records := make([]benchmarkRecord, 0, len(inputs))
	for _, input := range inputs {
		record := benchmarkInput(sorter, input.name, input.makeinput, iteration, maxsize)
		records = append(records, record)
	}
	return benchmarkResult{records: records}
}

func benchmarkInput(sorter sorter, inputName string, makeinput sizedInputFunc, iteration, maxsize uint) benchmarkRecord {
	sizedCounters := iterateSizedSort(sorter, makeinput, constIteration(iteration), maxsize)
	return benchmarkRecord{
		sortName:      sorter.name(),
		inputName:     inputName,
		sizedCounters: sizedCounters,
	}
}

type inputFunc func(iteration uint) []int

type testResult struct {
	pass      bool
	counters  []sortCounter
	errinput  []int
	erroutput []int
}

func testSort(sorter sorter) testResult {
	const maxiteration = 200
	result := testResult{
		counters: make([]sortCounter, 0, maxiteration),
	}
	var rawinput = make([]int, 0, maxiteration)
	for i := uint(0); i <= maxiteration; i++ {
		size := i
		input := fuzzInput(size)(i)
		rawinput = append(rawinput[:0], input...)
		counter := measureSort(sorter, aslist(input))
		result.counters = append(result.counters, counter)
		if !isSorted(input) {
			result.errinput = rawinput
			result.erroutput = input
			return result
		}
	}
	result.pass = true
	return result
}

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

type sortCounter struct {
	nlen  int
	nswap int
	nless int
	npeek int
	lapse time.Duration
}

func measureSort(sorter sorter, list list) sortCounter {
	tool := measurelist{list: list}
	startTime := time.Now()
	sorter.sort(&tool)
	lapse := time.Since(startTime)
	return sortCounter{
		nlen:  tool.nlen,
		nswap: tool.nswap,
		nless: tool.nless,
		npeek: tool.npeek,
		lapse: lapse,
	}
}

func iterateSort(sorter sorter, inputFunc inputFunc, iteration uint) []sortCounter {
	// Design: 함수 구조로 만들어져 있어 iterateSort와 iterateSizedSort는
	// 병렬화가 가능해 보인다.
	counters := make([]sortCounter, 0, iteration)
	for i := uint(1); i <= iteration; i++ {
		input := inputFunc(i)
		counter := measureSort(sorter, aslist(input))
		counters = append(counters, counter)
	}
	return counters
}

type sizedInputFunc func(size uint) inputFunc

type sizedIterationFunc func(size uint) (iteration uint)

func iterateSizedSort(sorter sorter, sizedInputFunc sizedInputFunc, sizedIterationFunc sizedIterationFunc, maxsize uint) [][]sortCounter {
	buf := make([][]sortCounter, 0, maxsize+1)
	for size := uint(0); size <= maxsize; size++ {
		inputFunc := sizedInputFunc(size)
		iteration := sizedIterationFunc(size)
		buf = append(buf, iterateSort(sorter, inputFunc, iteration))
	}
	return buf
}

type sortStat struct {
	averageLen   float64
	averageSwap  float64
	averageLess  float64
	averagePeek  float64
	averageLapse time.Duration
	iteration    uint
}

func averageCounters(counters []sortCounter) sortStat {
	var stat sortStat
	for _, counter := range counters {
		stat.averageLen += float64(counter.nlen)
		stat.averageSwap += float64(counter.nswap)
		stat.averageLess += float64(counter.nless)
		stat.averagePeek += float64(counter.npeek)
		stat.averageLapse += counter.lapse
	}
	n := len(counters)
	if n > 0 {
		stat.averageLen /= float64(n)
		stat.averageSwap /= float64(n)
		stat.averageLess /= float64(n)
		stat.averagePeek /= float64(n)
		stat.averageLapse /= time.Duration(n)
	}
	stat.iteration = uint(n)
	return stat
}

func constInput(s []int) inputFunc {
	t := make([]int, len(s))
	return func(iteration uint) []int {
		copy(t, s)
		return t
	}
}

func fuzzInput(size uint) inputFunc {
	s := make([]int, size)
	return func(iteration uint) []int {
		for i := range s {
			s[i] = rngSource.Intn(int(size))
		}
		return s
	}
}

func sortedInput(size uint) inputFunc {
	s := make([]int, size)
	return func(iteration uint) []int {
		for i := range s {
			s[i] = i
		}
		return s
	}
}

func reversedInput(size uint) inputFunc {
	s := make([]int, size)
	return func(iteration uint) []int {
		n := int(size - 1)
		for i := range s {
			s[i] = n
			n--
		}
		return s
	}
}

func almostSortedInput(size, swap uint) inputFunc {
	s := make([]int, size)
	return func(iteration uint) []int {
		for i := range s {
			s[i] = i
		}
		for i := 0; i < int(swap); i++ {
			i, j := rngSource.Intn(int(size)), rngSource.Intn(int(size))
			s[i], s[j] = s[j], s[i]
		}
		return s
	}
}

func constIteration(iteration uint) sizedIterationFunc {
	return func(size uint) uint {
		return iteration
	}
}

type drawingData struct {
	sortName  string
	inputName string
	samples   []sortStat
}

func newDrawingData(record benchmarkRecord) *drawingData {
	return &drawingData{
		sortName:  record.sortName,
		inputName: record.inputName,
		samples:   newSamples(record.sizedCounters),
	}
}

func newSamples(sizedCounters [][]sortCounter) []sortStat {
	samples := make([]sortStat, 0, len(sizedCounters))
	for _, counters := range sizedCounters {
		samples = append(samples, averageCounters(counters))
	}
	return samples
}

type drawing struct {
	imageWidth  vg.Length
	imageHeight vg.Length
}

func documentDrawing() *drawing {
	const imageWidth = 13 * vg.Centimeter
	const imageHeight = 4.63 * vg.Centimeter
	return &drawing{
		imageWidth:  imageWidth,
		imageHeight: imageHeight,
	}
}

func largeDrawing() *drawing {
	const imageWidth = 30 * vg.Centimeter
	const imageHeight = 25 * vg.Centimeter
	return &drawing{
		imageWidth:  imageWidth,
		imageHeight: imageHeight,
	}
}

func (w *drawing) drawSamples(data *drawingData, server serveY) (*plot.Plot, error) {
	pl, err := plot.New()
	if err != nil {
		return nil, err
	}

	// TODO: size scaling
	set := serveSet(server, data.samples)
	line, _, err := plotter.NewLinePoints(set)
	if err != nil {
		return nil, err
	}
	line.Width = 0.12 * vg.Centimeter
	line.Color = plotutil.Color(0)

	pl.Add(line)
	pl.Legend.Add(server.label(), line)
	pl.Title.Text = fmt.Sprintf("%s - %s", data.sortName, data.inputName)
	pl.Title.Font.Size = w.goodSize(pl.Title.Font, 0.6*vg.Centimeter, pl.Title.Text)

	pl.X.Label.Text = "Size"
	pl.X.Label.Font.Size = 0.6 * vg.Centimeter
	pl.Y.Label.Text = server.label()
	pl.Y.Label.Font.Size = 0.6 * vg.Centimeter

	n := numeration{i: 1}
	drawFunc(pl, identifyFunc, "y = x", n.inc())
	drawFunc(pl, quadraticFunc, "y = x²", n.inc())
	drawFunc(pl, xlogxFunc, "y = xlog(x)", n.inc())

	return pl, nil
}

func (w *drawing) goodSize(font vg.Font, initial vg.Length, text string) vg.Length {
	font.Size = initial
	for {
		size := font.Width(text)
		if size <= w.imageWidth {
			return font.Size
		}
		font.Size -= 0.5
	}
}

type numeration struct {
	i int
}

func (n *numeration) inc() int {
	const numerationModular = 10
	i := n.i
	n.i++
	return i % numerationModular
}

func drawFunc(pl *plot.Plot, fn func(float64) float64, legend string, style int) {
	function := plotter.NewFunction(fn)
	function.Color = plotutil.Color(style)
	function.Width = 0.1 * vg.Centimeter
	pl.Add(function)
	pl.Legend.Add(legend, function)
}

type serveY interface {
	serveY(sample sortStat) float64
	label() string
}

func serveSet(server serveY, samples []sortStat) plotter.XYs {
	set := make(plotter.XYs, len(samples))
	for i, sample := range samples {
		set[i].X = float64(i)
		set[i].Y = server.serveY(sample)
	}
	return set
}

type serveCompare struct{}

func (serveCompare) label() string {
	return "compare"
}

func (serveCompare) serveY(sample sortStat) float64 {
	return sample.averageLess
}

type serveSwap struct{}

func (serveSwap) label() string {
	return "swap"
}

func (serveSwap) serveY(sample sortStat) float64 {
	return sample.averageSwap
}

type serveAccess struct{}

func (serveAccess) label() string {
	return "access"
}

func (serveAccess) serveY(sample sortStat) float64 {
	return sample.averagePeek
}

type serveMicrosecondLapse struct{}

func (serveMicrosecondLapse) label() string {
	return "time"
}

func (serveMicrosecondLapse) serveY(sample sortStat) float64 {
	return float64(convMicroseconds(sample.averageLapse))
}

func convNanoseconds(d time.Duration) int64  { return int64(d) }
func convMicroseconds(d time.Duration) int64 { return int64(d) / 1e3 }
func convMilliseconds(d time.Duration) int64 { return int64(d) / 1e6 }

// func plotSave(label label, pl *plot.Plot) error {
// 	dir := label[0]
// 	os.Mkdir(dir, os.ModePerm)
// 	name := fmt.Sprintf("%s.jpg", label.String())
// 	path := filepath.Join(dir, name)
// 	return pl.Save(imageWidth, imageHeight, path)
// }

func identifyFunc(x float64) float64 { return x }

func quadraticFunc(x float64) float64 { return x * x }

func xlogxFunc(x float64) float64 {
	if x <= 1 {
		return 0
	}
	return x * math.Log2(x)
}

func runTest(sorter sorter) bool {
	r := testSort(sorter)
	if !r.pass {
		testFail(sorter, showcase{
			stat:      averageCounters(r.counters),
			errinput:  r.errinput,
			erroutput: r.erroutput,
		})
		return false
	}
	testPass(sorter, showcase{
		stat: averageCounters(r.counters),
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
		sorter.name(),
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

type command struct {
	fn interface{}
}

func (c command) Call() error {
	ctx := context.Background()
	return c.CallContext(ctx)
}

func (c command) CallContext(ctx context.Context) error {
	switch fn := c.fn.(type) {
	case func():
		fn()
		return nil
	case func() error:
		return fn()
	case func(context.Context):
		fn(ctx)
		return nil
	case func(context.Context) error:
		return fn(ctx)
	default:
		panic(fmt.Sprintf("invalid command of type %T", fn))
	}
}

func namespacedCommands(namespace interface{}) []command {
	r := reflect.ValueOf(namespace)
	for r.Kind() == reflect.Ptr || r.Kind() == reflect.Interface {
		r = r.Elem()
	}
	if r.NumMethod() == 0 {
		return nil
	}
	cmds := make([]command, 0, r.NumMethod())
	for i := 0; i < r.NumMethod(); i++ {
		fn := r.Method(i).Interface()
		if isCommand(fn) {
			cmds = append(cmds, command{fn: fn})
		}
	}
	return cmds
}

func isCommand(fn interface{}) bool {
	switch fn.(type) {
	case func():
		return true
	case func() error:
		return true
	case func(context.Context):
		return true
	case func(context.Context) error:
		return true
	}
	return false
}

func runAll(ctx context.Context, namespace interface{}) error {
	for _, cmd := range namespacedCommands(namespace) {
		err := cmd.CallContext(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func Tests(ctx context.Context) error {
	return runAll(ctx, Test{})
}

type Test mg.Namespace

func (Test) ssort() {
	runTest(ssort{})
}

func (Test) Bsort() {
	runTest(bsort{})
}

func (Test) Hsort() {
	runTest(hsort{})
}

// TODO: Aliases 사용해서 test-x 같은 방식 구상해보기

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

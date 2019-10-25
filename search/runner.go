package main

import (
	"algorithm-class/inst"
	"time"
)

type Runner interface {
	Name() string
	Do(*Timer, []int) inst.State
	// Insts() []inst.Kind
}

type perf struct {
	state  inst.State
	elapse time.Duration
}

type perfStat struct {
	state     inst.State
	elapse    time.Duration
	iteration uint
}

func (s *perfStat) averageInst(kind inst.Kind) float64 {
	if s.iteration == 0 {
		return 0.0
	}
	n := s.state.Get(kind)
	return float64(n) / float64(s.iteration)
}

func (s *perfStat) averageElapse() time.Duration {
	if s.iteration == 0 {
		return time.Duration(0)
	}
	return s.elapse / time.Duration(s.iteration)
}

func (s perfStat) addperf(perf perf) perfStat {
	for kind, n := range perf.state {
		s.state[kind] += n
	}
	s.elapse += perf.elapse
	return s
}

type perfClass struct {
	runner     Runner
	inputClass string
	stats      []perfStat
}

func (c perfClass) Size() int {
	return len(c.stats)
}

type sizedInput interface {
	class() string

	input(size int) []int
}

func timeitClass(runner Runner, size int, sizedInput sizedInput, iteration uint) perfClass {
	stats := make([]perfStat, 0, size)
	for i := 0; i <= size; i++ {
		var stat perfStat
		for it := uint(0); it <= iteration; it++ {
			perf := timeit(runner, sizedInput.input(i))
			stat = stat.addperf(perf)
		}
		stat.iteration = iteration
		stats = append(stats, stat)
	}
	return perfClass{
		runner:     runner,
		inputClass: sizedInput.class(),
		stats:      stats,
	}
}

func timeit(runner Runner, input []int) perf {
	timer := newTimer(time.Now())
	state := runner.Do(timer, input)
	timer.Stop()
	return perf{
		state:  state,
		elapse: timer.elapse,
	}
}

func timeitAll(runner Runner, size int, iteration uint) []perfClass {
	buf := make([]int, 0, size)
	var inputClasses = []sizedInput{
		fuzzInput{buf: buf},
		sortedInput{buf: buf},
		reversedInput{buf: buf},
		swappedSortedInput{buf: buf, ratio: 0.10},
	}
	perfs := make([]perfClass, 0, len(inputClasses))
	for _, inputClass := range inputClasses {
		perfcls := timeitClass(runner, size, inputClass, iteration)
		perfs = append(perfs, perfcls)
	}
	return perfs
}

type Timer struct {
	before time.Time
	elapse time.Duration
}

func newTimer(before time.Time) *Timer {
	return &Timer{before: before}
}

func (t *Timer) Stop() {
	if t.before.IsZero() {
		return
	}
	t.elapse += time.Since(t.before)
	t.before = time.Time{}
}

func (t *Timer) Start() {
	t.before = time.Now()
}

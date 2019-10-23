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

type perfClass struct {
	runner     Runner
	inputClass string
	perfs      []perf
}

func (c perfClass) Size() int {
	return len(c.perfs)
}

type sizedInput interface {
	class() string

	input(size int) []int
}

func timeitClass(runner Runner, size int, sizedInput sizedInput) perfClass {
	perfs := make([]perf, 0, size)
	for i := 0; i <= size; i++ {
		perf := timeit(runner, sizedInput.input(size))
		perfs = append(perfs, perf)
	}
	return perfClass{
		runner:     runner,
		inputClass: sizedInput.class(),
		perfs:      perfs,
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

func timeitAll(runner Runner, size int) []perfClass {
	buf := make([]int, 0, size)
	var inputClasses = []sizedInput{
		fuzzInput{buf: buf},
		sortedInput{buf: buf},
		reversedInput{buf: buf},
		swappedSortedInput{buf: buf, ratio: 0.10},
	}
	perfs := make([]perfClass, 0, len(inputClasses))
	for _, inputClass := range inputClasses {
		perfcls := timeitClass(runner, size, inputClass)
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

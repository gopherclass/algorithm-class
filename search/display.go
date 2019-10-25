package main

import (
	"algorithm-class/inst"
	avl "algorithm-class/search/AVL-tree"
	binaryTree "algorithm-class/search/binary-tree"
	rb "algorithm-class/search/red-black-tree"
	"time"
)

type redblacktree struct{}

func (redblacktree) Name() string {
	return "red-black tree"
}

func (redblacktree) Do(timer *Timer, input []int) inst.State {
	timer.Stop()
	tree := rb.NewTree()
	for _, x := range input {
		tree.Insert(nil, x)
	}
	ic := inst.NewCounter()
	timer.Start()
	for _, x := range input {
		tree.Search(ic, x)
	}
	return ic.State()
}

func (redblacktree) Illusts() []Illust {
	return []Illust{
		DisplayIndirect{},
		DisplayTime{},
	}
}

type avltree struct{}

func (avltree) Name() string {
	return "AVL tree"
}

func (avltree) Do(timer *Timer, input []int) inst.State {
	timer.Stop()
	tree := avl.NewTree()
	for _, x := range input {
		tree.Insert(nil, x)
	}
	ic := inst.NewCounter()
	timer.Start()
	for _, x := range input {
		tree.Search(ic, x)
	}
	return ic.State()
}

func (avltree) Illusts() []Illust {
	return []Illust{
		DisplayIndirect{},
		DisplayTime{},
	}
}

type binarytree struct{}

func (binarytree) Name() string {
	return "Binary tree"
}

func (binarytree) Do(timer *Timer, input []int) inst.State {
	timer.Stop()
	tree := binaryTree.NewTree()
	for _, x := range input {
		tree.Insert(nil, x)
	}
	ic := inst.NewCounter()
	timer.Start()
	for _, x := range input {
		tree.Search(ic, x)
	}
	return ic.State()
}

func (binarytree) Illusts() []Illust {
	return []Illust{
		DisplayIndirect{},
		DisplayTime{},
	}
}

type DisplayIndirect struct{}

func (DisplayIndirect) Fx(perf perfStat) float64 {
	return float64(perf.averageInst(inst.Indirect))
}

func (DisplayIndirect) Legend(cls perfClass) string {
	return cls.runner.Name()
}

func (DisplayIndirect) Tag() string {
	return "Indirect"
}

type DisplayTime struct{}

func (DisplayTime) Fx(perf perfStat) float64 {
	return float64(convMicroseconds(perf.averageElapse()))
}

func convNanoseconds(d time.Duration) int64  { return int64(d) }
func convMicroseconds(d time.Duration) int64 { return int64(d) / 1e3 }
func convMilliseconds(d time.Duration) int64 { return int64(d) / 1e6 }

func (DisplayTime) Legend(cls perfClass) string {
	return cls.runner.Name()
}

func (DisplayTime) Tag() string {
	return "Time"
}

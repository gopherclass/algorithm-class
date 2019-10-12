package main

import "time"

type sortCounter struct {
	nlen   int
	nswap  int
	nless  int
	npeek  int
	nset   int
	nslice int
	npush  int
	npop   int
	npos   int
	nnext  int
	nprev  int
	lapse  time.Duration
}

type sortStat struct {
	averageLen   float64
	averageSwap  float64
	averageLess  float64
	averagePeek  float64
	averageSet   float64
	averageSlice float64
	averagePush  float64
	averagePop   float64
	averagePos   float64
	averageNext  float64
	averagePrev  float64
	averageLapse time.Duration
	iteration    uint
}

func accCounter(stat *sortStat, c sortCounter) {
	stat.averageLen += float64(c.nlen)
	stat.averageSwap += float64(c.nswap)
	stat.averageLess += float64(c.nless)
	stat.averagePeek += float64(c.npeek)
	stat.averageSet += float64(c.nset)
	stat.averageSlice += float64(c.nslice)
	stat.averagePush += float64(c.npush)
	stat.averagePop += float64(c.npop)
	stat.averagePos += float64(c.npos)
	stat.averageNext += float64(c.nnext)
	stat.averagePrev += float64(c.nprev)
	stat.averageLapse += c.lapse
	stat.iteration++
}

func averageStat(stat *sortStat) {
	if stat.iteration == 0 {
		return
	}
	n := float64(stat.iteration)
	stat.averageLen /= n
	stat.averageSwap /= n
	stat.averageLess /= n
	stat.averagePeek /= n
	stat.averageSet /= n
	stat.averageSlice /= n
	stat.averagePush /= n
	stat.averagePop /= n
	stat.averagePos /= n
	stat.averageNext /= n
	stat.averagePrev /= n
	stat.averageLapse /= time.Duration(n)
}

func (c *sortCounter) Len() bool {
	if c == nil {
		return true
	}
	c.nlen++
	return true
}

func (c *sortCounter) Swap() bool {
	if c == nil {
		return true
	}
	c.nswap++
	return true
}

func (c *sortCounter) Less() bool {
	if c == nil {
		return true
	}
	c.nless++
	return true
}

func (c *sortCounter) Set() bool {
	if c == nil {
		return true
	}
	c.nset++
	return true
}

func (c *sortCounter) Peek() bool {
	if c == nil {
		return true
	}
	c.npeek++
	return true
}

func (c *sortCounter) Slice() bool {
	if c == nil {
		return true
	}
	c.nslice++
	return true
}

func (c *sortCounter) Push() bool {
	if c == nil {
		return true
	}
	c.npush++
	return true
}

func (c *sortCounter) Pop() bool {
	if c == nil {
		return true
	}
	c.npop++
	return true
}

func (c *sortCounter) Pos() bool {
	c.npos++
	return true
}

func (c *sortCounter) Next() bool {
	if c == nil {
		return true
	}
	c.nnext++
	return true
}

func (c *sortCounter) Prev() bool {
	if c == nil {
		return true
	}
	c.nprev++
	return true
}

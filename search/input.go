package main

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

var rngSource = rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

func effectiveBits(n int) int {
	return n & ^(^0 << treeBits)
}

type fuzzInput struct {
	buf []int
}

func (fuzzInput) class() string {
	return "fuzz input"
}

func (v fuzzInput) input(n int) []int {
	s := v.buf[:n]
	for i := range s {
		s[i] = effectiveBits(rngSource.Int())
	}
	return s
}

type sortedInput struct {
	buf []int
}

func (sortedInput) class() string { return "sorted input" }

func (v sortedInput) input(n int) []int {
	s := v.buf[:n]
	for i := range s {
		s[i] = i
	}
	return s
}

type reversedInput struct {
	buf []int
}

func (reversedInput) class() string { return "reversed input" }

func (v reversedInput) input(n int) []int {
	s := v.buf[:n]
	for i := range s {
		s[i] = len(s) - i - 1
	}
	return s
}

type swappedSortedInput struct {
	buf   []int
	ratio float64
}

func (v swappedSortedInput) class() string {
	return fmt.Sprintf("almost sorted input (%.1f swapped)", v.ratio)
}

func (v swappedSortedInput) input(n int) []int {
	s := v.buf[:n]
	for i := range s {
		s[i] = i
	}
	nswap := int(float64(n) * v.ratio)
	for i := 0; i < nswap; i++ {
		i, j := rngSource.Intn(len(s)), rngSource.Intn(len(s))
		s[i], s[j] = s[j], s[i]
	}
	return s
}

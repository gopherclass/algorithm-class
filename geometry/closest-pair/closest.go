package closestPair

import (
	"algorithm-class/geometry/plane"
	"math"
	"sort"

	"github.com/bradfitz/slice"
)

func Closest(s []plane.Point) []plane.Point {
	slice.Sort(s, func(i, j int) bool {
		return s[i].X <= s[j].X
	})
	return closest(s, 0, len(s))
}

func closest(s []plane.Point, i, j int) []plane.Point {
	if j-i < 2 {
		return nil
	}
	if j-i == 2 {
		return s[i:j]
	}
	m := (i + j) / 2
	r0 := closest(s, i, m)
	r1 := closest(s, m, j)
	d0 := dist(r0)
	d1 := dist(r1)
	d := int(math.Ceil(d0))
	if d1 < d0 {
		d = int(math.Ceil(d1))
	}
	left := sort.Search(len(s), func(i int) bool {
		return s[i].X >= s[m].X-d
	})
	if left < i {
		left = i
	}
	right := sort.Search(len(s), func(i int) bool {
		return s[i].X > s[m].X+d
	})
	if right > j {
		right = j
	}
	r2 := closestNaive(s[left:right])
	if d2 := dist(r2); d2 < min(d0, d1) {
		return r2
	}
	if d0 < d1 {
		return r0
	}
	return r1
}

func min(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

func dist(s []plane.Point) float64 {
	if len(s) != 2 {
		return math.Inf(1)
	}
	return dist2(s[0], s[1])
}

func dist2(a, b plane.Point) float64 {
	p := b.Sub(a)
	return math.Hypot(float64(p.X), float64(p.Y))
}

func closestNaive(s []plane.Point) []plane.Point {
	if len(s) < 2 {
		return nil
	}
	mi, mj, md := 0, 0, math.Inf(1)
	for i, a := range s {
		for j, b := range s {
			if i == j {
				continue
			}
			d := dist2(a, b)
			if d >= md {
				continue
			}
			mi, mj, md = i, j, d
		}
	}
	return []plane.Point{s[mi], s[mj]}
}

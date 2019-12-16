package packageWrapping

import (
	"algorithm-class/geometry/plane"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func PseudoArg(p plane.Point) float64 {
	if p.X == 0 && p.Y == 0 {
		return 0
	}
	t := float64(p.Y) / float64(abs(p.X)+abs(p.Y))
	if p.X < 0 {
		return 2 - t
	}
	if p.Y < 0 {
		return 4 + t
	}
	return t
}

func Theta(a, b plane.Point) float64 {
	return PseudoArg(b.Sub(a))
}

func Wrap(s []plane.Point) []plane.Point {
	var maxarg float64
	next := chiefPoint(s)
	for pivot := range s {
		s[pivot], s[next] = s[next], s[pivot]
		next = 0
		minarg := maxarg
		maxarg = 4.0
		comparePivot := func(i int) {
			arg := Theta(s[pivot], s[i])
			if minarg < arg && arg < maxarg {
				maxarg = arg
				next = i
			}
		}
		for i := pivot + 1; i < len(s); i++ {
			comparePivot(i)
		}
		comparePivot(0)
		if next <= 0 {
			return s[:pivot+1]
		}
	}
	return s
}

func chiefPoint(s []plane.Point) int {
	if len(s) == 0 {
		return -1
	}
	i := 0
	for v := range s {
		if s[v].Y < s[i].Y {
			i = v
		}
	}
	return i
}

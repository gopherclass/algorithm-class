package grahamScan

import (
	packageWrapping "algorithm-class/geometry/package-wrapping"
	"algorithm-class/geometry/plane"
	"algorithm-class/geometry/rotation"

	"github.com/bradfitz/slice"
)

func Scan(s []plane.Point) []plane.Point {
	i := chiefPoint(s)
	pivot := s[i]
	s[0], s[i] = s[i], s[0]
	slice.Sort(s, func(i, j int) bool {
		return packageWrapping.Theta(pivot, s[i]) <=
			packageWrapping.Theta(pivot, s[j])
	})
	top := 1
	for i := 2; i < len(s); i++ {
		for top >= 1 && rotation.Dir(s[top-1], s[top], s[i]) < 0 {
			top--
		}
		top++
		s[i], s[top] = s[top], s[i]
	}
	return s[:top+1]
}

func chiefPoint(s []plane.Point) int {
	if len(s) == 0 {
		return -1
	}
	i := 0
	for v := range s {
		if s[v].Y < s[i].Y || (s[v].Y == s[i].Y && s[i].X < s[v].X) {
			i = v
		}
	}
	return i
}

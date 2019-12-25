package inPolygon

import (
	"algorithm-class/geometry/plane"
)

// In은 점 x가 다각형 s 내부에 있는지 검사합니다. 경계에 위치한 점에 대하여
// 일관되지 않은 검사 결과를 반환할 수 있습니다.
func In(s []plane.Point, x plane.Point) bool {
	i := -1
	xMin, xMax := s[0].X, s[0].X
	for j := range s {
		if s[j].Y != x.Y {
			i = j
			break
		}
	}
	if i < 0 {
		return xMin <= x.X && x.X <= xMax
	}
	k := i

	crossings := 0
	check := func(j int) {
		var e, ray plane.Segment
		e.Min = s[i]
		e.Max = s[j]
		e = e.Canon()
		if e.Max.X < x.X {
			i = j
			return
		}
		ray.Min = x
		ray.Max.X = e.Max.X + 1
		ray.Max.Y = x.Y
		if (x.Y == e.Min.Y && x.X <= e.Min.X) || (x.Y == e.Max.Y && x.X <= e.Max.X) {
			return
		}
		i = j
		if ray.Intersect(e) {
			crossings++
		}
	}

	for j := k + 1; j < len(s); j++ {
		check(j)
	}
	for j := 0; j <= k; j++ {
		check(j)
	}
	return crossings%2 == 1
}

func OnBoundary(s []plane.Point, x plane.Point) bool {
	panic("not yet implemented")
}

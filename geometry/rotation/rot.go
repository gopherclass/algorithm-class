package rotation

import "algorithm-class/geometry/plane"

// Dir는 한 점 x가 x -> y -> z 순서로 이동할 때 어떤 방향으로 움직이는지
// 알려줍니다
//   Dir > 0 <=> 반시계방향
//   Dir = 0 <=> 방향 없음
//   Dir < 0 <=> 시계방향
//
func Dir(x, y, z plane.Point) int {
	return plane.Det(y.Sub(x), z.Sub(x))
}

func Sign(x int) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

func Intersect(r0, r1 plane.Segment) bool {
	return Dir(r0.Min, r0.Max, r1.Min)*Dir(r0.Min, r0.Max, r1.Max) <= 0 &&
		Dir(r1.Min, r1.Max, r0.Min)*Dir(r1.Min, r1.Max, r0.Max) <= 0
}

func StrictlyIntersect(r0, r1 plane.Segment) bool {
	panic("not yet implemented")
}

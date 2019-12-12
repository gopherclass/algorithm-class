package plane

// Segment는 양 끝 두점으로 선분을 표현합니다.
type Segment struct {
	// Min은 선분의 두 점 중 x 값이 더 작은 점입니다.
	Min Point
	// Max은 선분의 두 점 주 x 값이 더 큰 점입니다.
	Max Point
}

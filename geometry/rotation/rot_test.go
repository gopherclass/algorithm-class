package rotation

import (
	"algorithm-class/geometry/plane"
	"algorithm-class/random"
	"testing"

	"github.com/stretchr/testify/assert"
)

var rng = random.Now()

func TestDir(t *testing.T) {
	testcases := []struct {
		x, y, z  plane.Point
		expected int
	}{
		{plane.Point{0, 0}, plane.Point{0, 0}, plane.Point{0, 0}, 0},
		{plane.Point{0, 0}, plane.Point{1, 1}, plane.Point{0, 2}, 1},
		// TODO: 더 많은 테스트케이스가 필요하다.
	}
	for _, tc := range testcases {
		got := Dir(tc.x, tc.y, tc.z)
		assert.Equal(t, Sign(tc.expected), Sign(got), tc.x, tc.y, tc.z)
	}
}

func TestDirAgainst(t *testing.T) {
	for i := 0; i < 1000; i++ {
		x, y, z := rngPoint(), rngPoint(), rngPoint()
		expected := Dir(x, y, z)
		got := ccw(x, y, z)
		assert.Equal(t, Sign(expected), Sign(got), x, y, z)
	}
}

func rngPoint() plane.Point {
	return plane.Point{
		X: rng.Intn(1000),
		Y: rng.Intn(1000),
	}
}

func ccw(p0, p1, p2 plane.Point) int {
	p1 = p1.Sub(p0)
	p2 = p2.Sub(p0)
	if p1.X*p2.Y > p1.Y*p2.X {
		return 1
	}
	if p1.X*p2.Y < p1.Y*p2.X {
		return -1
	}
	if p1.X == 0 && p1.Y == 0 {
		return 0
	}
	if p1.X*p2.X < 0 || p1.Y*p2.Y < 0 {
		return -1
	}
	if p1.Dot(p1) < p2.Dot(p2) {
		return 1
	}
	return 0
}

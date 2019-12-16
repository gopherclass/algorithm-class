package packageWrapping

import (
	"algorithm-class/geometry/plane"
	"algorithm-class/random"
	"testing"

	"github.com/stretchr/testify/require"
)

var rng = random.Now()

func TestPseudoArg(t *testing.T) {
	require.Equal(t, 0.0, PseudoArg(plane.Point{0, 0}))
	for i := 0; i < 1000; i++ {
		x, y := rngXY()
		arg0 := PseudoArg(plane.Point{x, y})
		arg1 := PseudoArg(plane.Point{x*x - y*y, 2 * x * y})
		if y == 0 {
			require.Equal(t, 2.0, arg0, "x = %f, y = %f", x, y)
			require.Equal(t, 0.0, arg1, "x = %f, y = %f", x, y)
			continue
		}
		require.Less(t, arg0, arg1, "x = %f, y = %f", x, y)
	}
}

func rngXY() (int, int) {
	for {
		x, y := rng.Intn(2000)-1000, rng.Intn(1000)
		if x == 0 && y == 0 {
			continue
		}
		return x, y
	}
}

func TestWrap(t *testing.T) {
	t.Skip()
}

func TestExample(t *testing.T) {
	points, pointsMap := examplePoints()
	for _, p := range Wrap(points) {
		t.Log(p, pointsMap[p])
	}
}

func TestProblem1(t *testing.T) {
	points, pointsMap := problem53Points()
	for _, p := range Wrap(points) {
		t.Log(p, pointsMap[p])
	}
}

func examplePoints() ([]plane.Point, map[plane.Point]string) {
	points := []plane.Point{
		{4, 6},
		{12, 16},
		{3, 12},
		{10, 11},
		{14, 4},
		{1, 10},
		{13, 8},
		{6, 7},
		{8, 9},
		{7, 5},
		{15, 3},
		{16, 14},
		{2, 15},
		{11, 1},
		{9, 13},
		{5, 2},
	}
	chars := []string{
		"A",
		"B",
		"C",
		"D",
		"E",
		"F",
		"G",
		"H",
		"I",
		"J",
		"K",
		"L",
		"M",
		"N",
		"O",
		"P",
	}
	pointsMap := make(map[plane.Point]string)
	for i, p := range points {
		pointsMap[p] = chars[i]
	}
	return points, pointsMap
}

func problem53Points() ([]plane.Point, map[plane.Point]string) {
	points := []plane.Point{
		{5, 3},
		{7, 11},
		{6, 12},
		{8, 2},
		{1, 4},
		{4, 10},
		{3, 6},
		{2, 5},
		{11, 8},
		{12, 9},
		{10, 1},
		{9, 7},
	}
	chars := []string{
		"A",
		"B",
		"C",
		"D",
		"E",
		"F",
		"G",
		"H",
		"I",
		"J",
		"K",
		"L",
	}
	pointsMap := make(map[plane.Point]string)
	for i, p := range points {
		pointsMap[p] = chars[i]
	}
	return points, pointsMap
}

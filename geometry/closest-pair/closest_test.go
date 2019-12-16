package closestPair

import (
	"algorithm-class/geometry/plane"
	"testing"
)

func TestClosest(t *testing.T) {
	points, pointsMap := examplePoints()
	r := Closest(points)
	t.Log(r[0], pointsMap[r[0]])
	t.Log(r[1], pointsMap[r[1]])
	t.Log(dist(r))
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

func TestProblem3(t *testing.T) {
	points, pointsMap := examplePoints()
	r := Closest(points)
	t.Log(r[0], pointsMap[r[0]])
	t.Log(r[1], pointsMap[r[1]])
	t.Log(dist(r))
}

func problem54Points() ([]plane.Point, map[plane.Point]string) {
	points := []plane.Point{
		{5, 3},
		{7, 11},
		{3, 10},
		{1, 6},
		{9, 7},
		{12, 5},
		{10, 4},
		{2, 2},
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
	}
	pointsMap := make(map[plane.Point]string)
	for i, p := range points {
		pointsMap[p] = chars[i]
	}
	return points, pointsMap
}

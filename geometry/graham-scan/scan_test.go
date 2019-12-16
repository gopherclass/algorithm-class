package grahamScan

import (
	"algorithm-class/geometry/plane"
	"testing"
)

func TestScan(t *testing.T) {
	points, pointsMap := examplePoints()
	for _, p := range Scan(points) {
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

func TestProblem2(t *testing.T) {
	points, pointsMap := problem53Points()
	for _, p := range Scan(points) {
		t.Log(p, pointsMap[p])
	}
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

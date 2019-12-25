package inPolygon

import (
	"algorithm-class/geometry/plane"
	"testing"

	"github.com/stretchr/testify/assert"
)

type inTest struct {
	Polygon []plane.Point
	In      []plane.Point
	Out     []plane.Point
}

func TestIn(t *testing.T) {
	tests := []inTest{
		inTest{
			Polygon: []plane.Point{
				{2, 3}, {4, 3}, {5, 2}, {8, 2}, {9, 4},
				{10, 6}, {8, 7}, {8, 8}, {8, 9}, {7, 9},
				{6, 9}, {5, 9}, {4, 8}, {6, 7}, {6, 6},
				{7, 5}, {5, 5}, {4, 6}, {3, 7}, {2, 8},
				{2, 6}, {4, 5}, {3, 4},
			},
			In: []plane.Point{
				{3, 6},
				{4, 4},
				{5, 3},
				{5, 4},
				{5, 8},
				{6, 3},
				{6, 4},
				{6, 8},
				{7, 3},
				{7, 4},
				{7, 8},
				{8, 3},
				{8, 4},
				{8, 5},
				{9, 5},
			},
			Out: []plane.Point{
				{1, 10},
				{1, 1},
				{1, 2},
				{1, 3},
				{1, 4},
				{1, 5},
				{1, 6},
				{1, 7},
				{1, 8},
				{1, 9},
				{10, 10},
				{10, 1},
				{10, 2},
				{10, 3},
				{10, 4},
				{10, 4},
				{10, 5},
				{10, 7},
				{10, 8},
				{10, 9},
				{2, 10},
				{2, 2},
				{2, 4},
				{2, 5},
				{2, 9},
				{3, 2},
				{3, 5},
				{3, 8},
				{4, 2},
				{4, 7},
				{4, 9},
				{5, 6},
				{5, 7},
				{9, 10},
				{9, 2},
				{9, 3},
				{9, 7},
				{9, 8},
				{9, 9},
			},
		},
	}
	for _, tc := range tests {
		for _, x := range tc.In {
			assert.True(t, In(tc.Polygon, x), "%v", x)
		}
		for _, x := range tc.Out {
			assert.False(t, In(tc.Polygon, x), "%v", x)
		}
	}
}

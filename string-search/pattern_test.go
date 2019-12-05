package stringSearch

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPatternProblem1(t *testing.T) {
	type Tc struct {
		s string
		i int
	}

	type Test struct {
		patternString string
		pattern       *Pattern
		tcs           []Tc
	}

	tests := []Test{
		Test{
			patternString: "(A*B+AC)D",
			pattern: &Pattern{
				Chars: []byte{' ', 'A', ' ', 'B', ' ', ' ', 'A', 'C', 'D', ' '},
				Next1: []int{5, 2, 3, 4, 8, 6, 7, 8, 9, 0},
				Next2: []int{5, 2, 1, 4, 8, 2, 7, 8, 9, 0},
			},
			tcs: []Tc{
				{"AABD", 4},
				{"", -1},
				{"DE", -1},
				{"ABD", 3},
				{"AAABD", 5},
				{"AAABDD", 5},
				{"AD", -1},
				{"ACD", 3},
			},
		},
		Test{
			patternString: "(A+B)*C",
			pattern: &Pattern{
				Chars: []byte{' ', ' ', ' ', 'A', 'B', 'C', ' '},
				Next1: []int{1, 2, 3, 1, 1, 6, 0},
				Next2: []int{1, 5, 4, 1, 1, 6, 0},
			},
			tcs: []Tc{
				{"", -1},
				{"AC", 2},
				{"BC", 2},
				{"ABABC", 5},
				{"ABAAAAC", 7},
			},
		},
		Test{
			patternString: "(AB*+A*D)E",
			pattern: &Pattern{
				Chars: []byte{' ', 'A', ' ', 'B', 'E', ' ', ' ', 'A', 'D'},
				Next1: []int{1, 2, 3, 2, 5, 0, 7, 6, 4},
				Next2: []int{6, 2, 4, 2, 5, 0, 8, 0, 4},
			},
			tcs: []Tc{
				{"AE", 2},
				{"ABBBE", 5},
				{"AAAADE", 6},
			},
		},
		Test{
			patternString: "(A+B)*(C+D)*E",
			pattern: &Pattern{
				Chars: []byte{' ', ' ', ' ', 'E', ' ', 'A', 'B', ' ', 'C', 'D', ' ', ' ', ' '},
				Next1: []int{1, 2, 3, 4, 0, 7, 7, 1, 10, 10, 2, 5, 8},
				Next2: []int{1, 11, 12, 4, 0, 7, 7, 1, 10, 10, 2, 6, 9},
			},
			tcs: []Tc{
				{"E", 1},
				{"AE", 2},
				{"CE", 2},
				{"ACE", 3},
				{"BCE", 3},
				{"BBCCDDE", 7},
				{"ABCDE", 5},
				{"ABCDEE", 5},
				{"AC", -1},
				{"", -1},
				{"ABE", 3},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.patternString, func(t *testing.T) {
			for _, tc := range test.tcs {
				state := NewState(test.pattern, tc.s)
				for state.Next() {
				}
				if !assert.Equal(t, tc.i >= 0, state.Ok(), tc.s) {
					continue
				}
				if !(tc.i >= 0 && assert.Equal(t, tc.i, state.Pos())) {
					continue
				}
				logMatch(t, test.pattern, test.patternString, tc.s)
			}
		})
	}
}

func logMatch(t *testing.T, pattern *Pattern, patternString, str string) {
	if !testing.Verbose() {
		return
	}
	t.Logf("pattern = %s, string = %s", patternString, str)
	state := NewState(pattern, str)
	printDeque(t, state.q)
	for state.Next() {
		printDeque(t, state.q)
	}
	t.Logf("match = %t, position = %d", state.Ok(), state.Pos())
}

var rowPool sync.Pool

func getRow() []int {
	row, ok := rowPool.Get().([]int)
	if !ok {
		return nil
	}
	return row[:0]
}

func printDeque(t *testing.T, q *Deque) {
	row := getRow()
	defer rowPool.Put(row)
	for r := q.Tail; r != nil; r = r.Prev {
		row = append(row, r.Value)
	}
	t.Logf("%d", row)
}

func TestDeque(t *testing.T) {
	q := NewDeque()
	require.True(t, q.IsEmpty())
	q.Append(1)
	require.False(t, q.IsEmpty())
	require.NotNil(t, q.Head)
	require.NotNil(t, q.Tail)

	q.Append(2)
	require.Equal(t, q.Pop(), 2)

	q.Prepend(0)
	require.Equal(t, q.Head.Value, 0)
	require.Equal(t, q.Tail.Value, 1)

	require.Equal(t, q.Pop(), 1)
	require.Equal(t, q.Pop(), 0)
	require.True(t, q.IsEmpty())
	require.Nil(t, q.Head)
	require.Nil(t, q.Tail)
}

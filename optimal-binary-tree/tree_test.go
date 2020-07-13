package optimalBinaryTree

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

func TestNewPrior(t *testing.T) {
	r := NewPrior([]float64{0.3, 0.2, 0.4, 0.1})
	t.Log(r.Min())
	t.Log(r.Costs.s)
	t.Log(r.Decisions.s)
	t.Fail()
}

func TestExam(t *testing.T) {
	s := []float64{0.02, 0.26, 0.12, 0.21, 0.39}
	r := NewPrior(s)
	t.Log(r.Min())
	costs := r.Costs.s
	for len(costs) > 0 {
		t.Log(costs[:len(s)])
		costs = costs[len(s):]
	}
	decs := r.Decisions.s
	for len(decs) > 0 {
		t.Log(decs[:len(s)])
		decs = decs[len(s):]
	}
	testSimulation(t, r)
}

func TestTree(t *testing.T) {
	probabilities := []float64{0.3, 0.2, 0.4, 0.1}
	tree := NewPrior(probabilities).Tree()
	require.Equal(t, tree.Search(0), 2)
	require.Equal(t, tree.Search(1), 3)
	require.Equal(t, tree.Search(2), 1)
	require.Equal(t, tree.Search(3), 2)
}

func TestNewPrior2(t *testing.T) {
	r := NewPrior2([]float64{0.3, 0.2, 0.4, 0.1})
	t.Log(r.Min())
	t.Log(r.Costs.s)
	t.Log(r.Decisions.s)
	t.Fail()
}

func TestSimulation(t *testing.T) {
	testSimulation(t, NewPrior([]float64{0.3, 0.2, 0.4, 0.1}))
	testSimulation(t, NewPrior2([]float64{0.3, 0.2, 0.4, 0.1}))
}

func testSimulation(t *testing.T, prior *Prior) {
	const trials = 10000
	samples := simulate(prior.Tree(), prior.Probabilities, trials)

	exp := expectation(samples)
	std := standardDeviation(samples, exp, sampleSpace) / math.Sqrt(trials)

	low := exp - alpha5*std
	high := exp + alpha5*std
	t.Logf("CI = [%.4f, %.4f]", low, high)
	t.Logf("Min = %.4f, Exp = %.4f, Std = %.4f", prior.Min(), exp, std)

	if exp <= low || exp >= high {
		t.Fail()
	}
}

type widelyTraverse struct {
	cur []*Node
	buf []*Node
}

func (w *widelyTraverse) Hint(n int) {
	w.buf = make([]*Node, 0, (n+1)/2)
}

func (w *widelyTraverse) Next() bool {
	buf := w.buf[:0]
	for _, node := range w.cur {
		if node.Left != nil {
			buf = append(buf, node.Left)
		}
		if node.Right != nil {
			buf = append(buf, node.Right)
		}
	}
	w.cur, w.buf = buf, w.cur
	return len(buf) > 0
}

func newWidelyTraverse(cur ...*Node) widelyTraverse {
	return widelyTraverse{cur: cur}
}

func treeCosts(tree *Tree) []float64 {
	traverse := newWidelyTraverse(tree.root)
	traverse.Hint(tree.Size())
	discrete := make([]float64, 0, tree.Size())
For:
	u := probabilityUnion(traverse.cur)
	discrete = append(discrete, u)
	if traverse.Next() {
		goto For
	}
	return discrete
}

func probabilityUnion(nodes []*Node) float64 {
	var u float64
	for _, node := range nodes {
		u += node.Probability
	}
	return u
}

const alpha5 = 1.96

func simulate(tree *Tree, probabilities []float64, trials int) []float64 {
	samples := make([]float64, 0, trials)
	for i := 1; i <= trials; i++ {
		x := choose(probabilities)
		h := tree.Search(x)
		samples = append(samples, float64(h))
	}
	return samples
}

func choose(probabilities []float64) int {
	if len(probabilities) == 0 {
		panic("empty probability space")
	}
	scale := 1.0
	for i, p := range probabilities[:len(probabilities)-1] {
		if scale*rand.Float64() <= p {
			return i
		}
		scale -= p
	}
	return len(probabilities) - 1
}

func expectation(samples []float64) float64 {
	if len(samples) == 0 {
		return missingValue()
	}
	var exp float64
	for _, sample := range samples {
		if isMissingValue(sample) {
			continue
		}
		exp += sample
	}
	return exp / float64(len(samples))
}

func variance(samples []float64, exp float64, degrees int) float64 {
	switch len(samples) {
	case 0:
		return missingValue()
	case 1:
		return 0
	}
	if isMissingValue(exp) {
		exp = expectation(samples)
	}
	degrees = determineDegrees(len(samples), degrees)
	var v float64
	for _, sample := range samples {
		dev := sample - exp
		v += dev * dev
	}
	return v / float64(degrees)
}

func standardDeviation(samples []float64, exp float64, degrees int) float64 {
	return math.Sqrt(variance(samples, exp, degrees))
}

const (
	populationSpace = -1
	sampleSpace     = -2
)

func determineDegrees(samplesLen, degrees int) int {
	switch degrees {
	case populationSpace:
		return samplesLen
	case sampleSpace:
		return samplesLen - 1
	}
	return degrees
}

func missingValue() float64 {
	return math.NaN()
}

func isMissingValue(x float64) bool {
	return math.IsNaN(x) || math.IsInf(x, 0)
}

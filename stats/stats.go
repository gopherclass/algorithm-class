package stats

import (
	"fmt"
	"math"
)

const (
	PopulationSpace = -1
	SampleSpace     = -2
)

type Samples struct {
	Samples      []float64
	xExpectation float64
	xVariance    float64
}

func Sampling(tries int, sampling func() float64) *Samples {
	samples := make([]float64, 0, tries)
	for i := 1; i <= tries; i++ {
		sample := sampling()
		samples = append(samples, sample)
	}
	return NewSamples(samples)
}

func NewSamples(samples []float64) *Samples {
	return &Samples{
		Samples:      samples,
		xExpectation: MissingValue(),
		xVariance:    MissingValue(),
	}
}

func (s Samples) Copy() *Samples {
	return &s
}

func (s *Samples) Expectation() float64 {
	if IsMissingValue(s.xExpectation) {
		s.xExpectation = Expectation(s.Samples)
	}
	return s.xExpectation
}

func (s *Samples) Variance(n int) float64 {
	if IsMissingValue(s.xVariance) {
		s.xVariance = s.getVariance()
	}
	n = getDegree(len(s.Samples), n)
	return s.xVariance / float64(n)
}

func getDegree(samples, n int) int {
	switch n {
	case PopulationSpace:
		return samples
	case SampleSpace:
		return samples - 1
	}
	return n
}

func (s *Samples) getVariance() float64 {
	switch len(s.Samples) {
	case 0:
		return MissingValue()
	case 1:
		return 0
	}
	exp := s.Expectation()
	var v float64
	for _, sample := range s.Samples {
		dev := sample - exp
		v += dev * dev
	}
	return v
}

func (s *Samples) StandardDeviation(n int) float64 {
	return math.Sqrt(s.Variance(n))
}

type CI struct {
	Expectation float64
	Margin      float64
	Alpha       float64
}

func (i CI) String() string {
	return fmt.Sprintf("%.4f±%.4f", i.Expectation, i.Margin)
}

func SampleMean(s *Samples, alpha float64) CI {
	std := s.StandardDeviation(SampleSpace)
	return CI{
		Expectation: alpha,
		Margin:      alpha * std,
		Alpha:       alpha,
	}
}

func Expectation(samples []float64) float64 {
	if len(samples) == 0 {
		return MissingValue()
	}
	var exp float64
	for _, sample := range samples {
		if IsMissingValue(sample) {
			continue
		}
		exp += sample
	}
	return exp / float64(len(samples))
}

func MissingValue() float64 {
	return math.NaN()
}

func IsMissingValue(x float64) bool {
	return math.IsNaN(x) || math.IsInf(x, 0)
}

package stats

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpectation(t *testing.T) {
	require.Equal(t, 3.0, newSamples(1, 2, 3, 4, 5).Expectation())
}

func newSamples(samples ...float64) *Samples {
	return NewSamples(samples)
}

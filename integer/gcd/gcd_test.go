package gcd

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

func TestX(t *testing.T) {
	rng := newrng()
	for i := 0; i < 10000; i++ {
		a, b := rng.Int(), rng.Int()
		x, y, r := xgcd(a, b)
		require.Equal(t, r, a*x+b*y, "%d*%d + %d*%d != %d",
			a, x, b, y, r)
	}
}

func TestMuliplicativeInverse(t *testing.T) {
	rng := newrng()
	One := big.NewInt(1)
	for i := 0; i < 1000; i++ {
		a, b := twoCoprime(rng)
		i := MultiplicativeInverse(a, b)

		var x, y, z big.Int
		x.SetInt64(int64(a))
		y.SetInt64(int64(b))
		z.SetInt64(int64(i))
		x.Mul(&x, &z)
		x.Rem(&x, &y)
		if x.Sign() < 0 {
			x.Add(&x, &y)
		}
		require.True(t, x.Cmp(One) == 0, "%d*%d != 1 (mod %d)",
			a, i, b)
	}
}

func twoCoprime(rng *rand.Rand) (int, int) {
	for {
		a, b := rng.Int(), rng.Int()
		if IsCoprime(a, b) {
			return a, b
		}
	}
}

func newrng() *rand.Rand {
	rngSeed := uint64(time.Now().UnixNano())
	rngSource := rand.NewSource(rngSeed)
	rng := rand.New(rngSource)
	return rng
}

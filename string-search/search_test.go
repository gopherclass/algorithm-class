package stringSearch

import (
	"algorithm-class/inst"
	"flag"
	"strings"
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const alphabetSpace = "abcdefghijklmnopqrstuvwxyz "

var seedValue uint64

var rngSourceGlobal *rand.Rand

func init() {
	flag.Uint64Var(&seedValue, "seed", 0, "Seed value")
}

func rngSource() *rand.Rand {
	if rngSourceGlobal == nil {
		if seedValue == 0 {
			seedValue = uint64(time.Now().UnixNano())
		}
		rngSourceGlobal = rand.New(rand.NewSource(seedValue))
	}
	return rngSourceGlobal
}

func TestNaiveSearch(t *testing.T) {
	testSearchFunc(t, NaiveSearch, alphabet)
}

func testSearchFunc(t *testing.T, searchFunc func(ic *inst.Counter, str, pat string) int, charsTable string) {
	strBuilder := new(stringBuilder).setTable(charsTable)
	patBuilder := new(stringBuilder).setTable(charsTable)
	const Loop = 300
	for i := 0; i < Loop; i++ {
		patSize := rngSource().Intn(i+1) + 1
		str, pat := strBuilder.Build(i), patBuilder.Build(patSize)
		expected := strings.Index(str, pat)
		got := searchFunc(nil, str, pat)
		if expected != got {
			t.Fatalf("str = %q, pat = %q, %d expected, but got %d",
				str, pat, expected, got)
		}
	}
}

type stringBuilder struct {
	buf        strings.Builder
	charsTable string
}

func (r *stringBuilder) setTable(charsTable string) *stringBuilder {
	r.charsTable = charsTable
	return r
}

// Build is unsafe.
func (r *stringBuilder) Build(n int) string {
	buf := &r.buf
	buf.Reset()
	buf.Grow(n)
	r.buildString(buf, n)
	return buf.String()
}

func (r *stringBuilder) buildString(buf *strings.Builder, n int) {
	for i := 0; i < n; i++ {
		buf.WriteByte(r.char())
	}
}

func (r *stringBuilder) char() byte {
	n := len(r.charsTable)
	i := rngSource().Intn(n)
	return r.charsTable[i]
}

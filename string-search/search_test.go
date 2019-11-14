package stringSearch

import (
	"algorithm-class/inst"
	"flag"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

func TestNaiveSearchProblem1(t *testing.T) {
	if !testing.Verbose() {
		t.Skip()
	}
	str := "ababacabcbababcacacbcaababca"
	pat := "ababca"
	logCompare(t, NaiveSearch, str, pat)
	ensureTermination(t)
}

func TestKMPPrecomputedTableProblem2(t *testing.T) {
	t.Log(KMPPrecomputedTable(nil, "ababca"))
	t.Log(KMPPrecomputedTable(nil, "abababca"))
	t.Log(KMPPrecomputedTable(nil, "abcbabcbabc"))
	t.Log(KMPPrecomputedTable(nil, "abracadabra"))
	ensureTermination(t)
}

func TestKMPImprovedPrecomputedTableProblem3(t *testing.T) {
	t.Log(getKMPImprovedPrecomputedTable(nil, "ababca"))
	t.Log(getKMPImprovedPrecomputedTable(nil, "abababca"))
	t.Log(getKMPImprovedPrecomputedTable(nil, "abcbabcbabc"))
	t.Log(getKMPImprovedPrecomputedTable(nil, "abracadabra"))
	ensureTermination(t)
}

func getKMPImprovedPrecomputedTable(ic *inst.Counter, pat string) []int {
	next := KMPPrecomputedTable(ic, pat)
	next = KMPImprovedPrecomputedTable(ic, pat, next)
	return next
}

func TestKMPSearchProblem4(t *testing.T) {
	str := "ababacabcbababcacacaababca"
	pat := "ababca"
	logCompare(t, KMPSearch, str, pat)
	ensureTermination(t)
}

func TestKMPSearch(t *testing.T) {
	testSearchFunc(t, KMPSearch, alphabet)
}

func TestBoyerMooreSearchCase1(t *testing.T) {
	str := "pphxpbbj"
	pat := "ob"
	require.Equal(t, BoyerMooreSearch(nil, str, pat), strings.Index(str, pat))
	ensureTermination(t)
}

func TestPresentedBoyerMooreSearchCase1(t *testing.T) {
	str := "pphxpbbj"
	pat := "ob"
	require.Equal(t, presentedBoyerMooreSearch(nil, str, pat), strings.Index(str, pat))
	ensureTermination(t)
}

func TestBoyerMooreSearch(t *testing.T) {
	testSearchFunc(t, BoyerMooreSearch, alphabet)
}

func TestPresentedBoyerMooreSearch(t *testing.T) {
	testSearchFunc(t, presentedBoyerMooreSearch, alphabet)
}

func presentedBoyerMooreSearch(ic *inst.Counter, str, pat string) int {
	if len(str) < len(pat) {
		return -1
	}
	skip := BoyerMooreBadCharSkip(pat)
	i := len(pat) - 1
	j := len(pat) - 1
	for j >= 0 {
		for str[i] != pat[j] {
			s := skip[str[i]]
			if len(pat)-j > s {
				i += len(pat) - j
			} else {
				i += s
			}
			if i >= len(str) {
				return -1
			}
			j = len(pat) - 1
		}
		i--
		j--
	}
	return i + 1
}

func TestRabinKarpPower(t *testing.T) {
	require.Equal(t, _D%_Q, RabinKarpPowerHash(2))
	require.Equal(t, _D*_D%_Q, RabinKarpPowerHash(3))
	require.Equal(t, _D*_D*_D%_Q, RabinKarpPowerHash(4))
}

func TestRabinKarpHash(t *testing.T) {
	require.Equal(t, (uint64('a')*_D+uint64('b'))%_Q, RabinKarpHash("ab"))

	power := RabinKarpPowerHash(2)
	h := RabinKarpHash("ab")
	h = RabinKarpSlidingHash(h, power, uint64('a'), uint64('c'))
	require.Equal(t, RabinKarpHash("bc"), h)
	h = RabinKarpSlidingHash(h, power, uint64('b'), uint64('a'))
	require.Equal(t, RabinKarpHash("ca"), h)
	h = RabinKarpSlidingHash(h, power, uint64('c'), uint64('b'))
	require.Equal(t, RabinKarpHash("ab"), h)
}

func TestRabinKarpSearch(t *testing.T) {
	testSearchFunc(t, RabinKarpSearch, alphabet)
}

func TestRabinKarpSearchProblem6(t *testing.T) {
	str := "STRING STARTING CONSISTING"
	pat := "STING"
	index := logRabinKarpSearch(t, nil, str, pat)
	require.Equal(t, index, 21)
	t.Logf("str = %q, pat = %q: Index = %d",
		str, pat, index)
}

func logRabinKarpSearch(t *testing.T, ic *inst.Counter, str, pat string) int {
	rejectCounter(ic)
	if len(str) < len(pat) {
		return -1
	}

	index := func(s string, i int) uint64 {
		switch {
		case 'A' <= s[i] && s[i] <= 'Z':
			return uint64(s[i]) - 'A' + 1
		case s[i] == ' ':
			return 0
		default:
			panic(fmt.Sprintf("unsupported alphabet %c", s[i]))
		}
	}
	RabinKarpHash := func(s string) (h uint64) {
		for i := range s {
			h *= _D
			h += index(s, i)
			h %= _Q
		}
		return h
	}

	RabinKarpPowerHash := func(n int) (power uint64) {
		power = 1
		for i := 1; i < n; i++ {
			power *= _D
			power %= _Q
		}
		return power
	}

	RabinKarpSlidingHash := func(strHash, power, dropped, slided uint64) uint64 {
		strHash -= dropped * power
		strHash *= _D
		strHash += slided
		strHash %= _Q
		return strHash
	}

	patHash := RabinKarpHash(pat)
	strHash := RabinKarpHash(str[:len(pat)])
	power := RabinKarpPowerHash(len(pat))
	i := 0
	t.Logf("pattern hash: %d", patHash)
	t.Logf("string hash : %d", strHash)
	for {
		if strHash == patHash && str[i:i+len(pat)] == pat {
			return i
		}
		if i+len(pat) >= len(str) {
			break
		}
		strHash = RabinKarpSlidingHash(
			strHash,
			power,
			index(str, i),
			index(str, i+len(pat)),
		)
		t.Logf("string hash : %d", strHash)
		i++
	}
	return -1
}

type searchFunc func(ic *inst.Counter, str, pat string) int

func testSearchFunc(t *testing.T, searchFunc searchFunc, charsTable string) {
	strBuilder := new(stringBuilder).setTable(charsTable)
	patBuilder := new(stringBuilder).setTable(charsTable)
	const Loop = 500
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

func logCompare(t *testing.T, searchFunc searchFunc, str, pat string) {
	ic := inst.NewCounter()
	index := searchFunc(ic, str, pat)
	compare := ic.State()[inst.Compare]
	t.Logf("str = %q, pat = %q: Index = %d, Compare = %d",
		str, pat, index, compare)
}

func ensureTermination(t *testing.T) {}

// 1번 문제 10, 22 26아닌가?

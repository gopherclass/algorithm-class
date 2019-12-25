package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

func binomial(n, r, overflow int64) int64 {
	var x int64 = 1
	if n-r < r {
		r = n - r
	}
	for i := int64(1); i <= r; i++ {
		x *= n - r + i
		x /= i
		if x < 0 || overflow < x {
			return overflow + 1
		}
	}
	return x
}

type C struct {
	n, r int64
}

func solveMain(r *Reader) {
	n := r.Int64()
	var res []C
	for k := int64(0); binomial(2*k, k, n) <= n; k++ {
		lo, hi := 2*k, n

		for lo < hi {
			mi := (lo + hi) / 2
			if binomial(mi, k, n) < n {
				lo = mi + 1
			} else {
				hi = mi
			}
		}
		if binomial(lo, k, n) == n {
			res = append(res, C{lo, k})
			if 2*k < lo {
				res = append(res, C{lo, lo - k})
			}
		}
	}
	sort.Sort(byRes(res))
	fmt.Println(len(res))
	for i := range res {
		fmt.Println(res[i].n, res[i].r)
	}
}

type byRes []C

func (s byRes) Len() int      { return len(s) }
func (s byRes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byRes) Less(i, j int) bool {
	return s[i].n < s[j].n || (s[i].n == s[j].n && s[i].r <= s[j].r)
}

func main() {
	solveMain(Stdin())
}

type Reader struct {
	r *bufio.Reader
}

func (r *Reader) Int() int {
	n, err := strconv.Atoi(r.String())
	r.check(err)
	return n
}

func (r *Reader) Ints(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = r.Int()
	}
	return s
}

func (r *Reader) Int64() int64 {
	n, err := strconv.ParseInt(r.String(), 10, 64)
	r.check(err)
	return n
}

func (r *Reader) String() string {
	buf, err := r.Scan(r.IsSeparating)
	r.check(err)
	return string(buf)
}

func (r *Reader) Strings(n int) []string {
	s := make([]string, n)
	for i := range s {
		s[i] = r.String()
	}
	return s
}

func (*Reader) IsSeparating(r rune) bool {
	return r == ' ' || r == '\n'
}

func (r *Reader) Scan(isInvalid func(rune) bool) ([]byte, error) {
	var buf []byte
	var enc [utf8.UTFMax]byte
	for {
		r0, _, err := r.r.ReadRune()
		if err == io.EOF && len(buf) > 0 {
			return buf, nil
		}
		if err != nil {
			return buf, err
		}
		if isInvalid(r0) {
			return buf, nil
		}
		size := utf8.EncodeRune(enc[:], r0)
		buf = append(buf, enc[:size]...)
	}
}

func (*Reader) check(err error) {
	if err != nil {
		panic(err)
	}
}

func New(r io.Reader) *Reader {
	return &Reader{
		r: bufio.NewReader(r),
	}
}

func Stdin() *Reader {
	return New(os.Stdin)
}

func Example(s string) *Reader {
	return New(strings.NewReader(s))
}

// ---

// TODO
type Test struct {
	Input  string
	Output string
}

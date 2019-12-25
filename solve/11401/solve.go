package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

const P = 1000000007

func solveMain(r *Reader) {
	n, k := r.Int(), r.Int()

	factorial := make([]int64, 4000001)
	factorial[0] = 1
	factorial[1] = 1
	for i := 2; i <= n; i++ {
		factorial[i] = (factorial[i-1] * int64(i)) % P
	}

	inverseFactorial := make([]int64, 4000001)
	inverseFactorial[n] = int64(multiplicativeInverse(int(factorial[n]), P))
	for i := n - 1; i >= 0; i-- {
		inverseFactorial[i] = (int64(i+1) * inverseFactorial[i+1]) % P
	}

	C := (factorial[n] * inverseFactorial[k]) % P
	C = (C * inverseFactorial[n-k]) % P

	fmt.Println(C)
}

func main() {
	solveMain(Stdin())
}

// a*s + b*t = gcd(a, b) = r
func xgcd(a, b int) (s int, t int, r int) {
	s0, s1 := 1, 0
	t0, t1 := 0, 1
	r0, r1 := a, b
	for r1 != 0 {
		quotient := r0 / r1
		r0, r1 = r1, r0-quotient*r1
		s0, s1 = s1, s0-quotient*s1
		t0, t1 = t1, t0-quotient*t1
	}
	return s0, t0, r0
}

func gcd(a, b int) int {
	r0, r1 := a, b
	for r1 != 0 {
		quotient := r0 / r1
		r0, r1 = r1, r0-quotient*r1
	}
	return r0
}

func isCoprime(a, b int) bool {
	return gcd(a, b) == 1
}

// a*s = 1 (mod b)
func multiplicativeInverse(a, b int) (s int) {
	s, _, _ = xgcd(a, b)
	return s
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

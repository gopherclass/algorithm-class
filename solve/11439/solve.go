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

func getPrimes(n int64) []int64 {
	s := make([]int64, 1, n)
	s[0] = 2
Next:
	for i := int64(3); i <= n; i += 2 {
		for _, p := range s {
			if i%p == 0 {
				continue Next
			}
		}
		s = append(s, i)
	}
	return s
}

func factorization(pow []int64, n int64, primes []int64) {
	for i, p := range primes {
		m := n
		for m > 0 {
			m /= p
			pow[i] += m
		}
	}
}

func pow(x, y, m int64) int64 {
	var res int64 = 1
	b := x % m
	for y > 0 {
		if y&1 == 1 {
			res *= b
			res %= m
		}
		b *= b
		b %= m
		y >>= 1
	}
	return res
}

func solveMain(r *Reader) {
	N, K, M := r.Int64(), r.Int64(), r.Int64()
	primes := getPrimes(334)

	isPrime := true
	for _, p := range primes {
		if N%p == 0 {
			isPrime = false
			break
		}
	}
	if isPrime {
		primes = append(primes, N)
	}

	A := make([]int64, len(primes))
	B := make([]int64, len(primes))
	factorization(A, N, primes)
	factorization(B, K, primes)
	factorization(B, N-K, primes)

	for i := 0; i < len(primes); i++ {
		A[i] -= B[i]
	}

	var res int64 = 1
	for i := 0; i < len(primes); i++ {
		if A[i] > 0 {
			fmt.Println(primes[i], A[i], B[i])
		}
		res *= pow(int64(primes[i]), A[i], M)
		res %= M
	}
	fmt.Println(res)
}

func main() {
	solveMain(Example(`1111111 3 88`))
	return
	solveMain(Example(`282 8 190`))
	solveMain(Example(`242 8 190`))
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

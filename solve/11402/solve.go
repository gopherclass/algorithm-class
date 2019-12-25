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

func solveMain(r *Reader) {
	N, K, M := r.Int64(), r.Int64(), r.Int()

	nCr := make([][]int, M+1)
	nCr[0] = []int{1}
	for i := 1; i <= M; i++ {
		row := make([]int, i+1)
		row[0] = 1
		row[i] = 1
		prev := nCr[i-1]
		for j := 1; j < i; j++ {
			row[j] = (prev[j] + prev[j-1]) % M
		}
		nCr[i] = row
	}

	var m, n int
	var res int = 1
	for N > 0 || K > 0 {
		N, m = N/int64(M), int(N%int64(M))
		K, n = K/int64(M), int(K%int64(M))
		if m < n {
			res = 0
			break
		}
		res *= nCr[m][n]
		res %= M
	}
	fmt.Println(res)
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

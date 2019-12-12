// https://www.acmicpc.net/problem/11758
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

type Point struct {
	X, Y int
}

func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

func Det(p, q Point) int {
	return p.X*q.Y - p.Y*q.X
}

func solveMain(r *Reader) {
	x := Point{r.Int(), r.Int()}
	y := Point{r.Int(), r.Int()}
	z := Point{r.Int(), r.Int()}
	fmt.Println(Sign(Dir(x, y, z)))
}

func Dir(x, y, z Point) int {
	return Det(y.Sub(x), z.Sub(x))
}

func Sign(x int) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
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

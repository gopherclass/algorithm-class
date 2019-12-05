package randomize

import (
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

const (
	LowerAlphabets      = "abcdefghijklmnopqrstuvwxyz"
	UpperAlphabets      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LowerAlphabetsSpace = "abcdefghijklmnopqrstuvwxyz "
	UpperAlphabetsSpace = "ABCDEFGHIJKLMNOPQRSTUVWXYZ "
)

type RuneSource interface {
	Get() rune
	Len() int
}

type IntSource interface {
	Get() int
	Len() int
}

type String struct {
	Rand   *rand.Rand
	String string
}

func (s String) Len() int { return len(s.String) }
func (s String) Get() rune {
	i := s.Rand.Intn(len(s.String))
	return rune(s.String[i])
}

type Runes struct {
	Rand  *rand.Rand
	Runes []rune
}

func (s Runes) Len() int { return len(s.Runes) }
func (s Runes) Get() rune {
	i := s.Rand.Intn(len(s.Runes))
	return s.Runes[i]
}

type Ints struct {
	Rand  *rand.Rand
	Table []float64
}

func (s Ints) Len() int { return len(s.Table) }
func (s Ints) Get() int {
	u := 1.0
	for i, p := range s.Table {
		x := s.Rand.Float64()
		if u*x <= p {
			return i
		}
		u -= p
	}
	return -1
}

func Now() *rand.Rand {
	now := time.Now()
	return New(uint64(now.UnixNano()))
}

func New(seed uint64) *rand.Rand {
	src := rand.NewSource(seed)
	return rand.New(src)
}

type StringBuilder struct {
	buf strings.Builder
	src RuneSource
}

// unsafe
func (s *StringBuilder) Build(size int) string {
	s.buf.Reset()
	s.buf.Grow(size)
	for i := 0; i < size; i++ {
		s.buf.WriteRune(s.src.Get())
	}
	return s.buf.String()
}

func NewStringBuilder(src RuneSource) *StringBuilder {
	return &StringBuilder{src: src}
}

package stringSearch

import (
	"algorithm-class/inst"
)

func rejectCounter(ic *inst.Counter) {
	if ic != nil {
		panic("inst.Counter argument is rejected")
	}
}

func NaiveSearch(ic *inst.Counter, str, pat string) int {
	i, j := 0, 0
	for i < len(str) && j < len(pat) {
		ic.Once(inst.Compare)
		if str[i] != pat[j] {
			i = i - j + 1
			j = 0
		} else {
			i++
			j++
		}
	}
	if j >= len(pat) {
		return i - j
	}
	return -1
}

func InitNext(ic *inst.Counter, pat string) []int {
	return KMPPrecomputedTable(ic, pat)
}

// TODO: Reimplment KMPPrecoputedTable with recursive theory
func KMPPrecomputedTable(ic *inst.Counter, pat string) []int {
	next := make([]int, len(pat))
	for i, j := 0, -1; i < len(pat); i, j = i+1, j+1 {
		next[i] = j
		for j >= 0 && pat[i] != pat[j] {
			j = next[j]
		}
	}
	return next
}

func KMPImprovedPrecomputedTable(ic *inst.Counter, pat string, next []int) []int {
	for i := range next {
		for next[i] >= 0 && pat[i] == pat[next[i]] {
			next[i] = next[next[i]]
		}
	}
	return next
}

func KMPSearch(ic *inst.Counter, str, pat string) int {
	// TODO: Precomputed table에서 수행한 비교 횟수도 포함해야 하나?
	// TODO: inst.Counter 구조의 설계를 바꿔야 할 시점이 다시 오는 것 같다.
	next := KMPPrecomputedTable(ic, pat)
	return doKMPSearch(ic, str, pat, next)
}

func doKMPSearch(ic *inst.Counter, str, pat string, next []int) int {
	i, j := 0, 0
	for i < len(str) && j < len(pat) {
		ic.Once(inst.Compare)
		if str[i] == pat[j] {
			i++
			j++
		} else {
			j = next[j]
			if j < 0 {
				i++
				j = 0
			}
		}
	}
	if j >= len(pat) {
		return i - j
	}
	return -1
}

func KMPImprovedSearch(ic *inst.Counter, str, pat string) int {
	// TODO: Precomputed table에서 수행한 비교 횟수도 포함해야 하나?
	next := KMPPrecomputedTable(ic, pat)
	next = KMPImprovedPrecomputedTable(ic, pat, next)
	return doKMPSearch(ic, str, pat, next)
}

func BoyerMooreBadCharSkip(pat string) (badCharSkip [256]int) {
	for i := range badCharSkip {
		badCharSkip[i] = len(pat)
	}
	n := len(pat) - 1
	for i, r := range pat {
		badCharSkip[r] = n - i
	}
	return badCharSkip
}

func BoyerMooreSearch(ic *inst.Counter, str, pat string) int {
	badCharSkip := BoyerMooreBadCharSkip(pat)
	i := len(pat) - 1
	for i < len(str) {
		j := len(pat) - 1
		for j >= 0 && ic.Once(inst.Compare) && str[i] == pat[j] {
			i--
			j--
		}
		if j < 0 {
			return i + 1
		}
		i += maxInt(badCharSkip[str[i]], len(pat)-j)
	}
	return -1
}

func maxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}

const _Q uint64 = 33554393
const _D uint64 = 32

func RabinKarpHash(s string) (h uint64) {
	for i := range s {
		h *= _D
		h += uint64(s[i])
		h %= _Q
	}
	return h
}

func RabinKarpPowerHash(n int) (power uint64) {
	power = 1
	for i := 1; i < n; i++ {
		power *= _D
		power %= _Q
	}
	return power
}

func RabinKarpSlidingHash(strHash, power, dropped, slided uint64) uint64 {
	strHash -= dropped * power
	strHash *= _D
	strHash += slided
	strHash %= _Q
	return strHash
}

func RabinKarpSearch(ic *inst.Counter, str, pat string) int {
	rejectCounter(ic)
	if len(str) < len(pat) {
		return -1
	}
	patHash := RabinKarpHash(pat)
	strHash := RabinKarpHash(str[:len(pat)])
	power := RabinKarpPowerHash(len(pat))
	i := 0
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
			uint64(str[i]),
			uint64(str[i+len(pat)]),
		)
		i++
	}
	return -1
}

// Harris Corner detection
// Auto correlation

// Automatic Scale Selection
// Sigma blurring

// 이미지를 Blurring한 것은 이미지를 멀리서 본것과 비슷하며 이는 Scaling한 효과와 비슷하다.

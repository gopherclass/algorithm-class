package stringSearch

import "algorithm-class/inst"

func rejectCounter(ic *inst.Counter) {
	if ic != nil {
		panic("inst.Counter argument is rejected")
	}
}

func NaiveSearch(ic *inst.Counter, str, pat string) int {
	rejectCounter(ic)
	i, j := 0, 0
	for i < len(str) && j < len(pat) {
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

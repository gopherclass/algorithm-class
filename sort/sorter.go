package main

import "fmt"

type sorter interface {
	epithet() string
}

// TODO: 정렬 클래스 분류 관계를 어떻게 표현할 수 있을까?  interface로 표현할
// 수 있고 type code로도 표현할 수 있다. 지금은 정렬 클래스 갯수가 작으니까
// 상관없이 interface를 쓰지만 다루기 어려워지면 type code를 고려해본다.
type isqsort interface {
	isqsort()
}

type sequenceSorter interface {
	sort(s sequence, r lesser, c *sortCounter) source
	sorter
}

type vecSorter interface {
	sort(v vec, r lesser, c *sortCounter) source
	sorter
}

type lnkSorter interface {
	sort(l lnk, r lesser, c *sortCounter) source
	sorter
}

// func callSort(sorter sorter, src source, r lesser, c *sortCounter) source {
// 	switch sorter := sorter.(type) {
// 	case sequenceSorter:
// 		s, ok := src.(sequence)
// 		aver(ok, "source must be of sequence type")
// 		return sorter.sort(s, r, c)
// 	case vecSorter:
// 		v, ok := src.(vec)
// 		aver(ok, "source must be of vec type")
// 		return sorter.sort(v, r, c)
// 	case lnkSorter:
// 		l, ok := src.(lnk)
// 		aver(ok, "source must be of lnk type")
// 		return sorter.sort(l, r, c)
// 	default:
// 		panic("unrecognized sorter")
// 	}
// }

func aver(truth bool, format string, args ...interface{}) {
	if truth {
		return
	}
	panic(fmt.Sprintf(format, args...))
}

func convLink(s []int) lnk {
	l := new(aslnk)
	for _, x := range s {
		l.Push(x)
	}
	return l
}

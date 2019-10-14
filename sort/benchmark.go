package main

import "time"

const (
	rangeMin      = 0
	rangeMax      = 1000
	inputMaxsize  = 400
	sortIteration = 7
	inputBits     = 3
	inputRadix    = 10
)

type benchmarkResult struct {
	sortName string
	records  []benchmarkRecord
}

type benchmarkRecord struct {
	sortName  string
	inputName string // TODO: extends
	samples   []sortStat
}

func benchmark(sorter sorter, maxsize, iteration uint) benchmarkResult {
	swappedSorted := func(swapFactor float64) sizedInputFunc {
		return func(size uint) inputFunc {
			swap := uint(float64(size) * swapFactor)
			return almostSortedInput(size, swap)
		}
	}
	type inputType struct {
		name      string
		makeinput sizedInputFunc
	}
	var inputs = []inputType{
		{"fuzz input", fuzzInput},
		{"sorted input", sortedInput},
		{"reversed input", reversedInput},
		// {"almost sorted input (0.75 swapped)", swappedSorted(0.75)},
		{"almost sorted input (0.1 swapped)", swappedSorted(0.10)},
		// {"almost sorted input (0.25 swapped)", swappedSorted(0.25)},
		// {"almost sorted input (0.125 swapped)", swappedSorted(0.125)},
	}
	if _, ok := sorter.(isqsort); ok {
		killqsort := func(size uint) inputFunc {
			killer := antiqsort(sorter.(sequenceSorter), int(size))
			return constInput(killer)
		}
		inputs = append(inputs, inputType{"killing input", killqsort})
	}
	records := make([]benchmarkRecord, 0, len(inputs))
	for _, input := range inputs {
		record := benchmarkInput(sorter, input.name, input.makeinput, iteration, maxsize)
		records = append(records, record)
	}
	return benchmarkResult{
		sortName: sorter.epithet(),
		records:  records,
	}
}

func benchmarkInput(sorter sorter, inputName string, makeinput sizedInputFunc, iteration, maxsize uint) benchmarkRecord {
	samples := iterateSizedSort(sorter, makeinput, constIteration(iteration), maxsize)
	return benchmarkRecord{
		sortName:  sorter.epithet(),
		inputName: inputName,
		samples:   samples,
	}
}

type inputFunc func(iteration uint) []int

func measureSort(sorter sorter, src []int) (source, lesser, sortCounter) {
	var c sortCounter
	startTime := time.Now()
	r := &aslesser{defaultLesser(), &c}
	sorted := callSort(sorter, src, r, &c)
	c.lapse = time.Since(startTime)
	return sorted, r, c
}

func defaultLesser() lesser {
	return &lesserRadixInt{
		bits:  inputBits,
		radix: inputRadix,
	}
}

func callSort(sorter sorter, src []int, r lesser, c *sortCounter) source {
	switch sorter := sorter.(type) {
	case sequenceSorter:
		return sorter.sort(wrapSequence(asints(src), c), r, c)
	case vecSorter:
		v := asints(src)
		return sorter.sort(wrapVec(&v, c), r, c)
	case lnkSorter:
		return sorter.sort(wrapLink(convLink(src), c), r, c)
	default:
		panic("unrecognized sorter")
	}
}

func iterateSort(sorter sorter, inputFunc inputFunc, iteration uint) sortStat {
	// Design: 함수 구조로 만들어져 있어 iterateSort와 iterateSizedSort는
	// 병렬화가 가능해 보인다.
	var stat sortStat
	for i := uint(1); i <= iteration; i++ {
		input := inputFunc(i)
		_, _, counter := measureSort(sorter, input)
		accCounter(&stat, counter)
	}
	return stat
}

type sizedInputFunc func(size uint) inputFunc

type sizedIterationFunc func(size uint) (iteration uint)

func iterateSizedSort(sorter sorter, sizedInputFunc sizedInputFunc, sizedIterationFunc sizedIterationFunc, maxsize uint) []sortStat {
	buf := make([]sortStat, 0, maxsize+1)
	for size := uint(0); size <= maxsize; size++ {
		inputFunc := sizedInputFunc(size)
		iteration := sizedIterationFunc(size)
		buf = append(buf, iterateSort(sorter, inputFunc, iteration))
	}
	return buf
}

func constInput(s []int) inputFunc {
	t := make([]int, len(s))
	return func(iteration uint) []int {
		copy(t, s)
		return t
	}
}

func fuzzInput(size uint) inputFunc {
	s := make([]int, size)
	return func(iteration uint) []int {
		for i := range s {
			s[i] = rngSource.Intn(int(size))
		}
		return s
	}
}

func sortedInput(size uint) inputFunc {
	s := make([]int, size)
	return func(iteration uint) []int {
		for i := range s {
			s[i] = i
		}
		return s
	}
}

func reversedInput(size uint) inputFunc {
	s := make([]int, size)
	return func(iteration uint) []int {
		n := int(size - 1)
		for i := range s {
			s[i] = n
			n--
		}
		return s
	}
}

func almostSortedInput(size, swap uint) inputFunc {
	s := make([]int, size)
	return func(iteration uint) []int {
		for i := range s {
			s[i] = i
		}
		for i := 0; i < int(swap); i++ {
			i, j := rngSource.Intn(int(size)), rngSource.Intn(int(size))
			s[i], s[j] = s[j], s[i]
		}
		return s
	}
}

func constIteration(iteration uint) sizedIterationFunc {
	return func(size uint) uint {
		return iteration
	}
}

package main

import (
	"fmt"
	"os"
)

type instructionCounter struct {
	n int
}

func (i *instructionCounter) reset()  { i.n = 0 }
func (i *instructionCounter) inc()    { i.n++ }
func (i instructionCounter) now() int { return i.n }

var instCounter instructionCounter

func binarySearch(a []int, key int, left, right int) int {
	if left <= right {
		mid := (left + right) / 2
		instCounter.inc()
		if key == a[mid] {
			return mid
		}
		if key < a[mid] {
			return binarySearch(a, key, left, mid-1)
		}
		return binarySearch(a, key, mid+1, right)
	}
	return -1
}

func inspectBinarySearch(a []int, key int) (index int, insts int) {
	instCounter.reset()
	index = binarySearch(a, key, 0, len(a)-1)
	insts = instCounter.now()
	return
}

func test(a []int, key int, expected int) {
	got, insts := inspectBinarySearch(a, key)
	if got == expected {
		fmt.Printf("OK insts=%d\n", insts)
		return
	}
	fmt.Fprintf(os.Stderr, "FAIL: binarySearch(key=%d) is expected to return %d, but got %d",
		key, expected, got)
}

var a = []int{1, 3, 5, 6, 7, 8, 10, 11, 12, 13, 15}

func main() {
	fmt.Println("배열", a)
	for {
		fmt.Print("키 = ")
		var key int
		fmt.Scan(&key)
		index, insts := inspectBinarySearch(a, key)
		fmt.Printf("위치 = %d, 비교 횟수 = %d\n", index, insts)
	}
}

func testMain() {
	test(a, 3, 1)
	test(a, 5, 2)
	test(a, -123, -1)
	test(a, -12412331, -1)
	test(a, 123819, -1)
	test(a, 13, 9)
	test(a, 42, -1)
	test(a, 8, 5)
}

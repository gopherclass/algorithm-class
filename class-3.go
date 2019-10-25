//+build mage

package main

import (
	"fmt"
)

func isPerfect(n int) int {
	s := 0
	for d := 1; d < n; d++ {
		if n%d == 0 {
			s += d
		}
	}
	if s < n {
		return -1
	}
	if s > n {
		return 1
	}
	return 0
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n%2 == 0 {
		return false
	}
	for i := 3; i < n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func printNum(n int) {
	k := 1
	i := 1
	for i <= n {
		for j := 1; j <= k && i <= n; j++ {
			fmt.Print(i, " ")
			i++
		}
		fmt.Println()
		k *= 2
	}
}

func Problem1() {
	for i := 1; i <= 30; i++ {
		fmt.Println(i, "=>", what(isPerfect(i)))
	}
}

func what(v int) string {
	if v < 0 {
		return "부족수"
	}
	if v > 0 {
		return "과잉수"
	}
	return "완전수"
}

func Problem2() {
	for i := 1; i <= 100; i++ {
		if isPrime(i) {
			fmt.Println(i, "is a prime number")
		} else {
			fmt.Println(i, "is not a prime number")
		}
	}
}

func Problem3() {
	printNum(17)
	fmt.Println()

	printNum(30)
}

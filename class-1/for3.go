package main

import (
	"fmt"
)

func f(n int) int {
	return (n + 2) * (n + 1) * n / 6
}

func g(n int) int {
	x := 0
	for i := 1; i <= n; i++ {
		for j := 1; j <= i; j++ {
			for k := 1; k <= j; k++ {
				x++
			}
		}
	}
	return x
}

func main() {
	for i := 1; i <= 20; i++ {
		fmt.Println(i, f(i), g(i))
	}
}

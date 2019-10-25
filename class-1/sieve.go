package main

import "fmt"

func main() {
	for {
		fmt.Print("n = ")
		var n int
		fmt.Scan(&n)
		r := sieve(array(n))
		fmt.Println(r)
	}
}

func array(n int) []int {
	a := make([]int, n)
	for i := range a {
		a[i] = i
	}
	return a
}

func sieve(a []int) []int {
	n := len(a) - 1
	a[1] = 0
	x := 2
	for x <= n {
		y := 2 * x
		for y <= n {
			a[y] = 0
			y += x
		}
		x++
		for x <= n && a[x] == 0 {
			x++
		}
	}
	return a
}

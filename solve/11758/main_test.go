package main

func ExampleSolve1() {
	solveMain(Example(`1 1
5 5
7 3`))
	// Output: -1
}

func ExampleSolve2() {
	solveMain(Example(`1 1
3 3
5 5`))
	// Output: 0
}

func ExampleSolve3() {
	solveMain(Example(`1 1
7 3
5 5`))
	// Output: 1
}

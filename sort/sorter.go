package main

type sorter interface {
	sort(*sortCounter, []int) []int
	epithet() string
}

// TODO: 정렬 클래스 분류 관계를 어떻게 표현할 수 있을까?  interface로 표현할
// 수 있고 type code로도 표현할 수 있다. 지금은 정렬 클래스 갯수가 작으니까
// 상관없이 interface를 쓰지만 다루기 어려워지면 type code를 고려해본다.
type isqsort interface {
	isqsort()
}

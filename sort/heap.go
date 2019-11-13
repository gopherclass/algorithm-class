package main

type rel interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

type vec interface {
	rel
	Push(interface{})
	Pop() interface{}
}

func heapinit(r rel) {
	n := r.Len()
	for i := n / 2; i >= 0; i-- {
		heapdrop(r, i)
	}
}

func heapdrop(r rel, i int) {
	n := r.Len()
	for i < n {
		j := 2*i + 1
		if j >= n {
			break
		}
		if j+1 < n && r.Less(j+1, j) {
			j = j + 1
		}
		if r.Less(i, j) {
			break
		}
		r.Swap(i, j)
		i = j
	}
}

func heapfloat(r rel, i int) {
	for i > 0 {
		j := (i - 1) / 2
		if r.Less(i, j) {
			break
		}
		r.Swap(i, j)
		i = j
	}
}

func heappop(v vec) interface{} {
	n := v.Len()
	v.Swap(0, n-1)
	x := v.Pop()
	heapdrop(v, 0)
	return x
}

func heappush(v vec, x interface{}) {
	v.Push(x)
	heapfloat(v, v.Len()-1)
}

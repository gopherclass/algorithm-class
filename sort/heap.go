package main

func heapinit(s vecless) {
	n := s.Len()
	for i := n / 2; i >= 0; i-- {
		heapdrop(s, i)
	}
}

func heapdrop(s vecless, i int) {
	n := s.Len()
	for i < n {
		j := 2*i + 1
		if j >= n {
			break
		}
		if j+1 < n && s.Less(j+1, j) {
			j = j + 1
		}
		if s.Less(i, j) {
			break
		}
		s.Swap(i, j)
		i = j
	}
}

func heapfloat(s vecless, i int) {
	for i > 0 {
		j := (i - 1) / 2
		if s.Less(i, j) {
			break
		}
		s.Swap(i, j)
		i = j
	}
}

func heappop(v vecless) interface{} {
	n := v.Len()
	v.Swap(0, n-1)
	x := v.Pop()
	heapdrop(v, 0)
	return x
}

func heappush(v vecless, x interface{}) {
	v.Push(x)
	heapfloat(v, v.Len()-1)
}

package optimalBinaryTree

type Prior struct {
	Probabilities []float64
	Costs         squareFloat
	Decisions     squareInt
}

func (r *Prior) Min() float64 {
	return r.Costs.At(r.Costs.stride-1, 0)
}

type Tree struct {
	root *Node
	size int
}

func (t *Tree) Size() int {
	return t.size
}

func (t *Tree) Search(value int) int {
	return binarySearch(t.root, value, 0)
}

func binarySearch(x *Node, value int, cost int) int {
	if x == nil {
		return -1
	}
	if value < x.Value {
		return binarySearch(x.Left, value, cost+1)
	}
	if value > x.Value {
		return binarySearch(x.Right, value, cost+1)
	}
	return cost + 1
}

type Node struct {
	Value       int
	Probability float64
	Left, Right *Node
}

func (r *Prior) Tree() *Tree {
	size := len(r.Probabilities)
	root := r.buildNode(0, size-1)
	return &Tree{
		root: root,
		size: size,
	}
}

func (r *Prior) buildNode(i, j int) *Node {
	if i > j {
		return nil
	}
	split := r.Decisions.At(j, i)
	return &Node{
		Value:       split,
		Probability: r.Probabilities[split],
		Left:        r.buildNode(i, split-1),
		Right:       r.buildNode(split+1, j),
	}
}

func NewPrior(s []float64) *Prior {
	costs := newSquareFloat(len(s), len(s))
	decs := newSquareInt(len(s), len(s))
	sigma := consecutiveSum(s)
	for i := 0; i < len(s); i++ {
		costs.Set(i, i, s[i])
		decs.Set(i, i, i)
	}
	for i := 1; i < len(s); i++ {
		for x, y := i, 0; x < len(s) && y < len(s); x, y = x+1, y+1 {
			// y .. x
			zMin := -1
			for z := y; z <= x; z++ {
				var cost float64
				if z >= 1 {
					cost += costs.At(z-1, y)
				}
				if z+1 < len(s) {
					cost += costs.At(x, z+1)
				}
				if costs.SetMin(x, y, cost) {
					zMin = z
				}
			}
			costs.Add(x, y, sigma.At(x, y))
			decs.Set(x, y, zMin)
		}
	}
	return &Prior{
		Probabilities: s,
		Costs:         costs,
		Decisions:     decs,
	}
}

func NewPrior2(s []float64) *Prior {
	costs := newSquareFloat(len(s), len(s))
	decs := newSquareInt(len(s), len(s))
	sigma := consecutiveSum(s)
	for i := 0; i < len(s); i++ {
		costs.Set(i, i, s[i])
		decs.Set(i, i, i)
	}
	for i := 1; i < len(s); i++ {
		for x, y := i, 0; x < len(s) && y < len(s); x, y = x+1, y+1 {
			// y .. x
			zMin := -1
			for z := y; z <= x; z++ {
				cost := s[z]
				if z >= 1 {
					A := sigma.At(z-1, y)
					cost += A * costs.At(z-1, y)
					cost += A * A
				}
				if z+1 < len(s) {
					B := sigma.At(x, z+1)
					cost += B * costs.At(x, z+1)
					cost += B * B
				}
				costs.Add(x, y, sigma.At(x, y))
				if costs.SetMin(x, y, cost) {
					zMin = z
				}
			}
			decs.Set(x, y, zMin)
		}
	}
	return &Prior{
		Probabilities: s,
		Costs:         costs,
		Decisions:     decs,
	}
}

func consecutiveSum(s []float64) squareFloat {
	sigma := newSquareFloat(len(s), len(s))
	for y := 0; y < len(s); y++ {
		sigma.Set(y, y, s[y])
		for x := y + 1; x < len(s); x++ {
			sigma.Set(x, y, sigma.At(x-1, y)+s[x])
		}
	}
	return sigma
}

type squareInt struct {
	s      []int
	stride int
}

func (s squareInt) Offset(x, y int) int {
	i := y*s.stride + x
	return i
}

func (s squareInt) At(x, y int) int {
	return s.s[s.Offset(x, y)]
}

func (s squareInt) Set(x, y, v int) {
	s.s[s.Offset(x, y)] = v
}

func newSquareInt(row, stride int) squareInt {
	return squareInt{
		s:      make([]int, row*stride),
		stride: stride,
	}
}

type squareFloat struct {
	s      []float64
	stride int
}

func newSquareFloat(row, stride int) squareFloat {
	return squareFloat{
		s:      make([]float64, row*stride),
		stride: stride,
	}
}

func (s squareFloat) Offset(x, y int) int {
	i := y*s.stride + x
	return i
}

func (s squareFloat) At(x, y int) float64 {
	return s.s[s.Offset(x, y)]
}

func (s squareFloat) Set(x, y int, v float64) {
	s.s[s.Offset(x, y)] = v
}

func (s squareFloat) Add(x, y int, v float64) {
	s.s[s.Offset(x, y)] += v
}

func (s squareFloat) SetMin(x, y int, v float64) bool {
	i := s.Offset(x, y)
	if v < s.s[i] || s.s[i] <= 0 {
		s.s[i] = v
		return true
	}
	return false
}

package plane

type Point struct {
	X, Y int
}

func Pt(x, y int) Point {
	return Point{x, y}
}

func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

func (p Point) Mul(z int) Point {
	return Point{z * p.X, z * p.Y}
}

func (p Point) Dot(q Point) int {
	return p.X*q.X + p.Y*q.Y
}

func Det(p, q Point) int {
	return p.X*q.Y - p.Y*q.X
}

func Dir(x, y, z Point) int {
	return Det(y.Sub(x), z.Sub(x))
}

package packageWrapping

type Point struct {
	X, Y int
}

func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func atan2(a, b Point) float64 {
	return theta(b.Sub(a))
}

// TODO: Tested
func theta(p Point) float64 {
	if p.X == 0 && p.Y == 0 {
		return 0
	}
	t := p.Y / (abs(p.X) + abs(p.Y))
	if p.X < 0 {
		return 2 - t
	}
	if p.Y < 0 {
		return 4 + t
	}
	return t
}

// TODO: Tested
func Wrap(set []Point) []Point {
	if len(set) == 0 {
		return nil
	}
	i := lowestY(set)
	for i := 0; i < len(set); i++ {
		for j := 0; j < len(set); j++ {
			if atan2(x, 
		}
	}
}

func lowestY(set []Point) int {
	if len(set) == 0 {
		return -1
	}
	i, p := 0, set[0]
	for j, q := range set[1:] {
		if q.Y < p.Y {
			i, p = j, q
		}
	}
	return i
}


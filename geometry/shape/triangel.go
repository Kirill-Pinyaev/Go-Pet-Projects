package shape

import (
	"math"
)

type Triangle struct {
	side1, side2, side3 float64
	height              float64
}

func (t Triangle) Area() float64 {
	if t.height > 0 {
		return 0.5 * t.side1 * t.height
	}
	p := t.side1 + t.side2 + t.side3
	s := p / 2
	area := math.Sqrt(s * (s - t.side1) * (s - t.side2) * (s - t.side3))
	return area
}

func (t Triangle) Perimeter() float64 {
	if t.height > 0 {
		return (t.side1 + t.side2 + math.Hypot(t.side2, t.side1))
	}
	return t.side1 + t.side2 + t.side3
}

package shape

type Rectangle struct {
	side1 float64
	side2 float64
}

func (r Rectangle) Area() float64 {
	return r.side1 * r.side2
}

func (r Rectangle) Perimeter() float64 {
	return (r.side1 * 2) + (r.side2 * 2)
}

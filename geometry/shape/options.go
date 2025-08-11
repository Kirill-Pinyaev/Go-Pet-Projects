package shape

import (
	"errors"
	"strings"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type baseShape struct {
	color string
}

type Option func(*baseShape) error

func WithColor(c string) Option {
	return func(b *baseShape) error {
		if strings.TrimSpace(c) == "" {
			return errors.New("color: empty string")
		}
		b.color = c
		return nil
	}
}

func NewShape(kind string, params []float64, opt ...Option) (Shape, error) {
	if kind = strings.ToLower(strings.TrimSpace(kind)); kind == "" {
		return nil, errors.New("kind: empty")
	}
	var s Shape
	switch kind {
	case "circle":
		if len(params) != 1 {
			return nil, errors.New("circle: need 1 param (radius)")
		}
		r := params[0]
		if r <= 0 {
			return nil, errors.New("circle: radius must be > 0")
		}
		s = Circle{radius: r}
	case "rectangle":
		if len(params) != 2 {
			return nil, errors.New("rectangle: need 2 params (a, b)")
		}
		a, b := params[0], params[1]
		if a <= 0 || b <= 0 {
			return nil, errors.New("rectangle: sides must be > 0")
		}
		s = Rectangle{side1: a, side2: b}
	case "triangle":
		switch len(params) {
		case 2:
			base, h := params[0], params[1]
			if base <= 0 || h <= 0 {
				return nil, errors.New("triangle: base/height must be > 0")
			}
			s = &Triangle{side1: base, height: h}
		case 3:
			a, b, c := params[0], params[1], params[2]
			if a <= 0 || b <= 0 || c <= 0 ||
				a+b <= c || a+c <= b || b+c <= a {
				return nil, errors.New("triangle: invalid sides")
			}
			s = &Triangle{side1: a, side2: b, side3: c}
		default:
			return nil, errors.New("triangle: need 2 or 3 params")
		}
	}
	// применение опций
	return s, nil
}

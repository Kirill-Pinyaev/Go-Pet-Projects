package main

import (
	"flag"
	"fmt"
	"geometry/shape"
	"log"
)

var (
	kind  = flag.String("shape", "", "circle | rectangle | triangle")
	r     = flag.Float64("r", 0, "radius for circle")
	a     = flag.Float64("a", 0, "side a / width / base")
	b     = flag.Float64("b", 0, "side b / height")
	c     = flag.Float64("c", 0, "side c (triangle)")
	prec  = flag.Int("prec", 3, "digits after decimal point")
	color = flag.String("color", "", "optional color tag")
)

func main() {
	flag.Parse()
	var params []float64
	switch *kind {
	case "circle":
		params = []float64{*r}
	case "rectangle":
		params = []float64{*a, *b}
	case "triangle":
		if *c != 0 {
			params = []float64{*a, *b, *c}
		} else {
			params = []float64{*a, *b}
		}
	default:
		log.Fatalf("неизвестная фигура: %s", *kind)
	}

	fig, err := shape.NewShape(
		*kind,
		params,
		shape.WithColor(*color),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%.*f %.*f\n", *prec, fig.Area(), *prec, fig.Perimeter())
}

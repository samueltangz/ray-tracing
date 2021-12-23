package main

type Color struct {
	r float64
	g float64
	b float64
}

func NewColor(r, g, b float64) Color {
	return Color{r, g, b}
}

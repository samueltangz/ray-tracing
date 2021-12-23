package main

import (
	"fmt"
)

func main () {
	x := NewVector(1, 2, 3)
	y := NewVector(4, 5, 6)

	z := Add(x, y)

	fmt.Println(z, z.NormSquared())
}
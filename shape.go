package main

type Shape interface {
	// Intersect(u, du) finds a non-negative x such that u + x*du lies on the shape 
	Intersect(Vector, Vector)

}

type Sphere struct {
	center Vector
	radius float64
}


package main

type Ray struct {
	origin    Vector
	direction Vector
}

func (r Ray) At(t float64) Vector {
	return Add(r.origin, ScalarMul(t, r.direction))
}

func NewRay(origin, direction Vector) Ray {
	return Ray{
		origin:    origin,
		direction: direction,
	}
}
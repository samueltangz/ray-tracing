package main

import (
	"math"
)

func SolveQuadratic(a, b, c float64) []float64 {
	discriminant := math.Pow(b, 2) - 4*a*c
	if discriminant < 0 {
		return []float64{}
	}
	if discriminant == 0 {
		return []float64{
			-b / (2 * a),
		}
	}

	return []float64{
		(-b - math.Sqrt(discriminant)) / (2 * a),
		(-b + math.Sqrt(discriminant)) / (2 * a),
	}
}

// Transform converts a vector v into the given basis
func Transform(v Vector, basis [3]Vector) Vector {
	return NewVector(
		Dot(v, basis[0]),
		Dot(v, basis[1]),
		Dot(v, basis[2]),
	)
}

// InverseTransform converts a vector v from the given basis
func InverseTransform(v Vector, basis [3]Vector) Vector {
	return NewVector(
		Dot(v, NewVector(basis[0].x, basis[1].x, basis[2].x)),
		Dot(v, NewVector(basis[0].y, basis[1].y, basis[2].y)),
		Dot(v, NewVector(basis[0].z, basis[1].z, basis[2].z)),
	)
}

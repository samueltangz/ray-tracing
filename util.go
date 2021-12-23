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

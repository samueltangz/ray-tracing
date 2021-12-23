package main

import (
	"math"
)

const EPSILON float64 = 1e-8

type Shape interface {
	// Hit(r) finds a non-negative x such that r.origin + x*r.direction lies on the shape
	Hit(Ray) (bool, float64)

	// UnitNormal(v) returns a unit normal vector against the current point
	UnitNormal(Vector) Vector
}

// Sphere

type Sphere struct {
	center Vector
	radius float64
}

func (s Sphere) Hit(r Ray) (bool, float64) {
	oc := Sub(r.origin, s.center)

	a := r.direction.NormSquared()
	b := 2.0 * Dot(oc, r.direction)
	c := oc.NormSquared() - math.Pow(s.radius, 2)

	roots := SolveQuadratic(a, b, c)
	if len(roots) > 0 && roots[0] > EPSILON {
		return true, roots[0]
	} else if len(roots) > 1 && roots[1] > EPSILON {
		return true, roots[1]
	}
	return false, 0.0
}

func (s Sphere) UnitNormal(v Vector) Vector {
	return Sub(v, s.center).Unit()
}

// Triangle

type Triangle struct {
	v0  Vector
	dv1 Vector
	dv2 Vector
}

func NewTriangle(v0, v1, v2 Vector) Triangle {
	// BUG: if v0 and v1 is swapped then the thing would not be correctly rendered
	dv1 := Sub(v1, v0)
	dv2 := Sub(v2, v0)
	return Triangle{v0, dv1, dv2}
}

func (s Triangle) Hit(r Ray) (bool, float64) {
	ft := Dot(s.dv1, Cross(r.direction, s.dv2))

	dp1 := Dot(Sub(r.origin, s.v0), Cross(r.direction, s.dv2)) / ft
	dp2 := Dot(r.direction, Cross(Sub(r.origin, s.v0), s.dv1)) / ft
	dp3 := Dot(s.dv2, Cross(Sub(r.origin, s.v0), s.dv1)) / ft

	if dp1 < 0 || dp2 < 0 || dp1 > 1 || (dp1+dp2) > 1 {
		return false, 0.0
	}
	if dp3 <= EPSILON {
		return false, 0.0
	}
	return true, dp3
}

func (s Triangle) UnitNormal(v Vector) Vector {
	return Cross(s.dv1, s.dv2).Unit()
}

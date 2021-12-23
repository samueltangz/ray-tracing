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

// Ellipsoid

type Ellipsoid struct {
	center Vector
	basis  [3]Vector
	l2     Vector
}

func NewEllipsoid(center, u1, u2 Vector, lengths [3]float64) Ellipsoid {
	return Ellipsoid{
		center: center,
		basis: [3]Vector{
			u1.Unit(),
			u2.Unit(),
			Cross(u1, u2).Unit(),
		},
		l2: NewVector(
			math.Pow(lengths[0], 2),
			math.Pow(lengths[1], 2),
			math.Pow(lengths[2], 2),
		),
	}
}

func (s Ellipsoid) Hit(r Ray) (bool, float64) {
	oc := Sub(r.origin, s.center)

	a := ReduceSum(Div(Mul(r.direction, r.direction), s.l2))
	b := ReduceSum(Div(ScalarMul(2.0, Mul(r.direction, oc)), s.l2))
	c := ReduceSum(Div(Mul(oc, oc), s.l2)) - 1

	roots := SolveQuadratic(a, b, c)
	if len(roots) > 0 && roots[0] > EPSILON {
		return true, roots[0]
	} else if len(roots) > 1 && roots[1] > EPSILON {
		return true, roots[1]
	}
	return false, 0.0
}

func (s Ellipsoid) UnitNormal(v Vector) Vector {
	return InverseTransform(Div(Transform(Sub(v, s.center), s.basis), s.l2), s.basis).Unit()
}

// Unknown Shape (From hxp ctf)

type Shape1 struct {
	x Vector
	r float64
}

func (s Shape1) blood(y Vector, _s float64) (bool, float64) {
	u := math.Pow(s.r, 2) + math.Pow(_s, 2) - y.NormSquared()
	if u < 0 {
		return false, 0.0
	}
	v := _s - math.Sqrt(u)
	if v < 0 {
		return false, 0.0
	}
	return true, v
}

func (s Shape1) Hit(r Ray) (bool, float64) {
	q := Sub(s.x, r.origin)
	return s.blood(q, Dot(q, r.direction))
}

func (s Shape1) UnitNormal(v Vector) Vector {
	return Sub(v, s.x).Unit()
}

// Unknown Shape 2 (From hxp ctf)

type Shape2 struct {
	x Vector
	y float64
}

func NewShape2(x Vector, y float64) Shape2 {
	return Shape2{x, y}
}

func (s Shape2) Hit(r Ray) (bool, float64) {
	y1 := r.origin
	y2 := Dot(s.x.Unit(), r.direction)

	if y2 >= -EPSILON {
		return false, 0.0
	}
	q := (Dot(s.x.Unit(), y1) + s.y) / -y2
	if q < 0 {
		return false, 0.0
	}
	return true, q
}

func (s Shape2) UnitNormal(v Vector) Vector {
	return v.Unit()
}

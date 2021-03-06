package main

import (
	"math"
	"math/rand"
)

type Vector struct {
	x float64
	y float64
	z float64
}

func (v Vector) NormSquared() float64 {
	return Dot(v, v)
}

func (v Vector) Norm() float64 {
	return math.Sqrt(v.NormSquared())
}

func (v Vector) Neg() Vector {
	return NewVector(-v.x, -v.y, -v.z)
}

func (v Vector) Unit() Vector {
	return ScalarMul(1/v.Norm(), v)
}

func NewVector(x, y, z float64) Vector {
	return Vector{
		x: x,
		y: y,
		z: z,
	}
}

func Add(v1 Vector, v2 Vector) Vector {
	return NewVector(v1.x+v2.x, v1.y+v2.y, v1.z+v2.z)
}

func Sub(v1 Vector, v2 Vector) Vector {
	return NewVector(v1.x-v2.x, v1.y-v2.y, v1.z-v2.z)
}

func Mul(v1 Vector, v2 Vector) Vector {
	return NewVector(v1.x*v2.x, v1.y*v2.y, v1.z*v2.z)
}

func Div(v1 Vector, v2 Vector) Vector {
	return NewVector(v1.x/v2.x, v1.y/v2.y, v1.z/v2.z)
}

func ScalarMul(t float64, v Vector) Vector {
	return NewVector(t*v.x, t*v.y, t*v.z)
}

func Dot(v1 Vector, v2 Vector) float64 {
	return v1.x*v2.x + v1.y*v2.y + v1.z*v2.z
}

func Cross(v1 Vector, v2 Vector) Vector {
	return NewVector(
		v1.y*v2.z-v1.z*v2.y,
		v1.z*v2.x-v1.x*v2.z,
		v1.x*v2.y-v1.y*v2.x,
	)
}

func Equal(v1 Vector, v2 Vector) bool {
	// TODO: Check if there is needs to accept epsilon as error
	return v1.x == v2.x && v1.y == v2.y && v1.z == v2.z
}

func RandomVector(n float64) Vector {
	for {
		v := NewVector(
			2*rand.Float64()-1,
			2*rand.Float64()-1,
			2*rand.Float64()-1,
		)
		if v.NormSquared() > 1 {
			continue
		}
		return ScalarMul(n, v)
	}
}

// Reflect creates a vector v' that is reflected against the normal v
func Reflect(v Vector, n Vector) Vector {
	return Sub(v, ScalarMul(2*Dot(v, n), n))
}

func ReduceSum(v Vector) float64 {
	return v.x + v.y + v.z
}

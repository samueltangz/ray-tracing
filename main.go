package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/schollz/progressbar/v3"
)

type Screen struct {
	origin          Vector
	lowerLeftCorner Vector
	horizontal      Vector
	vertical        Vector
	width           int64
	height          int64
	objects         []Object
}

func NewScreen(lookFrom, lookAt, vup Vector, viewportWidth, viewportHeight float64, width, height int64, objects []Object) Screen {
	w := Sub(lookFrom, lookAt).Unit()
	u := Cross(vup, w)
	v := Cross(w, u)

	origin := lookFrom
	horizontal := ScalarMul(viewportWidth, u)
	vertical := ScalarMul(viewportHeight, v)

	lowerLeftCorner := Sub(origin, ScalarMul(0.5, horizontal))
	lowerLeftCorner  = Sub(lowerLeftCorner, ScalarMul(0.5, vertical))
	lowerLeftCorner  = Sub(lowerLeftCorner, w)

	return Screen{
		origin:          origin,
		lowerLeftCorner: lowerLeftCorner,
		horizontal:      horizontal,
		vertical:        vertical,
		width:           width,
		height:          height,
		objects:         objects,
	}
}

func (s Screen) hit(r Ray) (bool, Object, float64) {
	minX := math.Inf(1)
	minXAt := Object{}
	hit := false

	for _, object := range s.objects {
		ok, x := object.shape.Hit(r)
		if !ok {
			continue
		}
		if x > minX {
			continue
		}

		hit = true
		minX = x
		minXAt = object
	}

	return hit, minXAt, minX
}

func (s Screen) color(r Ray, depth int) [3]float64 {
	// too insignificant... let's skip it
	if depth >= 64 {
		return [3]float64{0, 0, 0}
	}

	ok, object, x := s.hit(r)
	if ok {
		v := r.At(x)
		n := object.shape.UnitNormal(v)
		n = object.material.Perbute(n)

		baseColor := object.color
		reflectColor := s.color(NewRay(v, n), depth+1)

		return [3]float64{
			baseColor.r * reflectColor[0],
			baseColor.g * reflectColor[1],
			baseColor.b * reflectColor[2],
		}
	}

	// Not hitting anything - show them the background.

	t := 0.5 * (r.direction.Unit().y + 1.0)
	return [3]float64{
		1.0 - 0.5*t,
		1.0 - 0.3*t,
		1.0,
	}
}

func clip(u float64) float64 {
	if u < 0 {
		return 0.0
	}
	if u > 1 {
		return 1.0
	}
	return math.Sqrt(u)
}

func RGB(u [3]float64) []int {
	return []int{
		int(255 * clip(u[0])),
		int(255 * clip(u[1])),
		int(255 * clip(u[2])),
	}
}

func (s *Screen) Render(filename string, antiAliasingFactor int) error {
	bar := progressbar.Default(s.width * s.height)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	fmt.Fprintf(f, "P3\n")
	fmt.Fprintf(f, "%d %d\n", s.width, s.height)
	fmt.Fprintf(f, "255\n")

	for y := s.height - 1; y >= 0; y-- {
		for x := int64(0); x < s.width; x++ {
			c := [3]float64{0.0, 0.0, 0.0}
			for i := 0; i < antiAliasingFactor; i++ {
				dx := rand.Float64() * 0.5
				dy := rand.Float64() * 0.5

				r := NewRay(
					s.origin,
					Sub(
						Add(
							s.lowerLeftCorner,
							Add(
								ScalarMul((float64(x)+dx)/float64(s.width), s.horizontal),
								ScalarMul((float64(y)+dy)/float64(s.height), s.vertical),
							),
						),
						s.origin,
					).Unit(),
				)
				dc := s.color(r, 0)
				for j := 0; j < 3; j++ {
					c[j] += dc[j] / float64(antiAliasingFactor)
				}
			}
			rgb := RGB(c)
			fmt.Fprintf(f, "%d %d %d\n", rgb[0], rgb[1], rgb[2])
			bar.Add(1)
		}
	}

	return nil
}

func main() {
	objects := []Object{
		Object{
			// r
			Sphere{NewVector(2.0, 3.7, 0.0), 1},
			Metal{},
			NewColor(1.0, 0.5, 0.5),
		},

		Object{
			// g
			Sphere{NewVector(0.0, 1.9, 1.0), 1},
			Metal{},
			NewColor(0.5, 1.0, 0.5),
		},

		Object{
			// b
			Sphere{NewVector(-2.0, 5.0, 0.0), 1},
			Metal{},
			NewColor(0.5, 0.5, 1.0),
		},

		Object{
			// y = 0
			Plane{NewVector(0.0, 1.0, 0.0), 0.0},
			Lambertian{},
			NewColor(0.8, 0.8, 0.8),
		},

		// A pair of blue triangles on the ground (pointing towards +x-axis)
		Object{
			NewTriangle(
				NewVector(3.0, 0.01, +0.0),
				NewVector(2.0, 0.01, -0.5),
				NewVector(2.0, 0.01, +0.5),
			),
			Metal{},
			NewColor(0, 0, 1),
		},
		Object{
			NewTriangle(
				NewVector(4.0, 0.01, +0.0),
				NewVector(3.0, 0.01, -0.5),
				NewVector(3.0, 0.01, +0.5),
			),
			Metal{},
			NewColor(0, 0, 1),
		},

		// A pair of green triangles on the ground (pointing towards +z-axis)
		Object{
			NewTriangle(
				NewVector(-0.5, 0.01, 2.0),
				NewVector(+0.0, 0.01, 3.0),
				NewVector(+0.5, 0.01, 2.0),
			),
			Metal{},
			NewColor(0, 0.5, 0),
		},
		Object{
			NewTriangle(
				NewVector(-0.5, 0.01, 3.0),
				NewVector(+0.0, 0.01, 4.0),
				NewVector(+0.5, 0.01, 3.0),
			),
			Metal{},
			NewColor(0, 0.5, 0),
		},
	}

	lookFrom := NewVector(2.5, 1.5, 2.5)
	lookAt := NewVector(2.0, 1.2, 2.0)
	vup := NewVector(0, 1, 0)

	screen := NewScreen(
		lookFrom,
		lookAt,
		vup,
		6.0,
		6.0,
		512,
		512,
		objects,
	)

	screen.Render(
		fmt.Sprintf("rendered.ppm"),
		256,
	)
}

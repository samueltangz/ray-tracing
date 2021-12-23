package main

import (
	"os"
	"fmt"
	"math"
	"math/rand"

	"github.com/schollz/progressbar/v3"
)

type Screen struct {
	origin          Vector

	lowerLeftCorner Vector
	horizontal		Vector
	vertical		Vector

	width           int64
	height          int64

	objects			[]Object
}

func NewScreen(origin, lowerLeftCorner, horizontal, vertical Vector, width, height int64, objects []Object) Screen {
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
		if !ok { continue }
		if x > minX { continue }
		
		hit = true
		minX = x
		minXAt = object
	}

	return hit, minXAt, minX
}

func (s Screen) color(r Ray, depth int) []float64 {
	// too insignificant... let's skip it
	if depth >= 16 {
		return []float64{0, 0, 0}
	}

	ok, object, x := s.hit(r)
	if ok {
		v := r.At(x)
		n := object.shape.UnitNormal(v)
		n = object.material.Perbute(n)
		
		baseColor := object.color
		reflectColor := s.color(NewRay(v, n), depth+1)

		return []float64{
			baseColor.r * reflectColor[0],
			baseColor.g * reflectColor[1],
			baseColor.b * reflectColor[2],
		}
	}

	// Not hitting anything - show them the background.

	t := 0.5 * (r.direction.Unit().y + 1.0)
	return []float64{
		1.0 - 0.5*t,
		1.0 - 0.3*t,
		1.0,
	}
}

func clip(u float64) float64 {
	if u < 0 { return 0.0 }
	if u > 1 { return 1.0 }
	return math.Sqrt(u)
}

func RGB(u []float64) []int {
	return []int{
		int(255 * clip(u[0])),
		int(255 * clip(u[1])),
		int(255 * clip(u[2])),
	}
}

func (s *Screen) Render(filename string) error {
	bar := progressbar.Default(s.width * s.height)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	fmt.Fprintf(f, "P3\n")
	fmt.Fprintf(f, "%d %d\n", s.width, s.height)
	fmt.Fprintf(f, "255\n")

	antiAliasingFactor := 1024

	for y := s.height-1; y >= 0; y-- {
		for x := int64(0); x < s.width; x++ {
			c := []float64{0.0, 0.0, 0.0}
			for i := 0; i < antiAliasingFactor; i++ {
				// TODO: anti-aliasing (scale up the above factor, too)
				dx := rand.Float64() * 0.5
				dy := rand.Float64() * 0.5
				
				r := NewRay(
					s.origin,
					Add(
						s.lowerLeftCorner,
						Add(
							ScalarMul((float64(x) + dx)/float64(s.width), s.horizontal),
							ScalarMul((float64(y) + dy)/float64(s.height), s.vertical),
						),
					),
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
		Object{ // center triangle
			NewTriangle(
				NewVector(-1.0, 0.0, -2.0),
				NewVector(1.0, 0.0,  -2.0),
				NewVector(0.0, 2.0,  -1.0),
			),
			Metal{},
			NewColor(0.3986, 0.3790, 0.3368),
		},
		Object{ // right triangle
			NewTriangle(
				NewVector(1.0, 0.0,  -2.0),
				NewVector(2.0, 0.0,  -1.0),
				NewVector(0.0, 2.0,  -1.0),
			),
			Metal{},
			NewColor(0.3986, 0.3790, 0.3368),
		},
		Object{ // left triangle
			NewTriangle(
				NewVector(-2.0, 0.0, -1.0),
				NewVector(-1.0, 0.0, -2.0),
				NewVector(0.0, 2.0,  -1.0),
			),
			Metal{},
			NewColor(0.3986, 0.3790, 0.3368),
		},
		Object{Sphere{NewVector(0.0, 0.0, -1.0),    0.5  }, Metal{},      NewColor(0.3986, 0.3790, 0.3368)},
		Object{Sphere{NewVector(0.0, -100.5, -1.0), 100.0}, Lambertian{}, NewColor(0.5320, 0.3014, 0.1507)},

		Object{ // center triangle
			NewTriangle(
				NewVector(1.0, 0.0,  2.0),
				NewVector(-1.0, 0.0, 2.0),
				NewVector(0.0, 2.0,  1.0),
			),
			Metal{},
			NewColor(0.3986, 0.3790, 0.3368),
		},
		Object{ // right triangle
			NewTriangle(
				NewVector(2.0, 0.0,  1.0),
				NewVector(1.0, 0.0,  2.0),
				NewVector(0.0, 2.0,  1.0),
			),
			Metal{},
			NewColor(0.3986, 0.3790, 0.3368),
		},
		Object{ // left triangle
			NewTriangle(
				NewVector(-1.0, 0.0, 2.0),
				NewVector(-2.0, 0.0, 1.0),
				NewVector(0.0, 2.0,  1.0),
			),
			Metal{},
			NewColor(0.3986, 0.3790, 0.3368),
		},
	}


	screen := NewScreen(
		NewVector(0.0, 0.3, 0.0),
		NewVector(-2.0, -2.0, -1.0),
		NewVector(4.0, 0.0, 0.0),
		NewVector(0.0, 4.0, 0.0),
		512,
		512,
		objects,
	)

	fmt.Println(screen.Render("render.ppm"))
}
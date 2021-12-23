package main

type Material interface {
	Perbute(Vector) Vector
}

// Lambertian

type Lambertian struct {}

func (m Lambertian) Perbute(v Vector) Vector {
	p := RandomVector(0.5)
	return Add(v, p)
}	

// Metal

type Metal struct {}

func (m Metal) Perbute(v Vector) Vector {
	return v
}	
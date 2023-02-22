package main

import "math/rand"

// Rect represents a rectangle.
type Rect struct {
	X0, Y0, X1, Y1 int
}

// Dx returns the difference between r's X coordinates.
func (r Rect) Dx() int {
	return r.X1 - r.X0
}

// Dy returns the difference between r's Y coordinates.
func (r Rect) Dy() int {
	return r.Y1 - r.Y0
}

// In checks if every point of s is a point of r.
func (r Rect) In(s Rect) bool {
	return r.X0 >= s.X0 && r.X1 <= s.X1 &&
		r.Y0 >= s.Y0 && r.Y1 <= s.Y1
}

// Apply applies f for each point in r.
func (r Rect) Apply(f func(p XY)) {
	for x := r.X0; x < r.X1; x++ {
		for y := r.Y0; y < r.Y1; y++ {
			f(XY{x, y})
		}
	}
}

// OddPoint returns a random point from r with odd x and y
// coordinates.
func (r Rect) OddPoint() XY {
	return XY{
		r.X0 + rand.Intn(r.Dx()/2)*2 + 1,
		r.Y0 + rand.Intn(r.Dy()/2)*2 + 1,
	}
}

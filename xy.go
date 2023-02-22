package main

// XY represents a position on the grid.
type XY struct {
	X int
	Y int
}

// Pos implements the Poser interface.
func (p XY) Pos() XY {
	return p
}

// Add adds all terms to p.
func (p XY) Add(terms ...XY) XY {
	for _, t := range terms {
		p.X += t.X
		p.Y += t.Y
	}
	return p
}

// Sub subtracts all terms from p.
func (p XY) Sub(terms ...XY) XY {
	for _, t := range terms {
		p.X -= t.X
		p.Y -= t.Y
	}
	return p
}

// Mul multiplies p by all factors.
func (p XY) Mul(factors ...int) XY {
	for _, f := range factors {
		p.X *= f
		p.Y *= f
	}
	return p
}

// In reports whether p is in the given bounds.
func (p XY) In(r Rect) bool {
	return r.X0 <= p.X && p.X < r.X1 &&
		r.Y0 <= p.Y && p.Y < r.Y1
}

// Neighbors returns all neighbors of p.
func (p XY) Neighbors() []XY {
	return []XY{p.N(), p.S(), p.W(), p.E(),
		p.N().W(), p.N().E(), p.S().W(), p.S().E()}
}

// Orthogonal returns all neighbors orthogonal to p.
func (p XY) Orthogonal() []XY {
	return []XY{p.N(), p.S(), p.W(), p.E()}
}

// Direction offsets.
var (
	North = XY{Y: -1}
	South = XY{Y: +1}
	West  = XY{X: -1}
	East  = XY{X: +1}
)

func (p XY) N() XY { return p.Add(North) }
func (p XY) S() XY { return p.Add(South) }
func (p XY) W() XY { return p.Add(West) }
func (p XY) E() XY { return p.Add(East) }

// FOV calculates the Field of View of p given the radius r.
// opaque returns true if its argument cannot pass light.
func (p XY) FOV(r int, opaque func(XY) bool) map[XY]bool {
	points := map[XY]bool{}
	for i := -r; i <= r; i++ {
		for j := -r; j <= r; j++ {
			if i*i+j*j < r*r {
				for _, q := range p.Line(p.Add(XY{i, j})) {
					points[q] = true
					if opaque(q) {
						break
					}
				}
			}
		}
	}
	return points
}

// Line returns all points of the line from p to q.
// Implemented using Bresenham's Line Algorithm.
func (p XY) Line(q XY) []XY {
	x1, y1 := p.X, p.Y
	x2, y2 := q.X, q.Y

	steep := abs(y2-y1) > abs(x2-x1)
	if steep {
		x1, y1 = y1, x1
		x2, y2 = y2, x2
	}

	reversed := false
	if x1 > x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
		reversed = true
	}

	dx := x2 - x1
	dy := y2 - y1
	error := dx / 2
	ystep := 1
	if y1 >= y2 {
		ystep = -1
	}

	points := []XY{}
	for x, y := x1, y1; x <= x2; x++ {
		if steep {
			points = append(points, XY{y, x})
		} else {
			points = append(points, XY{x, y})
		}
		error -= abs(dy)
		if error < 0 {
			y += ystep
			error += dx
		}
	}
	if reversed {
		for i, j := 0, len(points)-1; i < j; i, j = i+1, j-1 {
			points[i], points[j] = points[j], points[i]
		}
	}
	return points
}

package main

import (
	"image/color"
	"sort"
)

type ID int

// MakeID returns a new unique ID.
var MakeID = func() func() ID {
	i := ID(0)
	return func() ID {
		defer func() { i++ }()
		return i
	}
}()

// Direction offsets.
var (
	North = XY{Y: -1}
	South = XY{Y: +1}
	West  = XY{X: -1}
	East  = XY{X: +1}
)

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
func (p XY) In(x0, y0, x1, y1 int) bool {
	return x0 <= p.X && p.X < x1 &&
		y0 <= p.Y && p.Y < y1
}

// Neighbors returns all neighbors of p.
func (p XY) Neighbors() [8]XY {
	var n [8]XY
	i := 0
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			n[i] = p.Add(XY{x, y})
			i++
		}
	}
	return n
}

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

// abs returns the absolute value of n.
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type State struct {
	Tiles    map[XY]Tile
	Entities map[ID]Entity
	at       map[XY][]ID
}

func NewState() *State {
	return &State{
		Tiles:    map[XY]Tile{},
		Entities: map[ID]Entity{},
		at:       map[XY][]ID{},
	}
}

func (s *State) Update() {
	for _, e := range s.Entities {
		e.Update()
	}
	s.at = map[XY][]ID{}
	for id, e := range s.Entities {
		p := e.Pos()
		s.at[p] = append(s.at[p], id)
	}
}

func (s *State) Add(e Entity) {
	id := MakeID()
	s.Entities[id] = e
}

// EntitiesAt returns all entities at the given position sorted by ID
// in increasing order.
func (s *State) EntitiesAt(xy XY) []Entity {
	ids := []ID{}
	for _, id := range s.at[xy] {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
	e := make([]Entity, 0, len(ids))
	for _, id := range ids {
		e = append(e, s.Entities[id])
	}
	return e
}

// Poser is implemented by any value that has a position on the grid.
type Poser interface {
	Pos() XY
}

// Symboler is implemented by any value that has a symbol and a color.
type Symboler interface {
	Symbol() Symbol
}

// Entity is implemented by any value that has a position, a symbol,
// and can update its state.
type Entity interface {
	Poser
	Symboler
	Update()
}

// Symbol represents a colored character on the grid.
type Symbol struct {
	Color color.Color
	Char  rune
}

type Tile int

const (
	Wall Tile = iota
	Floor
)

// Opaque returns true if the tile can pass light.
func (t Tile) Opaque() bool {
	return tileOpaque[t]
}

var tileOpaque = [...]bool{
	Wall:  true,
	Floor: false,
}

// Symbol implements the Symboler interface.
func (t Tile) Symbol() Symbol {
	return tileSymbol[t]
}

var tileSymbol = [...]Symbol{
	Wall:  {color.RGBA{0x45, 0x23, 0x0d, 0xff}, '▒'},
	Floor: {color.RGBA{0x2a, 0x1d, 0x0d, 0xff}, '.'},
}

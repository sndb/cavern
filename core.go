package main

import (
	"image/color"
)

type ID int

var MakeID = func() func() ID {
	i := ID(0)
	return func() ID {
		defer func() { i++ }()
		return i
	}
}()

var (
	North = XY{Y: -1}
	South = XY{Y: +1}
	West  = XY{X: -1}
	East  = XY{X: +1}
)

type XY struct {
	X int
	Y int
}

func (xy XY) Pos() XY {
	return xy
}

func (xy XY) Add(terms ...XY) XY {
	for _, t := range terms {
		xy.X += t.X
		xy.Y += t.Y
	}
	return xy
}

func (xy XY) Sub(terms ...XY) XY {
	for _, t := range terms {
		xy.X -= t.X
		xy.Y -= t.Y
	}
	return xy
}

func (xy XY) Mul(factors ...int) XY {
	for _, f := range factors {
		xy.X *= f
		xy.Y *= f
	}
	return xy
}

func (xy XY) Neighbors() []XY {
	n := make([]XY, 0, 8)
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			n = append(n, xy.Add(XY{x, y}))
		}
	}
	return n
}

func (xy XY) In(x0, y0, x1, y1 int) bool {
	return xy.X >= x0 && xy.X <= x1 && xy.Y >= y0 && xy.Y <= y1
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

func (s *State) At(xy XY) []Entity {
	r := []Entity{}
	for _, id := range s.at[xy] {
		r = append(r, s.Entities[id])
	}
	return r
}

type Entity interface {
	Symbol() Symbol
	Pos() XY
	Update()
}

type Symbol struct {
	Color color.Color
	Char  rune
}

type Tile int

const (
	Wall Tile = iota
	Floor
)

var tileSymbol = [...]Symbol{
	Wall:  {color.RGBA{0x45, 0x23, 0x0d, 0xff}, '\u2591'},
	Floor: {color.RGBA{0x2a, 0x1d, 0x0d, 0xff}, '\u2027'},
}

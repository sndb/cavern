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
	Door
	Arch
)

// Opaque returns true if the tile can pass light.
func (t Tile) Opaque() bool {
	return tileData[t].opaque
}

func (t Tile) Passable() bool {
	return tileData[t].passable
}

var tileData = [...]struct {
	opaque   bool
	passable bool
}{
	Wall:  {true, false},
	Floor: {false, true},
	Door:  {true, true},
	Arch:  {false, true},
}

// Symbol implements the Symboler interface.
func (t Tile) Symbol() Symbol {
	return tileSymbol[t]
}

var tileSymbol = [...]Symbol{
	Wall:  {color.RGBA{0x45, 0x23, 0x0d, 0xff}, '▒'},
	Floor: {color.RGBA{0x2a, 0x1d, 0x0d, 0xff}, '.'},
	Door:  {color.RGBA{0xa5, 0x62, 0x43, 0xff}, 'Ṩ'},
	Arch:  {color.RGBA{0x45, 0x23, 0x0d, 0xff}, 'ṧ'},
}

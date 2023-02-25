package main

import (
	"image/color"
)

// ID of a game entity.
type ID int

// MakeID returns a new unique ID.
var MakeID = func() func() ID {
	i := ID(0)
	return func() ID {
		defer func() { i++ }()
		return i
	}
}()

// Entity is implemented by any value that has a position, a symbol,
// and can update its state.
type Entity interface {
	Poser
	Symboler
	Update()
}

// Poser is implemented by any value that has a position on the grid.
type Poser interface {
	Pos() XY
}

// Symboler is implemented by any value that has a symbol and a color.
type Symboler interface {
	Symbol() Symbol
}

// Symbol represents a colored character on the grid.
type Symbol struct {
	Color color.RGBA
	Char  rune
}

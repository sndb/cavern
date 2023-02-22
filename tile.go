package main

import "image/color"

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

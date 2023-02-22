package main

// Generator is implemented by any values that generates something in
// the given bounds.
type Generator interface {
	Generate(map[XY]Tile, Rect)
}

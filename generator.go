package main

// Generator is implemented by any value that generates a tiled
// structure in the given bounds.
type Generator interface {
	Generate(map[XY]Tile, Rect)
}

package main

type Cave struct {
	Maze MazeFunc
	RDE1 int
	Grow int
	RDE2 int
}

// Generate generates a continuous cave map.
func (c Cave) Generate(tiles map[XY]Tile, bounds Rect) {
	c.Maze(tiles, bounds)
	for i := 0; i < c.RDE1; i++ {
		removeDeadEnds(tiles)
	}
	for i := 0; i < c.Grow; i++ {
		growMap(tiles)
	}
	for i := 0; i < c.RDE2; i++ {
		removeDeadEnds(tiles)
	}
}

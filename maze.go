package main

import "math/rand"

// MazeFunc generates a maze in the given bounds.
type MazeFunc func(map[XY]Tile, Rect)

// MazeDFS generates a maze in the given bounds.
// Implemented using Depth-First Search.
func MazeDFS(tiles map[XY]Tile, bounds Rect) {
	p := mazeStartingPoint(tiles, bounds)
	var dfs func(p XY)
	dfs = func(p XY) {
		dirs := [...]XY{North, South, West, East}
		rand.Shuffle(len(dirs), func(i, j int) {
			dirs[i], dirs[j] = dirs[j], dirs[i]
		})
		for _, dir := range dirs {
			if q := p.Add(dir.Mul(2)); q.In(bounds) && tiles[q] == Wall {
				tiles[p.Add(dir)] = Floor
				tiles[q] = Floor
				dfs(q)
			}
		}
	}
	dfs(p)
}

// MazePrim generates a maze in the given bounds.
// Implemented using Prim's algorithm.
func MazePrim(tiles map[XY]Tile, bounds Rect) {
	check := []XY{mazeStartingPoint(tiles, bounds)}
	for len(check) > 0 {
		var xy XY
		xy, check = randPop(check)
		if tiles[xy] == Floor {
			continue
		}
		tiles[xy] = Floor

		dirs := []XY{North, South, West, East}
		rand.Shuffle(len(dirs), func(i, j int) {
			dirs[i], dirs[j] = dirs[j], dirs[i]
		})
		for _, dir := range dirs {
			p := xy.Add(dir.Mul(2))
			if p.In(bounds) && tiles[p] == Floor {
				tiles[xy.Add(dir)] = Floor
				break
			}
		}
		for _, dir := range dirs {
			p := xy.Add(dir.Mul(2))
			if p.In(bounds) && tiles[p] == Wall {
				check = append(check, p)
			}
		}
	}
}

// mazeStartingPoint returns an odd point in the given bounds which
// has a Wall tile.
func mazeStartingPoint(tiles map[XY]Tile, bounds Rect) XY {
	const max = 1000
	for i := 0; i < max; i++ {
		p := bounds.OddPoint()
		if tiles[p] == Wall {
			return p
		}
	}
	panic("cannot get maze starting point")
}

// removeDeadEnds removes the tiles that have only one floor neighbor.
// Returns the number of tiles removed.
func removeDeadEnds(tiles map[XY]Tile) int {
	r := 0
	ends := []XY{}
	for xy := range tiles {
		neighbors := 0
		for _, dir := range []XY{North, South, West, East} {
			if tiles[xy.Add(dir)] == Wall {
				neighbors++
			}
		}
		if neighbors == 3 {
			ends = append(ends, xy)
		}
	}
	for _, end := range ends {
		delete(tiles, end)
		r++
	}
	return r
}

// growMap grows the map using cellular automata.
func growMap(tiles map[XY]Tile) {
	walls := map[XY]bool{}
	for floor := range tiles {
		for _, neighbor := range floor.Neighbors() {
			if tiles[neighbor] == Wall {
				walls[neighbor] = true
			}
		}
	}
	for wall := range walls {
		n := 0
		for _, neighbor := range wall.Neighbors() {
			if tiles[neighbor] == Floor {
				n++
			}
		}
		if n > 3 {
			tiles[wall] = Floor
		}
	}
}

package main

import "math/rand"

// Prim generates a continuous cave map using the Prim's Algorithm.
func Prim(x0, y0, x1, y1 int) map[XY]Tile {
	tiles := map[XY]Tile{}
	generateMaze(x0, x1, y0, y1, tiles)
	removeDeadEnds(tiles, 7)
	growMap(tiles, 3)
	removeDeadEnds(tiles, 3)
	return tiles
}

// popRandom removes and returns a random element from the list and
// the updated list.
func popRandom[T any](slice []T) (T, []T) {
	i := rand.Intn(len(slice))
	elem := slice[i]
	slice[i] = slice[len(slice)-1]
	slice = slice[:len(slice)-1]
	return elem, slice
}

// generateMaze generates a perfect maze in the given bounds.
func generateMaze(x0 int, x1 int, y0 int, y1 int, tiles map[XY]Tile) {
	check := []XY{{
		x0 + rand.Intn((x1-x0+1)/2)*2,
		y0 + rand.Intn((y1-y0+1)/2)*2,
	}}
	for len(check) > 0 {
		var xy XY
		xy, check = popRandom(check)
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
			if p.In(x0, y0, x1, y1) && tiles[p] == Floor {
				tiles[xy.Add(dir)] = Floor
				break
			}
		}
		for _, dir := range dirs {
			p := xy.Add(dir.Mul(2))
			if p.In(x0, y0, x1, y1) && tiles[p] == Wall {
				check = append(check, p)
			}
		}
	}
}

// removeDeadEnds removes the tiles that have only one floor neighbor.
func removeDeadEnds(tiles map[XY]Tile, n int) {
	for i := 0; i < n; i++ {
		ends := []XY{}
		for xy := range tiles {
			neighbors := 0
			for _, dir := range []XY{North, South, West, East} {
				if tiles[xy.Add(dir)] == Floor {
					neighbors++
				}
			}
			if neighbors == 1 {
				ends = append(ends, xy)
			}
		}
		for _, end := range ends {
			delete(tiles, end)
		}
	}
}

// growMap grows the map using cellular automata.
func growMap(tiles map[XY]Tile, n int) {
	for i := 0; i < n; i++ {
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
}

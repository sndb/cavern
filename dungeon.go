package main

import "math/rand"

type Dungeon struct {
	Maze         MazeFunc
	MaxRoomSize  XY
	RoomAttempts int
	Sparsity     float64
}

// Generate generates a continuous dungeon consisting of rooms and corridors.
func (d Dungeon) Generate(tiles map[XY]Tile, bounds Rect) {
	// The maze is region 0, the rooms are regions [1, d.RoomAttempts).
	regions := map[XY]int{}
	for i := 0; i < d.RoomAttempts; i++ {
		if r, ok := Room(tiles, bounds, d.MaxRoomSize); ok {
			r.Apply(func(p XY) {
				regions[p] = i + 1
			})
		}
	}
	d.Maze(tiles, bounds)

	conns := findConnectors(tiles, regions, bounds)
	rand.Shuffle(len(conns), func(i, j int) {
		conns[i], conns[j] = conns[j], conns[i]
	})

	merged := map[int]bool{}
	for len(conns) > 0 {
		// Merge regions if unmerged or the sparsity is high.
		conn := conns[len(conns)-1]
		conns = conns[:len(conns)-1]
		if merged[regions[conn.a]] && merged[regions[conn.b]] &&
			rand.Float64() > d.Sparsity {
			continue
		}
		passages := []Tile{Door, Arch}
		pass := passages[rand.Intn(len(passages))]
		tiles[conn.mid] = pass

		// Make all neighbor passages equal.
		for _, q := range conn.mid.Orthogonal() {
			if in(passages, tiles[q]) {
				floodFill(tiles, q, pass)
			}
		}

		merged[regions[conn.a]] = true
		merged[regions[conn.b]] = true
	}
	for removeDeadEnds(tiles) != 0 {
	}
}

// floodFill fills all tiles of the same type connected to p with t.
func floodFill(tiles map[XY]Tile, p XY, t Tile) {
	target := tiles[p]
	queue := []XY{p}
	seen := map[XY]bool{}
	for len(queue) > 0 {
		x := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		if seen[x] {
			continue
		}
		seen[x] = true

		tiles[x] = t
		for _, q := range x.Orthogonal() {
			if tiles[q] == target {
				queue = append(queue, q)
			}
		}
	}
}

// Room attemps to generate a room of the specified maximum size in
// the given bounds.
func Room(tiles map[XY]Tile, bounds Rect, maxSize XY) (r Rect, ok bool) {
	const min = 3
	p := bounds.OddPoint()
	r = Rect{
		p.X,
		p.Y,
		p.X + min + rand.Intn((maxSize.X-min+1)/2)*2,
		p.Y + min + rand.Intn((maxSize.Y-min+1)/2)*2,
	}
	if !r.In(bounds) {
		return r, false
	}

	good := true
	r.Apply(func(p XY) {
		if tiles[p] == Floor {
			good = false
			return
		}
	})
	if !good {
		return r, false
	}
	r.Apply(func(p XY) {
		tiles[p] = Floor
	})
	return r, true
}

type connector struct {
	mid, a, b XY
}

// findConnectors returns all connectors that can be used to merge different regions.
func findConnectors(tiles map[XY]Tile, regions map[XY]int, bounds Rect) []connector {
	r := []connector{}
	bounds.Apply(func(p XY) {
		if tiles[p] == Floor {
			return
		}
		n, s, w, e := p.N(), p.S(), p.W(), p.E()
		tp, tn, ts, tw, te := tiles[p], tiles[n], tiles[s], tiles[w], tiles[e]
		if !tp.Passable() && tw.Passable() && te.Passable() && regions[w] != regions[e] {
			r = append(r, connector{p, w, e})
		}
		if !tp.Passable() && tn.Passable() && ts.Passable() && regions[n] != regions[s] {
			r = append(r, connector{p, n, s})
		}
	})
	return r
}

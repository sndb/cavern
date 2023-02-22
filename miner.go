package main

import (
	"image/color"
	"math/rand"
)

type Miner struct {
	XY
	Bounds Rect
	Energy int
	State  *State
}

func (m *Miner) Update() {
	// Do nothing if dead.
	if m.Energy <= 0 {
		return
	}
	m.Energy--

	// Move.
	p := m.XY.Add([]XY{{}, North, South, West, East}[rand.Intn(5)])
	if !p.In(m.Bounds) {
		// Recurse if out of bounds.
		m.Update()
		return
	}
	m.XY = p

	// Dig.
	if m.State.Tiles[m.XY] == Wall {
		m.State.Tiles[m.XY] = Floor
		if rand.Float64() < 0.3 {
			m.State.Add(&Stone{m.XY})
		}
	}

	// Spawn another miner.
	if rand.Float64() < 0.1 {
		m.State.Add(&Miner{p, m.Bounds, m.Energy, m.State})
	}

	// Die if surrounded by empty space.
	empty := 0
	for _, neigh := range p.Neighbors() {
		if m.State.Tiles[neigh] == Floor {
			empty++
		}
	}
	if empty == 8 {
		m.Energy = 0
	}
}

func (m *Miner) Symbol() Symbol {
	if m.Energy <= 0 {
		return Symbol{color.RGBA{0x55, 0x0f, 0x0a, 0xff}, 'Ḳ'}
	}
	return Symbol{color.RGBA{0xef, 0xac, 0x28, 0xff}, 'Ḳ'}
}

type Stone struct {
	XY
}

func (s *Stone) Update() {}

func (s *Stone) Symbol() Symbol {
	return Symbol{color.RGBA{0x45, 0x23, 0x0d, 0xff}, '●'}
}

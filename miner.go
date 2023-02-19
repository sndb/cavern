package main

import (
	"image"
	"image/color"
	"math/rand"
)

type Miner struct {
	XY
	Bounds image.Rectangle
	Energy int
	State  *State
}

func (m *Miner) Update() {
	if m.Energy <= 0 {
		return
	}
	m.Energy--

	xy := m.XY.Add([]XY{{}, North, South, West, East}[rand.Intn(5)])
	if !xy.In(
		m.Bounds.Min.X, m.Bounds.Min.Y,
		m.Bounds.Max.X, m.Bounds.Max.Y,
	) {
		m.Update()
		return
	}
	m.XY = xy
	if m.State.Tiles[m.XY] == Wall {
		m.State.Tiles[m.XY] = Floor
		if rand.Float64() < 0.3 {
			m.State.Add(&Stone{m.XY})
		}
	}
}

func (m *Miner) Symbol() Symbol {
	return Symbol{color.RGBA{0xef, 0xac, 0x28, 0xff}, '\u263a'}
}

type Stone struct {
	XY
}

func (s *Stone) Update() {}

func (s *Stone) Symbol() Symbol {
	return Symbol{color.RGBA{0x45, 0x23, 0x0d, 0xff}, '\u25cf'}
}

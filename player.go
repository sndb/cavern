package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	XY
	Explored map[XY]bool
	FOV      map[XY]bool
	Radius   int
	State    *State
}

func (p *Player) Symbol() Symbol {
	return Symbol{color.RGBA{0xef, 0xac, 0x28, 0xff}, '☺'}
}

func (p *Player) Update() {
	p.FOV = p.XY.FOV(p.Radius, func(xy XY) bool { return p.State.Tiles[xy].Opaque() })
	for xy := range p.FOV {
		p.Explored[xy] = true
	}
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyUp):
		v := p.XY.Add(North)
		if !p.State.Tiles[v].Opaque() {
			p.XY = v
		}
	case ebiten.IsKeyPressed(ebiten.KeyDown):
		v := p.XY.Add(South)
		if !p.State.Tiles[v].Opaque() {
			p.XY = v
		}
	case ebiten.IsKeyPressed(ebiten.KeyLeft):
		v := p.XY.Add(West)
		if !p.State.Tiles[v].Opaque() {
			p.XY = v
		}
	case ebiten.IsKeyPressed(ebiten.KeyRight):
		v := p.XY.Add(East)
		if !p.State.Tiles[v].Opaque() {
			p.XY = v
		}
	}
}

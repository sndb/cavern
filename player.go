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
	Updated  bool
}

func (p *Player) Symbol() Symbol {
	return Symbol{color.RGBA{0xef, 0xac, 0x28, 0xff}, '☺'}
}

func (p *Player) Update() {
	if p.Updated {
		return
	}
	p.FOV = p.XY.FOV(p.Radius, func(xy XY) bool { return p.State.Tiles[xy].Opaque() })
	for xy := range p.FOV {
		p.Explored[xy] = true
	}
	var offset XY
	switch {
	case pressed(ebiten.KeyUp, ebiten.KeyNumpad8, ebiten.KeyW):
		offset = North
	case pressed(ebiten.KeyDown, ebiten.KeyNumpad2, ebiten.KeyX):
		offset = South
	case pressed(ebiten.KeyLeft, ebiten.KeyNumpad4, ebiten.KeyA):
		offset = West
	case pressed(ebiten.KeyRight, ebiten.KeyNumpad6, ebiten.KeyD):
		offset = East
	case pressed(ebiten.KeyNumpad7, ebiten.KeyQ):
		offset = North.Add(West)
	case pressed(ebiten.KeyNumpad9, ebiten.KeyE):
		offset = North.Add(East)
	case pressed(ebiten.KeyNumpad1, ebiten.KeyZ):
		offset = South.Add(West)
	case pressed(ebiten.KeyNumpad3, ebiten.KeyC):
		offset = South.Add(East)
	case pressed(ebiten.KeyNumpad5, ebiten.KeyS):
		p.Updated = true
	}
	var zero XY
	if q := p.XY.Add(offset); offset != zero && p.State.Tiles[q].Passable() {
		p.XY = q
		p.Updated = true
	}
}

func pressed(ks ...ebiten.Key) bool {
	for _, k := range ks {
		if ebiten.IsKeyPressed(k) {
			return true
		}
	}
	return false
}

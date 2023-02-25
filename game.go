package main

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	State    *State
	Terminal *Terminal
	Player   *Player
	Bounds   Rect
	Start    time.Time
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Terminal.Set(screen)
	background := color.RGBA{0x15, 0x0f, 0x0a, 0xff}
	dy, dx := g.Terminal.Dimensions.Y, g.Terminal.Dimensions.X
	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			// Dock camera to the edge of the map.
			p := XY{x + g.Player.X - dx/2, y + g.Player.Y - dy/2}
			switch {
			case g.Player.X < dx/2:
				p.X = x
			case g.Player.X >= g.Bounds.X1-dx/2:
				p.X = x + g.Bounds.X1 - dx
			}
			switch {
			case g.Player.Y < dy/2:
				p.Y = y
			case g.Player.Y >= g.Bounds.Y1-dy/2:
				p.Y = y + g.Bounds.Y1 - dy
			}

			c := Cell{Bg: background}
			if ents := g.State.EntitiesAt(p); g.Player.FOV[p] && len(ents) > 0 {
				ent := displayedEntity(g.Start, ents).Symbol()
				c.Fg = ent.Color
				c.Symbol = ent.Char
			} else {
				tile := g.State.Tiles[p].Symbol()
				c.Fg = tile.Color
				c.Symbol = tile.Char
			}
			op := &ebiten.DrawImageOptions{}
			switch {
			case g.Player.FOV[p]:
				// Visible; draw as is.
			case g.Player.Explored[p]:
				// Explored; draw shadowed.
				op.ColorM.ChangeHSV(math.Pi, 0.5, 0.75)
			default:
				// Unexplored; draw black.
				op.ColorM.ChangeHSV(0, 0, 0)
			}
			g.Terminal.Print(XY{x, y}, c, op)
		}
	}
}

func displayedEntity(gameStart time.Time, e []Entity) Entity {
	const period = 400
	dt := int(time.Since(gameStart).Milliseconds())
	return e[(dt/period)%len(e)]
}

func (g *Game) Update() error {
	g.Player.Update()
	if g.Player.Updated {
		g.State.Update()
		g.Player.Updated = false
	}
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.Terminal.Layout()
}

package main

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	State    *State
	Drawer   *Drawer
	Player   *Player
	Viewport XY
	Bounds   Rect
	Start    time.Time
}

func (g *Game) Draw(screen *ebiten.Image) {
	background := color.RGBA{0x15, 0x0f, 0x0a, 0xff}
	grid := Grid{Cells: []Cell{}, Columns: g.Viewport.X}
	for y := -g.Viewport.Y / 2; y <= g.Viewport.Y/2; y++ {
		for x := -g.Viewport.X / 2; x <= g.Viewport.X/2; x++ {
			pos := g.Player.Pos().Add(XY{x, y})
			cell := Cell{Bg: background}
			if ents := g.State.EntitiesAt(pos); g.Player.FOV[pos] && len(ents) > 0 {
				ent := displayedEntity(g.Start, ents).Symbol()
				cell.Fg = ent.Color
				cell.Symbol = ent.Char
			} else {
				tile := g.State.Tiles[pos].Symbol()
				cell.Fg = tile.Color
				cell.Symbol = tile.Char
			}
			switch {
			case g.Player.FOV[pos]:
				cell.Style = CellFOV
			case g.Player.Explored[pos]:
				cell.Style = CellExplored
			}
			grid.Cells = append(grid.Cells, cell)
		}
	}
	g.Drawer.Grid(screen, grid)
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
	return g.Drawer.TileWidth * g.Viewport.X,
		g.Drawer.TileHeight * g.Viewport.Y
}

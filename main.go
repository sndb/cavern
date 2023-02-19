package main

import (
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	State  *State
	Drawer *Drawer
	Width  int
	Height int
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x15, 0x0f, 0x0a, 0xff})
	for x := 0; x < g.Width; x++ {
		for y := 0; y < g.Height; y++ {
			if len(g.State.At(XY{x, y})) > 0 {
				continue
			}
			i := g.Drawer.Symbol(tileSymbol[g.State.Tiles[XY{x, y}]])
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(
				float64(x*g.Drawer.TileWidth),
				float64(y*g.Drawer.TileHeight),
			)
			screen.DrawImage(i, op)
		}
	}
	for _, e := range g.State.Entities {
		i := g.Drawer.Symbol(e.Symbol())
		xy := e.Pos()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(
			float64(xy.X*g.Drawer.TileWidth),
			float64(xy.Y*g.Drawer.TileHeight),
		)
		screen.DrawImage(i, op)
	}
}

func (g *Game) Update() error {
	g.State.Update()
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.Drawer.TileWidth * g.Width, g.Drawer.TileHeight * g.Height
}

func main() {
	// Initialize game.
	state := NewState()
	drawer := NewDrawer()
	game := &Game{state, drawer, 121, 61}
	state.Tiles = Prim(1, 1, game.Width-1, game.Height-1)

	floors := []XY{}
	for xy, tile := range state.Tiles {
		if tile == Floor {
			floors = append(floors, xy)
		}
	}

	// Add entities.
	for i := 0; i < 10; i++ {
		xy := floors[rand.Intn(len(floors))]
		game.State.Add(&Miner{
			XY:     xy,
			Bounds: image.Rect(1, 1, game.Width-1, game.Height-1),
			Energy: rand.Intn(1000),
			State:  state,
		})
	}

	// Run the game.
	width, height := game.Layout(0, 0)
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Cave")
	ebiten.SetTPS(30)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

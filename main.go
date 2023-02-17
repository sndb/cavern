package main

import (
	"image"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Grid    *Grid
	Walkers []*Walker
}

func (g *Game) Update() error {
	for _, w := range g.Walkers {
		g.Grid.Fill(w.p)
		w.Walk()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(Black)
	wall := g.Grid.Cell(screen, Dark)
	space := g.Grid.Cell(screen, Brown)
	walker := g.Grid.Cell(screen, White)
	for p := range g.Grid.points {
		space.Draw(screen, p)
		for x := -1; x <= 1; x++ {
			for y := -1; y <= 1; y++ {
				p := p.Add(image.Point{x, y})
				if !g.Grid.points[p] {
					wall.Draw(screen, p)
				}
			}
		}
	}
	for _, w := range g.Walkers {
		walker.Draw(screen, w.p)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	const (
		walkersN = 4
		width    = 320
		height   = 240
		spread   = 20
	)
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("Maze")
	ebiten.SetTPS(120)

	game := &Game{}
	game.Grid = NewGrid(width, height)
	for i := 0; i < walkersN; i++ {
		w := &Walker{
			p: image.Point{width / 2, height / 2},
			r: image.Rect(1, 1, width-1, height-1),
		}
		w.p.X += -spread/2 + rand.Intn(spread)
		w.p.Y += -spread/2 + rand.Intn(spread)
		game.Walkers = append(game.Walkers, w)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Initialize the game.
	game := &Game{
		State:    NewState(),
		Drawer:   NewDrawer(),
		Viewport: XY{81, 61},
		Bounds:   Rect{0, 0, 81, 81},
		Start:    time.Now(),
	}

	// Generate a map.
	gens := []Generator{
		Dungeon{MazeDFS, XY{15, 15}, 100, 0.02},
		Cave{MazeDFS, 400, 2, 3},
		Cave{MazePrim, 7, 3, 3},
	}
	gens[rand.Intn(len(gens))].Generate(game.State.Tiles, game.Bounds)

	// Add the player.
	game.Player = NewPlayer(game.State.RandomPosition(), 20, game.State)
	game.State.Add(game.Player)

	// Add some monsters.
	for i := 0; i < 0; i++ {
		game.State.Add(&Miner{
			XY:     game.State.RandomPosition(),
			Bounds: game.Bounds.Inset(1),
			Energy: rand.Intn(1000),
			State:  game.State,
		})
	}

	// Run the game.
	width, height := game.Layout(0, 0)
	ebiten.SetWindowSize(2*width, 2*height)
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	ebiten.SetWindowTitle("Cave")
	ebiten.SetTPS(12)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

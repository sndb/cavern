package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	State  *State
	Drawer *Drawer
	Player *Player
	Width  int
	Height int
	Start  time.Time
	DrawFn func(*ebiten.Image)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawFn(screen)
}

func (g *Game) DrawPlayer(screen *ebiten.Image) {
	screen.Fill(color.Black)
	for p := range g.Player.Explored {
		tile := g.State.Tiles[p]
		sym := tile.Symbol()
		if len(g.State.EntitiesAt(p)) > 0 {
			sym.Char = ' '
		}
		op := &ebiten.DrawImageOptions{}
		if !g.Player.FOV[p] {
			op.ColorM.ChangeHSV(math.Pi, 0.5, 0.75)
		}
		bg := g.Drawer.Background(color.RGBA{0x15, 0x0f, 0x0a, 0xff})
		fg := g.Drawer.Symbol(sym)
		g.Drawer.DrawWithOptions(screen, bg, p, op)
		g.Drawer.DrawWithOptions(screen, fg, p, op)
	}
	for p := range g.Player.FOV {
		entities := g.State.EntitiesAt(p)
		if len(entities) == 0 {
			continue
		}
		dt := int(time.Since(g.Start).Milliseconds())
		e := entities[(dt/400)%len(entities)]
		g.Drawer.Draw(screen, g.Drawer.Symbol(e.Symbol()), e.Pos())
	}
}

func (g *Game) DrawFull(screen *ebiten.Image) {
	screen.Fill(color.Black)
	for x := 0; x < g.Width; x++ {
		for y := 0; y < g.Height; y++ {
			p := XY{x, y}
			tile := g.State.Tiles[p]
			sym := tile.Symbol()
			entities := g.State.EntitiesAt(p)
			if len(entities) > 0 {
				sym.Char = ' '
			}
			bg := g.Drawer.Background(color.RGBA{0x15, 0x0f, 0x0a, 0xff})
			fg := g.Drawer.Symbol(sym)
			g.Drawer.Draw(screen, bg, p)
			g.Drawer.Draw(screen, fg, p)

			if len(entities) == 0 {
				continue
			}
			dt := int(time.Since(g.Start).Milliseconds())
			e := entities[(dt/400)%len(entities)]
			g.Drawer.Draw(screen, g.Drawer.Symbol(e.Symbol()), e.Pos())
		}
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.DrawFn = g.DrawPlayer
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		g.DrawFn = g.DrawFull
	}

	g.Player.Update()
	if g.Player.Updated {
		g.State.Update()
		g.Player.Updated = false
	}
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.Drawer.TileWidth * g.Width, g.Drawer.TileHeight * g.Height
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Initialize game.
	state := NewState()
	drawer := NewDrawer()
	game := &Game{
		State:  state,
		Drawer: drawer,
		Width:  81,
		Height: 61,
		Start:  time.Now(),
	}
	game.DrawFn = game.DrawPlayer
	generators := []Generator{
		Dungeon{MazeDFS, XY{15, 15}, 50, 0.02},
		Cave{MazeDFS, 400, 2, 3},
		Cave{MazePrim, 7, 3, 3},
	}
	gen := generators[rand.Intn(len(generators))]
	gen.Generate(state.Tiles, Rect{0, 0, game.Width, game.Height})

	// Find empty tiles.
	empty := []XY{}
	for xy, tile := range state.Tiles {
		if tile == Floor {
			empty = append(empty, xy)
		}
	}

	// Add entities.
	for i := 0; i < 0; i++ {
		xy := empty[rand.Intn(len(empty))]
		game.State.Add(&Miner{
			XY:     xy,
			Bounds: Rect{1, 1, game.Width - 1, game.Height - 1},
			Energy: rand.Intn(1000),
			State:  state,
		})
	}

	// Add player.
	player := &Player{
		XY:       empty[rand.Intn(len(empty))],
		Explored: map[XY]bool{},
		FOV:      map[XY]bool{},
		Radius:   20,
		State:    state,
		Updated:  true,
	}
	game.Player = player
	game.State.Add(player)

	// Run the game.
	width, height := game.Layout(0, 0)
	ebiten.SetWindowSize(2*width, 2*height)
	ebiten.SetWindowTitle("Cave")
	ebiten.SetTPS(12)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

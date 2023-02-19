package main

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed font.ttf
var fontData []byte

type Drawer struct {
	TileWidth  int
	TileHeight int
	Font       font.Face
}

func NewDrawer() *Drawer {
	const (
		fontSize   = 12
		tileWidth  = 8
		tileHeight = 12
	)
	f, err := opentype.Parse(fontData)
	if err != nil {
		panic(err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		panic(err)
	}
	return &Drawer{
		Font:       face,
		TileWidth:  tileWidth,
		TileHeight: tileHeight,
	}
}

var drawCache = map[rune]*ebiten.Image{}

func (d *Drawer) Symbol(s Symbol) *ebiten.Image {
	if i, ok := drawCache[s.Char]; ok {
		return i
	}
	i := ebiten.NewImage(d.TileWidth, d.TileHeight)
	text.Draw(i, string(s.Char), d.Font, 0, 10, s.Color)
	drawCache[s.Char] = i
	return i
}

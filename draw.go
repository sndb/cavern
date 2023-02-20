package main

import (
	_ "embed"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed font.ttf
var fontData []byte

type drawerCacheKey struct {
	c          rune
	r, g, b, a uint32
}

type Drawer struct {
	TileWidth       int
	TileHeight      int
	Font            font.Face
	backgroundCache map[drawerCacheKey]*ebiten.Image
	symbolCache     map[drawerCacheKey]*ebiten.Image
}

func NewDrawer() *Drawer {
	const (
		fontSize   = 8
		tileWidth  = 6
		tileHeight = 8
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
		Font:            face,
		TileWidth:       tileWidth,
		TileHeight:      tileHeight,
		backgroundCache: map[drawerCacheKey]*ebiten.Image{},
		symbolCache:     map[drawerCacheKey]*ebiten.Image{},
	}
}

func (d *Drawer) Symbol(s Symbol) *ebiten.Image {
	r, g, b, a := s.Color.RGBA()
	key := drawerCacheKey{s.Char, r, g, b, a}
	if i, ok := d.symbolCache[key]; ok {
		return i
	}

	i := ebiten.NewImage(d.TileWidth, d.TileHeight)
	text.Draw(i, string(s.Char), d.Font, 0, 8, s.Color)
	d.symbolCache[key] = i
	return i
}

func (d *Drawer) Background(c color.Color) *ebiten.Image {
	r, g, b, a := c.RGBA()
	key := drawerCacheKey{' ', r, g, b, a}
	if i, ok := d.backgroundCache[key]; ok {
		return i
	}

	i := ebiten.NewImage(d.TileWidth, d.TileHeight)
	i.Fill(c)
	d.backgroundCache[key] = i
	return i
}

func (d *Drawer) Draw(dst, src *ebiten.Image, p XY) {
	d.DrawWithOptions(dst, src, p, &ebiten.DrawImageOptions{})
}

func (d *Drawer) DrawWithOptions(dst, src *ebiten.Image, p XY, op *ebiten.DrawImageOptions) {
	op.GeoM.Reset()
	op.GeoM.Translate(
		float64(p.X*d.TileWidth),
		float64(p.Y*d.TileHeight),
	)
	dst.DrawImage(src, op)
}

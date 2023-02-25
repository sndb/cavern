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

func ParseFont(data []byte, size float64) font.Face {
	f, err := opentype.Parse(fontData)
	if err != nil {
		panic(err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		panic(err)
	}
	return face
}

type Terminal struct {
	TileSize   XY
	Dimensions XY
	Buffer     *ebiten.Image
	Font       font.Face
}

func (t *Terminal) Set(buffer *ebiten.Image) {
	buffer.Fill(color.Black)
	t.Buffer = buffer
}

func (t *Terminal) Layout() (width, height int) {
	return t.Dimensions.X * t.TileSize.X,
		t.Dimensions.Y * t.TileSize.Y
}

func (t *Terminal) Print(p XY, c Cell, op *ebiten.DrawImageOptions) {
	if op == nil {
		op = &ebiten.DrawImageOptions{}
	}
	op.GeoM.Reset()
	op.GeoM.Translate(float64(p.X*t.TileSize.X), float64(p.Y*t.TileSize.Y))
	t.Buffer.DrawImage(t.getCellImage(c), op)
}

type Cell struct {
	Fg, Bg color.RGBA
	Symbol rune
}

var cellImageCache = map[Cell]*ebiten.Image{}

func (t *Terminal) getCellImage(c Cell) *ebiten.Image {
	if img, ok := cellImageCache[c]; ok {
		return img
	}
	img := ebiten.NewImage(t.TileSize.X, t.TileSize.Y)
	img.Fill(c.Bg)
	drawCenteredText(img,
		t.Font, string(c.Symbol), c.Fg,
		t.TileSize.X/2, t.TileSize.Y/2,
	)
	cellImageCache[c] = img
	return img
}

func drawCenteredText(screen *ebiten.Image, font font.Face, s string, clr color.Color, cx, cy int) {
	bounds := text.BoundString(font, s)
	x, y := cx-bounds.Min.X-bounds.Dx()/2, cy-bounds.Min.Y-bounds.Dy()/2
	text.Draw(screen, s, font, x, y, clr)
}

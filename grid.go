package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Grid struct {
	points map[image.Point]bool
	width  int
	height int
}

func NewGrid(width, height int) *Grid {
	return &Grid{make(map[image.Point]bool), width, height}
}

func (g *Grid) Fill(p image.Point) {
	g.points[p] = true
}

func (g *Grid) Cell(screen *ebiten.Image, c color.Color) *Cell {
	width, height := screen.Size()
	width /= g.width
	height /= g.height

	cell := &Cell{ebiten.NewImage(width, height)}
	cell.Fill(c)
	return cell
}

type Cell struct {
	*ebiten.Image
}

func (c *Cell) Draw(screen *ebiten.Image, p image.Point) {
	width, height := c.Size()
	tx, ty := width*p.X, height*p.Y
	t := ebiten.GeoM{}
	t.Translate(float64(tx), float64(ty))
	screen.DrawImage(c.Image, &ebiten.DrawImageOptions{GeoM: t})
}

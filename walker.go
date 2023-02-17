package main

import (
	"image"
	"math/rand"
)

type Walker struct {
	p image.Point
	r image.Rectangle
}

func (w *Walker) Walk() {
	offset := image.Point{}
	if rand.Intn(2) == 0 {
		offset.X += -1 + rand.Intn(3)
	} else {
		offset.Y += -1 + rand.Intn(3)
	}
	moved := w.p.Add(offset)
	if !moved.In(w.r) {
		w.Walk()
		return
	}
	w.p = moved
}

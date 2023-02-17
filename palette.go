package main

type RGB uint32

func (c RGB) RGBA() (r, g, b, a uint32) {
	v := uint32(c)
	r = v >> 16
	r |= r << 8
	g = v >> 8 & 0xff
	g |= g << 8
	b = v & 0xff
	b |= b << 8
	a = 0xffff
	return
}

const (
	Black  RGB = 0x2c1f1a
	Dark   RGB = 0x7d4735
	Brown  RGB = 0xb36d41
	Yellow RGB = 0xd89840
	White  RGB = 0xeee4bd
	Salad  RGB = 0xc4c499
	Blue   RGB = 0x7f9081
	Green  RGB = 0x61806a
)

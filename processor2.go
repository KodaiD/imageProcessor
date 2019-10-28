package main

import (
	"image"
	"image/color"
)

func averagecolor(img image.Image, ) (uint8, uint8, uint8) {
	tileSize := 10
	red, green, blue := 0, 0, 0
	for j := y; j <= y +tileSize; j++ {
		for i := x; i <= x + tileSize; i++ {
			r, g, b, _ := img.At(i, j).RGBA()
			red += int(r)
			green += int(g)
			blue += int(b)
		}
	}
	redAV := uint8(red / tileSize)
	greenAV := uint8(green / tileSize)
	blueAV := uint8(blue / tileSize)
	return redAV, greenAV, blueAV
}
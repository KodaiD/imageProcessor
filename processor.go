package main

import (
	"image"
	"image/color"
)

func averageColor(img image.RGBA) (uint8, uint8, uint8) {
	red, blue, green := 0, 0, 0
	bounds := img.Bounds()
	totalPixel := (bounds.Max.X - bounds.Min.X ) * (bounds.Max.Y - bounds.Min.Y)
	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		for x := bounds.Min.X; x <= bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			red += int(r)
			green += int(g)
			blue += int(b)
		}
	}
	redAv := uint8(red / totalPixel)
	greenAv := uint8(green / totalPixel)
	blueAv := uint8(blue / totalPixel)

	return redAv, greenAv, blueAv
}

func trimming(img image.Image, bounds image.Rectangle) image.RGBA {
	newImage := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		for x := bounds.Min.X; x <= bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			newImage.Set(x, y, color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: uint8(a),
			})
		}
	}
	return *newImage
}
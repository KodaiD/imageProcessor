package main

import (
	"bytes"
	"encoding/base64"
	"html/template"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"
)

func Mosaic(w http.ResponseWriter, r *http.Request) {
	tileSize, _ := strconv.Atoi(r.FormValue("tile_size"))
	r.ParseMultipartForm(10485760)
	file, _, _ := r.FormFile("image")
	defer file.Close()
	original, _, _ := image.Decode(file)
	bounds := original.Bounds()
	newImage := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y <= bounds.Max.Y; y += tileSize {
		for x := bounds.Min.X; x <= bounds.Max.X; x += tileSize {
			var s int
			var redSum, greenSum, blueSum float64
			for j := y; j < y + tileSize; j++ {
				for i := x; i < x + tileSize; i++ {
					if j < bounds.Max.Y && i < bounds.Max.X {
						c := color.RGBAModel.Convert(original.At(i, j)).(color.RGBA)
						red, green, blue := c.R, c.G, c.B
						redSum += float64(red)
						greenSum += float64(green)
						blueSum += float64(blue)
						s = tileSize * tileSize
					} else if j <= bounds.Max.Y && i > bounds.Max.X {
						s = tileSize * (tileSize - (x + tileSize - bounds.Max.X))
					} else if i <= bounds.Max.X && j > bounds.Max.Y {
						s = tileSize * (tileSize - (y + tileSize - bounds.Max.Y))
					}
				}
			}
			redAV := uint8(redSum / float64(s))
			greenAV := uint8(greenSum / float64(s))
			blueAV := uint8(blueSum / float64(s))

			for j := y; j < y + tileSize; j++ {
				for i := x; i < x + tileSize; i++ {
					if i < bounds.Max.X && j < bounds.Max.Y {
						newImage.Set(i, j, color.RGBA{redAV, greenAV, blueAV, 255})
					}
				}
			}
		}
	}


	// 元の画像
	buf1 := new(bytes.Buffer)
	err := jpeg.Encode(buf1, original, nil)
	if err != nil {
		png.Encode(buf1, original)
	}
	originalStr := base64.StdEncoding.EncodeToString(buf1.Bytes())

	// 加工後
	buf2 := new(bytes.Buffer)
	jpeg.Encode(buf2, newImage, nil)
	newStr := base64.StdEncoding.EncodeToString(buf2.Bytes())

	images := map[string]string{
		"original": originalStr,
		"new":   newStr,
	}

	t, _ := template.ParseFiles("templates/result1.html")
	t.Execute(w, images)
}
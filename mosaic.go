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

	for y := bounds.Min.Y; y < bounds.Max.Y; y += tileSize {
		for x := bounds.Min.X; x < bounds.Max.X; x += tileSize {
			redAV, greenAV, blueAV := getAverageColor(original, x, y, bounds, tileSize)
			for j := y; j < y + tileSize; j++ {
				for i := x; i < x + tileSize; i++ {
					if i < bounds.Max.X && j < bounds.Max.Y {
						newImage.Set(i, j, color.RGBA{redAV, greenAV, blueAV, 255})
					}
				}
			}
		}
	}

	// before process
	buf1 := new(bytes.Buffer)
	err := jpeg.Encode(buf1, original, nil)
	if err != nil {
		png.Encode(buf1, original)
	}
	originalStr := base64.StdEncoding.EncodeToString(buf1.Bytes())

	// after process
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

func getAverageColor(original image.Image, x int, y int, bounds image.Rectangle, tileSize int) (uint8, uint8, uint8) {
	var s int
	var redSum, greenSum, blueSum int

	for j := y; j < y + tileSize; j++ {
		for i := x; i < x + tileSize; i++ {
			if j < bounds.Max.Y && i < bounds.Max.X {
				c := color.RGBAModel.Convert(original.At(i, j)).(color.RGBA)
				red, green, blue := int(c.R), int(c.G), int(c.B)
				redSum += red
				greenSum += green
				blueSum += blue
			}
		}
	}

	if y+tileSize < bounds.Max.Y && x+tileSize < bounds.Max.X {
		s = tileSize * tileSize
	} else if y+tileSize < bounds.Max.Y {
		s = tileSize * (tileSize - (x + tileSize - bounds.Max.X))
	} else if x+tileSize < bounds.Max.X {
		s = tileSize * (tileSize - (y + tileSize - bounds.Max.Y))
	} else {
		s = (tileSize - (x + tileSize - bounds.Max.X)) * (tileSize - (y + tileSize - bounds.Max.Y))
	}

	redAV := uint8(redSum / s)
	greenAV := uint8(greenSum / s)
	blueAV := uint8(blueSum / s)

	return redAV, greenAV, blueAV
}


//func partition(original image.Image, tileSize, x1, y1, x2, y2 int) <-chan image.Image {
//	c := make(chan image.Image)
//	sp := image.Point{}
//	go func() {
//		newImage := image.NewRGBA(image.Rect(x1, y1, x2, y2))
//		newImage.SubImage()
//	}()
//
//}
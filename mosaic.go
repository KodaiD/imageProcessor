package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func Mosaic(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	tileSize, _ := strconv.Atoi(r.FormValue("tile_size"))
	r.ParseMultipartForm(10485760)
	file, _, _ := r.FormFile("image")
	defer file.Close()

	original, _, _ := image.Decode(file)
	bounds := original.Bounds()
	c1 := partition(original, tileSize, bounds.Min.X, bounds.Min.Y, bounds.Max.X/2, bounds.Max.Y/2)
	c2 := partition(original, tileSize, bounds.Max.X/2, bounds.Min.Y, bounds.Max.X, bounds.Max.Y/2)
	c3 := partition(original, tileSize, bounds.Min.X, bounds.Max.Y/2, bounds.Max.X/2, bounds.Max.Y)
	c4 := partition(original, tileSize, bounds.Max.X/2, bounds.Max.Y/2, bounds.Max.X, bounds.Max.Y)
	c := combine(bounds, c1, c2, c3, c4)
	newImage := <-c
	end := time.Now()
	fmt.Println(end.Sub(start))

	// before process
	buf1 := new(bytes.Buffer)
	err := jpeg.Encode(buf1, original, nil)
	if err != nil {
		png.Encode(buf1, original)
	}
	originalStr := base64.StdEncoding.EncodeToString(buf1.Bytes())

	images := map[string]string{
		"original": originalStr,
		"new":   newImage,
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

func partition(original image.Image, tileSize, x1, y1, x2, y2 int) <-chan image.Image {
	c := make(chan image.Image)
	go func() {
		newImage := image.NewRGBA(image.Rect(x1, y1, x2, y2))
		for y := y1; y < y2; y += tileSize {
			for x := x1; x < x2; x += tileSize {
				redAV, greenAV, blueAV := getAverageColor(original, x, y, image.Rect(x1, y1, x2, y2), tileSize)
				for j := y; j < y + tileSize; j++ {
					for i := x; i < x + tileSize; i++ {
						if i < x2 && j < y2 {
							newImage.Set(i, j, color.RGBA{redAV, greenAV, blueAV, 255})
						}
					}
				}
			}
		}
		c <- newImage
	}()
	return c
}

func combine(r image.Rectangle, c1, c2, c3, c4 <-chan image.Image) <-chan string {
	c := make(chan string)
	go func() {
		var wg sync.WaitGroup
		newImage := image.NewRGBA(r)
		copyImage := func(dst draw.Image, r image.Rectangle, src image.Image, startPoint image.Point) {
			draw.Draw(dst, r, src, startPoint, draw.Src)
			wg.Done()
		}
		wg.Add(4)
		var s1, s2, s3, s4 image.Image
		var ok1, ok2, ok3, ok4 bool
		for {
			select {
			case s1, ok1 = <-c1:
				go copyImage(newImage, s1.Bounds(), s1, image.Point{r.Min.X, r.Min.Y})
			case s2, ok2 = <-c2:
				go copyImage(newImage, s2.Bounds(), s2, image.Point{r.Max.X / 2, r.Min.Y})
			case s3, ok3 = <-c3:
				go copyImage(newImage, s3.Bounds(), s3, image.Point{r.Min.X, r.Max.Y / 2})
			case s4, ok4 = <-c4:
				go copyImage(newImage, s4.Bounds(), s4, image.Point{r.Max.X / 2, r.Max.Y / 2})
			}
			if ok1 && ok2 && ok3 && ok4 {
				break
			}
		}
		// wait till all copy goroutines are complete
		wg.Wait()
		buf2 := new(bytes.Buffer)
		jpeg.Encode(buf2, newImage, nil)
		c <- base64.StdEncoding.EncodeToString(buf2.Bytes())
	}()
	return c
}
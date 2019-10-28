package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"net/http"
	"strconv"
)

func main()  {
	router := mux.NewRouter()
	router.HandleFunc("/", route)
	router.HandleFunc("/mode", mode)
	router.HandleFunc("/mono", mono)
	router.HandleFunc("/mosaic", mosaic)
	http.ListenAndServe(":8080", router)
}

func mode(w http.ResponseWriter, r *http.Request)  {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error: ", err)
	}
	if r.Form["mode"][0] == "mosaic" {
		t, _ := template.ParseFiles("studio1.html")
		t.Execute(w, nil)
	} else if r.Form["mode"][0] == "mono" {
		t, _ := template.ParseFiles("studio2.html")
		t.Execute(w, nil)
	} else {
		t, _ := template.ParseFiles("home.html")
		t.Execute(w, nil)
	}
}

func route(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("home.html")
	t.Execute(w, nil)
}

func mosaic(w http.ResponseWriter, r *http.Request) {
	tileSize, _ := strconv.Atoi(r.FormValue("tile_size"))
	r.ParseMultipartForm(10485760)
	file, _, _ := r.FormFile("image")
	defer file.Close()
	original, _, _ := image.Decode(file)
	bounds := original.Bounds()
	newImage := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y <= bounds.Max.Y; y += tileSize {
		for x := bounds.Min.X; x <= bounds.Max.X; x += tileSize {
			tile := image.NewRGBA(image.Rect(x, y, x + tileSize, y + tileSize))
			tileBounds := tile.Bounds()

			red, green, blue := 0, 0, 0
			for j := y; j <= y +tileSize; j++ {
				for i := x; i <= x + tileSize; i++ {
					r, g, b, _ := original.At(i, j).RGBA()
					red += int(r)
					green += int(g)
					blue += int(b)
				}
			}
			redAV := uint8(red / tileSize)
			greenAV := uint8(green / tileSize)
			blueAV := uint8(blue / tileSize)

			// 平均色で塗りつぶす
			for j := tileBounds.Min.Y; j <= tileBounds.Max.Y; j++ {
				for i := tileBounds.Min.X; i <= tileBounds.Max.X; i++ {
					a := color.RGBAModel.Convert(original.At(i, j)).(color.RGBA).A
					tile.Set(i, j, color.RGBA{R: redAV, G: greenAV, B: blueAV, A: a})
				}
			}
			//fmt.Println(tile)

			t := tile.SubImage(tile.Bounds())
			tileBounds = image.Rect(x, y, x + tileSize, y + tileSize)
			//fmt.Println(t)
			//fmt.Println(tileBounds)
			draw.Draw(newImage, tileBounds, t, image.Point{x,y}, draw.Src)
		}
	}
	fmt.Println(newImage)


	// 元の画像
	buf1 := new(bytes.Buffer)
	jpeg.Encode(buf1, original, nil)
	originalStr := base64.StdEncoding.EncodeToString(buf1.Bytes())

	// 加工後
	buf2 := new(bytes.Buffer)
	jpeg.Encode(buf2, newImage, nil)
	newStr := base64.StdEncoding.EncodeToString(buf2.Bytes())

	images := map[string]string{
		"original": originalStr,
		"new":   newStr,
	}

	t, _ := template.ParseFiles("result1.html")
	t.Execute(w, images)
}

func mono(w http.ResponseWriter, r *http.Request) {
	// フォームから画像を受け取る
	r.ParseMultipartForm(10485760)
	file, _, _ := r.FormFile("image")
	defer file.Close()
	original, _, _ := image.Decode(file)
	bounds := original.Bounds()
	newImage := image.NewGray16(bounds)
	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		for x := bounds.Min.X; x <= bounds.Max.X; x++ {
			gray, _ := color.Gray16Model.Convert(original.At(x, y)).(color.Gray16)
			newImage.Set(x, y, gray)
		}
	}

	// 元の画像
	buf1 := new(bytes.Buffer)
	jpeg.Encode(buf1, original, nil)
	originalStr := base64.StdEncoding.EncodeToString(buf1.Bytes())

	// 加工後
	buf2 := new(bytes.Buffer)
	jpeg.Encode(buf2, newImage, nil)
	newStr := base64.StdEncoding.EncodeToString(buf2.Bytes())

	images := map[string]string{
		"original": originalStr,
		"new":   newStr,
	}

	t, _ := template.ParseFiles("result2.html")
	t.Execute(w, images)
}
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
)

func Mono(w http.ResponseWriter, r *http.Request) {
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

	t, _ := template.ParseFiles("templates/result2.html")
	t.Execute(w, images)
}
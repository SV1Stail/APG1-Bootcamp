package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func main() {
	width := 300
	height := 300
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	for x := 0; x < 300; x++ {
		for y := 0; y < 300; y++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
			img.Set(x+200, y+200, color.RGBA{255, 255, 255, 255})
			img.Set(x+100, y+100, color.RGBA{255, 255, 255, 255})
			img.Set(x, y+200, color.RGBA{255, 255, 255, 255})
			img.Set(x+200, y, color.RGBA{255, 255, 255, 255})

		}
	}
	radius := 100
	for i := 0; i < 3; i++ {
		for angle := 0.0; angle < 2*math.Pi; angle += 0.01 {
			x := width/2 + int(float64(radius)*math.Cos(angle))
			y := height/2 + int(float64(radius)*math.Sin(angle))
			img.Set(x, y, color.RGBA{255, 0, 0, 255})

			radius -= 50
			x = width/2 + int(float64(radius)*math.Cos(angle))
			y = height/2 + int(float64(radius)*math.Sin(angle))
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
			radius += 100
			x = width/2 + int(float64(radius)*math.Cos(angle))
			y = height/2 + int(float64(radius)*math.Sin(angle))
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
			radius -= 50
		}
		radius -= 1
	}

	file, err := os.Create("amazing_logo.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Кодируем изображение в PNG формат и записываем в файл
	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}

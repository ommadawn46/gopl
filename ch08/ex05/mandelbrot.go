package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
)

type ImageParam struct {
	xmin, ymin, xmax, ymax float64
	width, height int
}

type Result struct {
	x, y int
	c color.Color
}

func main() {
	img := genImageInParallel(ImageParam{-2, -2, +2, +2, 1024, 1024}, 16)
	png.Encode(os.Stdout, img)
}

func genImageInParallel(p ImageParam, gNum int) *image.RGBA {
	results := make(chan Result, 1024)
	img := image.NewRGBA(image.Rect(0, 0, p.width, p.height))
	for g := 0; g < gNum; g++ {
		go func(g int) {
			i := g
			for i < p.width * p.height {
				px, py := i % p.width, i / p.width
				y := float64(py)/float64(p.height)*(p.ymax-p.ymin) + p.ymin
				x := float64(px)/float64(p.width)*(p.xmax-p.xmin) + p.xmin
				z := complex(x, y)
				results <- Result{px, py, mandelbrot(z)}
				i += gNum
			}
		}(g)
	}
	for i := 0; i < p.width * p.height; i++ {
		r := <-results
		img.Set(r.x, r.y, r.c)
	}
	return img
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return nToColor(contrast * n)
		}
	}
	return color.Black
}

func nToColor(n uint8) color.Color {
	var r, g, b uint8
	quot := 0xff / 4.0
	if n <= uint8(quot) {
		// Red To Yellow
		c := math.Max(float64(n), 0) / quot
		r = 0xff
		g = uint8(math.Min(math.Max(255*c, 0), 255))
		b = 0x0
	} else if n <= uint8(quot*2) {
		// Yellow To Blue-Green
		c := (float64(n) - quot) / quot
		r = uint8(math.Min(math.Max(255*(1.0-c), 0), 255))
		g = 0xff
		b = uint8(math.Min(math.Max(255*c, 0), 255))
	} else if n <= uint8(quot*3) {
		// Blue-Green To Blue
		c := (float64(n) - quot*2) / quot
		r = 0
		g = uint8(math.Min(math.Max(255*(1.0-c), 0), 255))
		b = 0xff
	} else {
		// Blue To Black
		c := (math.Min(float64(n), 0xff) - quot*3) / quot
		r = 0
		g = 0
		b = uint8(math.Min(math.Max(255*(1.0-c), 0), 255))
	}
	return color.RGBA{r, g, b, 0xff}
}

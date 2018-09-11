package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
)

func main() {
	const (
		ssRate                 = 5
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
		sWidth, sHeight        = width * ssRate, height * ssRate
	)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < sWidth; py += ssRate {
		for px := 0; px < sHeight; px += ssRate {
			var subPixels []color.Color
			for py_ := py; py_ < py+ssRate; py_++ {
				for px_ := px; px_ < px+ssRate; px_++ {
					y := float64(py_)/sWidth*(ymax-ymin) + ymin
					x := float64(px_)/sHeight*(xmax-xmin) + xmin
					z := complex(x, y)
					subPixels = append(subPixels, mandelbrot(z))
				}
			}
			img.Set(px/ssRate, py/ssRate, colorsAverage(subPixels))
		}
	}
	png.Encode(os.Stdout, img)
}

func colorsAverage(colors []color.Color) color.Color {
	var r, g, b, a int
	for _, c := range colors {
		r_, g_, b_, a_ := c.RGBA()
		r += int(r_)
		g += int(g_)
		b += int(b_)
		a += int(a_)
	}
	n := len(colors)
	return color.RGBA64{uint16(r / n), uint16(g / n), uint16(b / n), uint16(a / n)}
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

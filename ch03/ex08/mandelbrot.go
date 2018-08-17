package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/big"
	"math/cmplx"
	"os"
)

func main() {
	const (
		xmin, ymin = 0.3700028, 0.3700022
		scale = 0.0000005
		xmax, ymax = xmin+scale, ymin+scale
		width, height          = 1024, 1024
	)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			img.Set(px, py, mandelbrotFloat(z))
		}
	}
	png.Encode(os.Stdout, img)
}

const iterations = 200
const contrast = 15

func mandelbrot64(z complex128) color.Color {
	var v complex64
	z_ := complex64(z)
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z_
		if cmplx.Abs(complex128(v)) > 2 {
			return nToColor(contrast * n)
		}
	}
	return color.Black
}

func mandelbrot128(z complex128) color.Color {
	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return nToColor(contrast * n)
		}
	}
	return color.Black
}

func mandelbrotFloat(z complex128) color.Color {
	vr, vi := big.NewFloat(0), big.NewFloat(0)
	zr, zi := big.NewFloat(real(z)), big.NewFloat(imag(z))
	for n := uint8(0); n < iterations; n++ {
		vr_, vi_ := &big.Float{}, &big.Float{}
		// vr = vr*vr - vi*vi + zr
		vr_.Mul(vr, vr).Sub(vr_, vi_.Mul(vi, vi)).Add(vr_, zr)
		// vi = 2*vr*vi + zi
		vi_.Mul(vi, vr).Mul(vi_, big.NewFloat(2)).Add(vi_, zi)
		vr, vi = vr_, vi_
		// sqSum = vr*vr + vi*vi
		sqSum := &big.Float{}
		sqSum = sqSum.Mul(vr, vr).Add(sqSum, (&big.Float{}).Mul(vi, vi))
		if sqSum.Cmp(big.NewFloat(4)) > 0 {
			return nToColor(contrast * n)
		}
	}
	return color.Black
}

func mandelbrotRat(z complex128) color.Color {
	vr, vi := big.NewRat(0, 1), big.NewRat(0, 1)
	zr, zi := new(big.Rat).SetFloat64(real(z)), new(big.Rat).SetFloat64(imag(z))
	for n := uint8(0); n < iterations; n++ {
		vr_, vi_ := &big.Rat{}, &big.Rat{}
		// vr = vr*vr - vi*vi + zr
		vr_.Mul(vr, vr).Sub(vr_, vi_.Mul(vi, vi)).Add(vr_, zr)
		// vi = 2*vr*vi + zi
		vi_.Mul(vi, vr).Mul(vi_, big.NewRat(2, 1)).Add(vi_, zi)
		vr, vi = vr_, vi_
		// sqSum = vr*vr + vi*vi
		sqSum := &big.Rat{}
		sqSum = sqSum.Mul(vr, vr).Add(sqSum, (&big.Rat{}).Mul(vi, vi))
		if sqSum.Cmp(big.NewRat(4, 1)) > 0 {
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

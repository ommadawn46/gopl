package main

import (
	"image/color"
	"testing"
)

var temp color.Color

func BenchmarkComplex64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		temp = mandelbrot64(complex(float64(i), float64(i)))
	}
}

func BenchmarkComplex128(b *testing.B) {
	for i := 0; i < b.N; i++ {
		temp = mandelbrot128(complex(float64(i), float64(i)))
	}
}

func BenchmarkBigFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		temp = mandelbrotFloat(complex(float64(i), float64(i)))
	}
}

func BenchmarkBigRat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		temp = mandelbrotRat(complex(float64(i), float64(i)))
	}
}

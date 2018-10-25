package main

import (
	"image"
	"testing"
)

var temp image.RGBA

func Benchmark1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := float64(i)
		temp = *genImageInParallel(ImageParam{-f, -f, f, f, 1024, 1024}, 1)
	}
}

func Benchmark2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := float64(i)
		temp = *genImageInParallel(ImageParam{-f, -f, f, f, 1024, 1024}, 2)
	}
}

func Benchmark4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := float64(i)
		temp = *genImageInParallel(ImageParam{-f, -f, f, f, 1024, 1024}, 4)
	}
}

func Benchmark8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := float64(i)
		temp = *genImageInParallel(ImageParam{-f, -f, f, f, 1024, 1024}, 8)
	}
}

func Benchmark16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := float64(i)
		temp = *genImageInParallel(ImageParam{-f, -f, f, f, 1024, 1024}, 16)
	}
}

func Benchmark32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := float64(i)
		temp = *genImageInParallel(ImageParam{-f, -f, f, f, 1024, 1024}, 32)
	}
}

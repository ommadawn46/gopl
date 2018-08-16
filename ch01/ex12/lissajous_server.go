package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var palette = []color.Color{color.Black}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < 0xff; i++ {
		palette = append(palette, color.RGBA{
			gradation(i), gradation(i * 2), gradation(i * 4), 0xff,
		})
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cycles := 18
		res := 0.001
		size := 300
		nframes := 64
		delay := 8

		if err := r.ParseForm(); err != nil {
			log.Print(err)
		} else {
			if params, ok := r.Form["cycles"]; ok {
				if val, err := strconv.Atoi(params[0]); err == nil {
					cycles = val
				}
			}
			if params, ok := r.Form["res"]; ok {
				if val, err := strconv.ParseFloat(params[0], 64); err == nil {
					res = val
				}
			}
			if params, ok := r.Form["size"]; ok {
				if val, err := strconv.Atoi(params[0]); err == nil {
					size = val
				}
			}
			if params, ok := r.Form["nframes"]; ok {
				if val, err := strconv.Atoi(params[0]); err == nil {
					nframes = val
				}
			}
			if params, ok := r.Form["delay"]; ok {
				if val, err := strconv.Atoi(params[0]); err == nil {
					delay = val
				}
			}
		}
		lissajous(w, cycles, res, size, nframes, delay)
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func gradation(n int) uint8 {
	return uint8(math.Abs(float64(0x100 - n%0x200)))
}

func lissajous(out io.Writer, cycles int, res float64, size int, nframes int, delay int) {
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(
				size+int(x*float64(size)+0.5),
				size+int(y*float64(size)+0.5),
				gradation(int(t/res)),
			)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}

package main

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"math/cmplx"
	"net/http"
	"strconv"
)

var colors = [][]float64{
	{0, 1, 0}, {1, 0, 0}, {0, 0, 1}, {1, 1, 0}, {0, 1, 1}, {1, 0, 1},
	{1, 0.5, 0}, {0, 1, 0.5}, {1, 0, 0.5}, {0.5, 1, 0}, {0, 0.5, 1}, {0.5, 0, 1},
	{1, 0.33, 0.66}, {0.33, 1, 0.66}, {0.33, 0.66, 1}, {0.66, 0.33, 0.33}, {0.33, 0.66, 0.33}, {0.33, 0.33, 0.66},
	{1, 0.33, 0.33}, {0.33, 1, 0.33}, {0.33, 0.33, 1}, {1, 0.66, 0.33}, {0.66, 1, 0.33}, {0.66, 0.33, 1},
	{0.75, 0.50, 0.25}, {0.50, 0.75, 0.25}, {0.50, 0.25, 0.75}, {0.75, 0.25, 0.50}, {0.25, 0.75, 0.50}, {0.25, 0.50, 0.75},
	{0.75, 0, 0.25}, {0, 0.75, 0.25}, {0, 0.25, 0.75}, {0.75, 0.25, 0}, {0.25, 0.75, 0}, {0.25, 0, 0.75},
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var (
			width, height = 1024, 1024
			cx, cy        = 0.0, 0.0
			scale         = 10.0
			contrast      = uint8(5)
			d             = 4
		)
		if err := r.ParseForm(); err != nil {
			log.Print(err)
		} else {
			if params, ok := r.Form["d"]; ok {
				if val, err := strconv.Atoi(params[0]); err == nil {
					if 0 < val && val <= 36 {
						d = val
					}
				}
			}
			if params, ok := r.Form["x"]; ok {
				if val, err := strconv.ParseFloat(params[0], 64); err == nil {
					cx = val
				}
			}
			if params, ok := r.Form["y"]; ok {
				if val, err := strconv.ParseFloat(params[0], 64); err == nil {
					cy = val
				}
			}
			if params, ok := r.Form["scale"]; ok {
				if val, err := strconv.ParseFloat(params[0], 64); err == nil {
					if 0 < val {
						scale = val
					}
				}
			}
			if params, ok := r.Form["contrast"]; ok {
				if val, err := strconv.Atoi(params[0]); err == nil {
					if 0 < val {
						contrast = uint8(val)
					}
				}
			}
		}
		renderFractal(w, width, height, d, cx, cy, scale, contrast)
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func renderFractal(out io.Writer, width int, height int,
	d int, cx float64, cy float64, scale float64, contrast uint8) {
	f := func(x complex128) complex128 { return fnD(x, d) }
	roots := findRoots(f, -scale, -scale, scale, scale)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := 2*float64(py)*scale/float64(height) - scale
		for px := 0; px < width; px++ {
			x := 2*float64(px)*scale/float64(width) - scale
			z := complex(cx+x, cy+y)
			img.Set(px, py, calcColor(f, roots, z, contrast))
		}
	}
	png.Encode(out, img)
}

func fnD(x complex128, n int) complex128 {
	s := 1 + 0i
	for i := 0; i < n; i++ {
		s *= x
	}
	return s - 1
}

func findRoots(f func(complex128) complex128, xmin, ymin, xmax, ymax float64) []complex128 {
	var roots []complex128
	const threshold = 1.0 / (1 << 8)

	appendIfNotExists := func(r complex128) {
		for _, r_ := range roots {
			if calcComplexDist(r, r_) < threshold {
				return
			}
		}
		roots = append(roots, r)
	}

	for x := xmin; x <= xmax; x += (xmax - xmin) / 0xff {
		for y := ymin; y <= ymax; y += (ymax - ymin) / 0xff {
			z := complex(x, y)
			if r, _, err := solveWithNewton(f, z); err == nil {
				appendIfNotExists(r)
			}
		}
	}
	return roots
}

func calcColor(f func(complex128) complex128, roots []complex128, z complex128, contrast uint8) color.Color {
	if r, n, err := solveWithNewton(f, z); err == nil {
		return getColorFromN(int(n)*int(contrast), searchNearRootIndex(r, roots))
	}
	return color.Black
}

func getColorFromN(n int, rootIndex int) color.Color {
	baseColor := colors[rootIndex]
	brightness := 255 - math.Min(math.Max(float64(n), 0), 255)

	r := uint8(baseColor[0] * brightness)
	g := uint8(baseColor[1] * brightness)
	b := uint8(baseColor[2] * brightness)

	return color.RGBA{r, g, b, 0xff}
}

func solveWithNewton(f func(complex128) complex128, z complex128) (complex128, uint8, error) {
	const iterations = 0xff
	const threshold = 1.0 / (1 << 16)
	v := z
	for n := uint8(0); n < iterations; n++ {
		v -= f(v) / df(f, v)
		if cmplx.Abs(f(v)) < threshold {
			return v, n, nil
		}
	}
	return cmplx.NaN(), 0, errors.New("No solution found.")
}

func df(f func(complex128) complex128, x complex128) complex128 {
	const epsilon = 1.0 / (1 << 32)
	return (f(x+epsilon) - f(x)) / epsilon
}

func searchNearRootIndex(x complex128, roots []complex128) int {
	nearRootIndex := 0
	nearDist := calcComplexDist(roots[0], x)
	for i := 1; i < len(roots); i++ {
		dist := calcComplexDist(roots[i], x)
		if dist < nearDist {
			nearRootIndex, nearDist = i, dist
		}
	}
	return nearRootIndex
}

func calcComplexDist(a, b complex128) float64 {
	real_d := real(a) - real(b)
	imag_d := imag(a) - imag(b)
	return real_d*real_d + imag_d*imag_d
}

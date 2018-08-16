package main

import (
	"errors"
	"fmt"
	"math"
)

const (
	width, height = 600, 320
	cells         = 100
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func main() {
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, z0, err := corner(i+1, j)
			if err != nil {
				continue
			}
			bx, by, z1, err := corner(i, j)
			if err != nil {
				continue
			}
			cx, cy, z2, err := corner(i, j+1)
			if err != nil {
				continue
			}
			dx, dy, z3, err := corner(i+1, j+1)
			if err != nil {
				continue
			}
			r, g, b := zToRGB((z0 + z1 + z2 + z3) / 4)
			fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g' "+
				"style='stroke: #%02x%02x%02x'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy,
				r, g, b)
		}
	}
	fmt.Println("</svg>")
}

func zToRGB(z float64) (int, int, int) {
	c := (z + 0.2) / 1.1
	r := int(math.Min(math.Max(255*c, 0), 255))
	g := 0
	b := int(math.Min(math.Max(255*(1.0-c), 0), 255))
	return r, g, b
}

func corner(i, j int) (float64, float64, float64, error) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	z := f(x, y)
	if math.IsNaN(z) || math.IsInf(z, 0) {
		return 0, 0, 0, errors.New("Invalid Value")
	}

	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, z, nil
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}

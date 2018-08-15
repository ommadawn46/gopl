package main

import (
  "image"
  "image/color"
  "image/png"
  "math"
  "math/cmplx"
  "os"
  "errors"
  "flag"
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
  d := flag.Int("d", 4, "It specifies the order of equation.")
  flag.Parse()

  const (
    xmin, ymin, xmax, ymax = -10, -10, +10, +10
    width, height = 1024, 1024
  )
  f := func(x complex128) complex128 { return fnD(x, *d) }
  roots := findRoots(f, xmin, ymin, xmax, ymax)

  img := image.NewRGBA(image.Rect(0, 0, width, height))
  for py := 0; py < height; py++ {
    y := float64(py) / height*(ymax-ymin) + ymin
    for px := 0; px < width; px++ {
      x := float64(px) / width*(xmax-xmin) + xmin
      z := complex(x, y)
      img.Set(px, py, calcColor(f, roots, z))
    }
  }
  png.Encode(os.Stdout, img)
}

func fnD(x complex128, n int) complex128 {
  s := 1+0i
  for i := 0; i < n; i++ {
    s *= x
  }
  return s - 1
}

func findRoots(f func(complex128)complex128, xmin, ymin, xmax, ymax float64) []complex128 {
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

func calcColor(f func(complex128)complex128, roots []complex128, z complex128) color.Color {
  const contrast = 15
  if r, n, err := solveWithNewton(f, z); err == nil {
    return getColorFromN(n * contrast, searchNearRootIndex(r, roots))
  }
  return color.Black
}

func getColorFromN(n uint8, rootIndex int) color.Color {
  baseColor := colors[rootIndex]
  brightness := 255 - math.Min(math.Max(float64(n), 0), 255)

  r := uint8(baseColor[0]*brightness)
  g := uint8(baseColor[1]*brightness)
  b := uint8(baseColor[2]*brightness)

  return color.RGBA{r, g, b, 0xff}
}

func solveWithNewton(f func(complex128)complex128, z complex128) (complex128, uint8, error) {
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

func df(f func(complex128)complex128, x complex128) complex128 {
  const epsilon = 1.0 / (1 << 32)
  return (f(x+epsilon) - f(x)) / epsilon
}

func searchNearRootIndex(x complex128, roots []complex128) int {
  nearRootIndex := 0
  nearDist := calcComplexDist(roots[0], x)
  for i := 0; i < len(roots); i++ {
    dist := calcComplexDist(roots[i], x)
    if dist < nearDist {
      nearRootIndex, nearDist = i, dist
    }
  }
  return nearRootIndex
}

func calcComplexDist(a, b complex128) float64 {
  real_d := real(a)-real(b)
  imag_d := imag(a)-imag(b)
  return real_d*real_d + imag_d*imag_d
}

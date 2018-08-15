package main

import (
  "io"
  "math"
  "log"
  "net/http"
  "strconv"
  "fmt"
  "errors"
)

func main() {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
    var (
      width, height = 600, 320
      cells = 100
      xyrange = 30.0
      xyscale = float64(width) / 2 / xyrange
      zscale = float64(height) * 0.4
      angle = math.Pi / 6
      lowColor = "B"
      highColor = "R"
    )
    if err := r.ParseForm(); err != nil {
      log.Print(err)
    } else {
      if params, ok := r.Form["width"]; ok {
        if val, err := strconv.Atoi(params[0]); err == nil {
          width = val
        }
      }
      if params, ok := r.Form["height"]; ok {
        if val, err := strconv.Atoi(params[0]); err == nil {
          height = val
        }
      }
      if params, ok := r.Form["cells"]; ok {
        if val, err := strconv.Atoi(params[0]); err == nil {
          cells = val
        }
      }
      if params, ok := r.Form["xyrange"]; ok {
        if val, err := strconv.ParseFloat(params[0], 64); err == nil {
          xyrange = val
        }
      }
      if params, ok := r.Form["xyscale"]; ok {
        if val, err := strconv.ParseFloat(params[0], 64); err == nil {
          xyscale = float64(width) / val / xyrange
        }
      }
      if params, ok := r.Form["zscale"]; ok {
        if val, err := strconv.ParseFloat(params[0], 64); err == nil {
          zscale = float64(height) * val
        }
      }
      if params, ok := r.Form["angle"]; ok {
        if val, err := strconv.ParseFloat(params[0], 64); err == nil {
          angle = val
        }
      }
      if params, ok := r.Form["lowColor"]; ok {
        lowColor = params[0]
      }
      if params, ok := r.Form["highColor"]; ok {
        highColor = params[0]
      }
    }
    w.Header().Set("Content-Type", "image/svg+xml")
    surfaceRender(w, width, height, cells, xyrange,
      xyscale, zscale, angle, lowColor, highColor)
  })
  log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func surfaceRender(out io.Writer, width int, height int,
  cells int, xyrange float64, xyscale float64, zscale float64,
  angle float64, lowColor string, highColor string) {
  fmt.Fprintf(out, "<svg xmlns='http://www.w3.org/2000/svg' " +
    "style='stroke: grey; fill: white; stroke-width: 0.7' " +
    "width='%d' height='%d'>", width, height)

  corner := func (i, j int) (float64, float64, float64, error) {
    x := xyrange * (float64(i)/float64(cells) - 0.5)
    y := xyrange * (float64(j)/float64(cells) - 0.5)
    z := f(x, y)
    if math.IsNaN(z) || math.IsInf(z, 0) {
      return 0, 0, 0, errors.New("Invalid Value")
    }

    sx := float64(width)/2 + (x-y)*math.Cos(angle)*xyscale
    sy := float64(height)/2 + (x+y)*math.Sin(angle)*xyscale - z*zscale
    return sx, sy, z, nil
  }

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
        r, g, b := zToRGB((z0 + z1 + z2 + z3) / 4, lowColor, highColor)
        fmt.Fprintf(out, "<polygon points='%g,%g %g,%g %g,%g %g,%g' " +
          "style='stroke: #%02x%02x%02x'/>\n",
          ax, ay, bx, by, cx, cy, dx, dy,
          r, g, b)
    }
  }
  fmt.Fprintln(out, "</svg>")
}

func zToRGB(z float64, lowColor string, highColor string) (int, int, int) {
  c := (z + 0.2) / 1.1
  r, g, b := 0, 0, 0
  switch lowColor {
  case "R": r = int(math.Min(math.Max(255*(1.0-c), 0), 255))
  case "G": g = int(math.Min(math.Max(255*(1.0-c), 0), 255))
  case "B": b = int(math.Min(math.Max(255*(1.0-c), 0), 255))
  }
  switch highColor {
  case "R": r = int(math.Min(math.Max(255*c, 0), 255))
  case "G": g = int(math.Min(math.Max(255*c, 0), 255))
  case "B": b = int(math.Min(math.Max(255*c, 0), 255))
  }
  return r, g, b
}

func f(x, y float64) float64 {
  r := math.Hypot(x, y)
  return math.Sin(r) / r
}

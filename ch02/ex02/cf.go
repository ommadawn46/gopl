package main

import (
  "fmt"
  "os"
  "bufio"
  "strconv"
  "flag"

  "./lenconv"
  "./tempconv"
  "./weightconv"
)

var (
  l = flag.Bool("l", false, "convert value as length")
  t = flag.Bool("t", false, "convert value as temperature")
  w = flag.Bool("w", false, "convert value as weight")
)

func printValueConv(opt string, val float64){
  switch opt {
  case "l":
    m := lenconv.Meter(val)
    ft := lenconv.Feet(val)
    fmt.Printf("%s = %s, %s = %s\n",
      m, lenconv.MToFt(m), ft, lenconv.FtToM(ft))
  case "t":
    f := tempconv.Fahrenheit(val)
    c := tempconv.Celsius(val)
    fmt.Printf("%s = %s, %s = %s\n",
      f, tempconv.FToC(f), c, tempconv.CToF(c))
  case "w":
    kg := weightconv.KiloGram(val)
    p := weightconv.Pound(val)
    fmt.Printf("%s = %s, %s = %s\n",
      kg, weightconv.KgToP(kg), p, weightconv.PToKg(p))
  }
}

func main(){
  flag.Parse()

  if flag.NArg() > 1 {
    for _, arg := range flag.Args() {
      val, err := strconv.ParseFloat(arg, 64)
      if err != nil {
        fmt.Fprintf(os.Stderr, "cf: %v\n", err)
        os.Exit(1)
      }
      opt := ""
      switch {
        case *l: opt = "l"
        case *t: opt = "t"
        case *w: opt = "w"
      }
      printValueConv(opt, val)
    }
  } else {
    input := bufio.NewScanner(os.Stdin)
    fmt.Print("Option > ")
    opt := ""
    for input.Scan() {
      opt = input.Text()
      break
    }
    fmt.Print("Value > ")
    val := 0.0
    for input.Scan() {
      val, _ = strconv.ParseFloat(input.Text(), 64)
      break
    }
    printValueConv(opt, val)
  }
}

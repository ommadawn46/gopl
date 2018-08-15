package main

import (
  "crypto/sha256"
  "crypto/sha512"
  "fmt"
  "flag"
)

var (
  algo = flag.String("a", "sha256", "Hash Algorithm")
  message = flag.String("m", "", "Message")
)

func main() {
  flag.Parse()

  m := []byte(*message)
  x := ""
  switch *algo {
  case "sha384":
    x = fmt.Sprintf("%x", sha512.Sum384(m))
  case "sha512":
    x = fmt.Sprintf("%x", sha512.Sum512(m))
  default:
    x = fmt.Sprintf("%x", sha256.Sum256(m))
  }

  fmt.Println(x)
}

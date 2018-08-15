package main

import (
  "crypto/sha256"
  "encoding/binary"
  "fmt"
)

func PopCount(x uint64) int {
  var count int
  for x > 0 {
    x &= (x-1)
    count ++
  }
  return count
}

func main() {
  c1 := sha256.Sum256([]byte("ABCDEFGH"))
  c2 := sha256.Sum256([]byte("abcdefgh"))
  x1 := binary.BigEndian.Uint64(c1[:])
  x2 := binary.BigEndian.Uint64(c2[:])
  fmt.Println(PopCount(x1 ^ x2))
}

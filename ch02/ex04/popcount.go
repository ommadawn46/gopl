package popcount

var pc [256]byte

func init() {
  for i := range pc {
    pc[i] = pc[i/2] + byte(i&1)
  }
}

func PopCountTable(x uint64) int {
  return int(
    pc[byte(x>>(0*8))] +
    pc[byte(x>>(1*8))] +
    pc[byte(x>>(2*8))] +
    pc[byte(x>>(3*8))] +
    pc[byte(x>>(4*8))] +
    pc[byte(x>>(5*8))] +
    pc[byte(x>>(6*8))] +
    pc[byte(x>>(7*8))])
}

func PopCountBitShift(x uint64) int {
  var count int
  for i := 0; i < 64; i++ {
    count += int(x & 1)
    x >>= 1
  }
  return count
}

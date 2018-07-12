package popcount

import (
  "testing"
)

func BenchmarkPopCountTable(b *testing.B) {
  for i := uint64(0); i < uint64(b.N); i++ {
    PopCountTable(i)
  }
}

func BenchmarkPopCountLowerBitClear(b *testing.B) {
  for i := uint64(0); i < uint64(b.N); i++ {
    PopCountLowerBitClear(i)
  }
}

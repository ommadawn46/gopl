package popcount

import (
	"testing"
)

var temp int

// init()で予め表を初期化しておくと表に基づく方法が圧倒的に早くなってしまう
// そのため毎回初期化するようにしている
func benchmarkPopCountTable(b *testing.B, n int) {
	for i := uint64(0); i < uint64(b.N); i++ {
		for j := range pc {
			pc[j] = pc[j/2] + byte(j&1)
		}
		temp = 0
		for j := 0; j < n; j++ {
			temp += PopCountTable(i)
		}
	}
}

func benchmarkPopCount(b *testing.B, f func(x uint64) int, n int) {
	for i := uint64(0); i < uint64(b.N); i++ {
		temp = 0
		for j := 0; j < n; j++ {
			temp += f(i)
		}
	}
}

func BenchmarkPopCountTable1(b *testing.B) {
	benchmarkPopCountTable(b, 1)
}

func BenchmarkPopCountTable10(b *testing.B) {
	benchmarkPopCountTable(b, 10)
}

func BenchmarkPopCountTable100(b *testing.B) {
	benchmarkPopCountTable(b, 100)
}

func BenchmarkPopCountTable1000(b *testing.B) {
	benchmarkPopCountTable(b, 1000)
}

func BenchmarkPopCountTable10000(b *testing.B) {
	benchmarkPopCountTable(b, 10000)
}

func BenchmarkPopCountBitShift1(b *testing.B) {
	benchmarkPopCount(b, PopCountBitShift, 1)
}

func BenchmarkPopCountBitShift10(b *testing.B) {
	benchmarkPopCount(b, PopCountBitShift, 10)
}

func BenchmarkPopCountBitShift100(b *testing.B) {
	benchmarkPopCount(b, PopCountBitShift, 100)
}

func BenchmarkPopCountBitShift1000(b *testing.B) {
	benchmarkPopCount(b, PopCountBitShift, 1000)
}

func BenchmarkPopCountBitShift10000(b *testing.B) {
	benchmarkPopCount(b, PopCountBitShift, 10000)
}

func BenchmarkPopCountLowerBitClear1(b *testing.B) {
	benchmarkPopCount(b, PopCountLowerBitClear, 1)
}

func BenchmarkPopCountLowerBitClear10(b *testing.B) {
	benchmarkPopCount(b, PopCountLowerBitClear, 10)
}

func BenchmarkPopCountLowerBitClear100(b *testing.B) {
	benchmarkPopCount(b, PopCountLowerBitClear, 100)
}

func BenchmarkPopCountLowerBitClear1000(b *testing.B) {
	benchmarkPopCount(b, PopCountLowerBitClear, 1000)
}

func BenchmarkPopCountLowerBitClear10000(b *testing.B) {
	benchmarkPopCount(b, PopCountLowerBitClear, 10000)
}

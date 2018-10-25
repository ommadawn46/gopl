package popcount

import "sync"

var pc [256]byte

var initTableOnce sync.Once

func initTable() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func PopCountTable(x uint64) int {
	initTableOnce.Do(initTable)
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

func PopCountLoop(x uint64) int {
	initTableOnce.Do(initTable)
	var count int
	for i := 0; i < 8; i++ {
		count += int(pc[byte(x>>uint64(i*8))])
	}
	return count
}

func PopCountBitShift(x uint64) int {
	var count int
	for i := 0; i < 64; i++ {
		count += int(x & 1)
		x >>= 1
	}
	return count
}

func PopCountLowerBitClear(x uint64) int {
	var count int
	for x > 0 {
		x &= (x - 1)
		count++
	}
	return count
}

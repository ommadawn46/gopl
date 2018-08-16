package popcount

var pc [256]byte
var pc_loop [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
	for i := range pc_loop {
		pc_loop[i] = pc_loop[i/2] + byte(i&1)
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

func PopCountLoop(x uint64) int {
	var count int
	for i := 0; i < 8; i++ {
		count += int(pc_loop[byte(x>>uint64(i*8))])
	}
	return count
}

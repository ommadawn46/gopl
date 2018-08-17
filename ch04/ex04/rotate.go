package rotate

func rotate(s []int, r int) {
	n := len(s)
	if n <= 0 {
		return
	}
	for r < 0 {
		r += n
	}
	r %= n
	t := make([]int, r)
	copy(t, s[:r])
	copy(s, s[r:])
	copy(s[n-r:], t)
}

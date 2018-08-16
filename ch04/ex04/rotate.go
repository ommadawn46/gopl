package rotate

func rotate(s []int, r int) {
	if len(s) <= 0 {
		return
	}
	r = r % len(s)
	t := make([]int, r)
	copy(t, s[:r])
	copy(s, s[r:])
	copy(s[len(s)-r:], t)
}

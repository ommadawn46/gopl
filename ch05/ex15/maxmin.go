package maxmin

func max(vals ...int) int {
	if len(vals) == 0 {
		return 0
	}
	max := vals[0]
	for _, val := range vals[1:] {
		if val > max {
			max = val
		}
	}
	return max
}

func max_(max int, vals ...int) int {
	if len(vals) == 0 {
		return max
	}
	for _, val := range vals {
		if val > max {
			max = val
		}
	}
	return max
}

func min(vals ...int) int {
	if len(vals) == 0 {
		return 0
	}
	min := vals[0]
	for _, val := range vals[1:] {
		if val < min {
			min = val
		}
	}
	return min
}

func min_(min int, vals ...int) int {
	if len(vals) == 0 {
		return min
	}
	for _, val := range vals {
		if val < min {
			min = val
		}
	}
	return min
}

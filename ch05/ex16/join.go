package join

func join(sep string, args ...string) string {
	if len(args) == 0 {
		return ""
	}
	res := args[0]
	for _, s := range args[1:] {
		res += sep + s
	}
	return res
}

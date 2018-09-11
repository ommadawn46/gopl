package expand

import "regexp"

var pattern = regexp.MustCompile(`\$(\S*)`)

func expand(s string, f func(string) string) string {
	return pattern.ReplaceAllStringFunc(s, func(matched string) string {
		return f(pattern.FindStringSubmatch(matched)[1])
	})
}

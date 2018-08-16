package comma

import (
	"bytes"
	"strings"
)

func comma(s string) string {
	if len(s) <= 0 {
		return ""
	}

	out := new(bytes.Buffer)
	if s[0] == '+' || s[0] == '-' {
		out.WriteString(s[:1])
		s = s[1:]
	}

	interger, decimal := s, ""
	if sp := strings.Split(s, "."); len(sp) >= 2 {
		interger, decimal = sp[0], sp[1]
	}

	strLen := len(interger)
	if strLen%3 > 0 {
		out.WriteString(interger[:strLen%3] + ",")
	}
	sep := ""
	for i := strLen % 3; i < len(interger); i += 3 {
		out.WriteString(sep + interger[i:i+3])
		sep = ","
	}

	if decimal != "" {
		out.WriteString("." + decimal)
	}
	return out.String()
}

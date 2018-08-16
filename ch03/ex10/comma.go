package comma

import (
	"bytes"
)

func comma(s string) string {
	out := new(bytes.Buffer)
	strLen := len(s)
	if strLen%3 > 0 {
		out.WriteString(s[:strLen%3] + ",")
	}
	sep := ""
	for i := strLen % 3; i < len(s); i += 3 {
		out.WriteString(sep + s[i:i+3])
		sep = ","
	}
	return out.String()
}

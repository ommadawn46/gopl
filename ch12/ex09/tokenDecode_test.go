package sexpr

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestTokenDecode(t *testing.T) {
	for _, test := range []struct {
		input string
		want  []Token
	}{
		{
			"",
			[]Token{},
		},
		{
			`("A" "B" "C")`,
			[]Token{StartList{}, String("A"), String("B"), String("C"), EndList{}},
		},
		{
			`("A" "B" "C" (1 2 3))`,
			[]Token{StartList{}, String("A"), String("B"), String("C"), StartList{}, Int(1), Int(2), Int(3), EndList{}, EndList{}},
		},
		{
			`  ( "A"   "B"  "C" (  1  2   3  )   ) `,
			[]Token{StartList{}, String("A"), String("B"), String("C"), StartList{}, Int(1), Int(2), Int(3), EndList{}, EndList{}},
		},
	} {
		r := bytes.NewBuffer([]byte(test.input))
		decoder := NewDecoder(r)
		got := []Token{}

		token, err := decoder.Token()
		for err != io.EOF {
			got = append(got, token)
			token, err = decoder.Token()
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("\ngot: \n%s\nwant: \n%s\n", got, test.want)
		}
	}
}

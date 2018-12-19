package params

import (
	"testing"
)

func TestPack(t *testing.T) {
	type Data struct {
		Labels     []string `http:"l"`
		MaxResults int      `http:"max"`
		Exact      bool     `http:"x"`
	}

	for _, test := range []struct {
		input Data
		want  string
	}{
		{
			Data{[]string{}, 0, false},
			"max=0&x=false",
		},
		{
			Data{[]string{"golang"}, 10, false},
			"l=golang&max=10&x=false",
		},
		{
			Data{[]string{"golang", "programming"}, 100, true},
			"l=golang&l=programming&max=100&x=true",
		},
	} {
		got, err := Pack(&test.input)
		if err != nil {
			t.Fatalf("Pack failed: %v", err)
		}
		if got != test.want {
			t.Errorf("\ngot: \n%s\nwant: \n%s\n", got, test.want)
		}
	}
}

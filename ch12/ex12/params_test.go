package params

import (
	"fmt"
	"net/http"
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

func TestValidate(t *testing.T) {
	type ValidateData struct {
		Name        string `http:"n"`
		Email       string `http:"e" validate:"email"`
		PhoneNumber string `http:"p" validate:"phonenumber"`
	}

	for _, test := range []struct {
		input    string
		wantData ValidateData
		wantErr  error
	}{
		{
			"http://example.com/",
			ValidateData{},
			nil,
		},
		{
			"http://example.com/?n=bob&e=hoge@fuga.piyo&p=090-1234-5678",
			ValidateData{Name: "bob", Email: "hoge@fuga.piyo", PhoneNumber: "090-1234-5678"},
			nil,
		},
		{
			"http://example.com/?n=bob&e=hoge_fuga_piyo&p=090-1234-5678",
			ValidateData{},
			fmt.Errorf(`e: failed to validate "hoge_fuga_piyo" with "email"`),
		},
		{
			"http://example.com/?n=bob&e=hoge@fuga.piyo&p=090@1234@5678",
			ValidateData{},
			fmt.Errorf(`p: failed to validate "090@1234@5678" with "phonenumber"`),
		},
	} {
		var got ValidateData
		req, err := http.NewRequest("GET", test.input, nil)
		if err != nil {
			t.Fatalf("New request failed: %v", err)
		}
		err = Unpack(req, &got)
		if got != test.wantData {
			t.Errorf("\ngot: \n%s\nwant: \n%s\n", got, test.wantData)
		}
		if !(err == nil && test.wantErr == nil ||
			err != nil && test.wantErr != nil && err.Error() == test.wantErr.Error()) {
			t.Errorf("\ngot: \n%s\nwant: \n%s\n", err, test.wantErr)
		}
	}
}

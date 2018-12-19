package sexpr

import (
	"reflect"
	"testing"
)

type TaggedStruct struct {
	StringField string `sexpr:"sf"`
	IntField    int
}

func TestMarhalWithTag(t *testing.T) {
	for _, test := range []struct {
		input TaggedStruct
		want  string
	}{
		{
			TaggedStruct{},
			`()`,
		},
		{
			TaggedStruct{StringField: "abc"},
			`((sf "abc"))`,
		},
		{
			TaggedStruct{"abc", 123},
			`((sf "abc")
 (IntField 123))`,
		},
	} {
		got, err := Marshal(test.input)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		if string(got) != test.want {
			t.Errorf("\ngot: \n%s\nwant: \n%s\n", got, test.want)
		}
	}
}

func TestUnmarhalWithTag(t *testing.T) {
	for _, test := range []struct {
		input string
		want  TaggedStruct
	}{
		{
			`()`,
			TaggedStruct{},
		},
		{
			`((sf "abc"))`,
			TaggedStruct{StringField: "abc"},
		},
		{
			`((sf "abc")
 (IntField 123))`,
			TaggedStruct{"abc", 123},
		},
	} {
		got := TaggedStruct{}
		err := Unmarshal([]byte(test.input), &got)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("\ngot: \n%v\nwant: \n%v\n", got, test.want)
		}
	}
}

func TestCanReverseWithTag(t *testing.T) {
	type Movie struct {
		Title    string            `sexpr:"t"`
		Subtitle string            `sexpr:"st"`
		Year     int               `sexpr:"y"`
		Actor    map[string]string `sexpr:"a"`
		Oscars   []string          `sexpr:"o"`
		Sequel   *string           `sexpr:"s"`
	}
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}

	// Encode it
	data, err := Marshal(strangelove)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	t.Logf("Marshal() = \n%s\n", data)

	// Decode it
	var movie Movie
	if err := Unmarshal(data, &movie); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	t.Logf("Unmarshal() = %+v\n", movie)

	// Check equality.
	if !reflect.DeepEqual(movie, strangelove) {
		t.Fatal("not equal")
	}
}

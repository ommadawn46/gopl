package sexpr

import (
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type wrapper struct {
		Interf1 interface{}
		Interf2 interface{}
		Interf3 interface{}
		Interf4 interface{}
		Interf5 interface{}
	}
	type TestStruct struct {
		StringField    string
		InterfaceField interface{}
	}
	INTERFACES["sexpr.TestStruct"] = reflect.TypeOf(TestStruct{})
	for _, test := range []struct {
		input string
		want  wrapper
	}{
		{`nil`, wrapper{}},
		{`()`, wrapper{}},
		{`(
			 (Interf1 ("[]int" (1 2 3)))
			 (Interf2 ("[]string" ("a" "b" "c")))
			 (Interf3 nil)
       (Interf4 ("[]complex128" (#C(1.2 3.4) #C(5.6 7.8) #C(9.10 11.12))))
			 (Interf5 nil)
			)`,
			wrapper{
				Interf1: []int{1, 2, 3},
				Interf2: []string{"a", "b", "c"},
				Interf3: nil,
				Interf4: []complex128{complex(1.2, 3.4), complex(5.6, 7.8), complex(9.10, 11.12)},
				Interf5: nil,
			},
		},
		{`(
			 (Interf1 ("[]int" (1 2 3)))
			 (Interf2 ("[]string" ("a" "b" "c")))
			 (Interf3 ("[]float64" (1.23 4.56 7.89)))
       (Interf4 ("[]complex128" (#C(1.2 3.4) #C(5.6 7.8) #C(9.10 11.12))))
			 (Interf5 ("sexpr.TestStruct" ((StringField "abc")
                                     (InterfaceField ("sexpr.TestStruct" ((StringField "def")
                                                                          (InterfaceField nil))))))
		  )
		)`,
			wrapper{
				Interf1: []int{1, 2, 3},
				Interf2: []string{"a", "b", "c"},
				Interf3: []float64{1.23, 4.56, 7.89},
				Interf4: []complex128{complex(1.2, 3.4), complex(5.6, 7.8), complex(9.10, 11.12)},
				Interf5: TestStruct{StringField: "abc", InterfaceField: TestStruct{StringField: "def"}},
			}},
	} {
		got := wrapper{}
		err := Unmarshal([]byte(test.input), &got)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("\ngot: \n%v\nwant: \n%v\n", got, test.want)
		}
	}
}

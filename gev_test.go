package gev

import (
	"reflect"
	"testing"
)

func Test_parse(t *testing.T) {
	type Example struct {
		In  string
		T   reflect.Type
		Out interface{}
		Err error
	}

	str, i64, f64, b := "hello", int64(202), float64(2.32), true
	inputs := []Example{
		Example{In: "hello", T: reflect.TypeOf(""), Out: str},
		Example{In: "hello", T: reflect.TypeOf(&str), Out: &str},
		Example{In: "hello", T: reflect.TypeOf([]byte{}), Out: []byte("hello")},

		Example{In: "202", T: reflect.TypeOf(int64(0)), Out: i64},
		Example{In: "202", T: reflect.TypeOf(&i64), Out: &i64},
		Example{In: "", T: reflect.TypeOf(&i64), Out: nil},

		Example{In: "2.320", T: reflect.TypeOf(float64(0.0)), Out: f64},
		Example{In: "2.320", T: reflect.TypeOf(&f64), Out: &f64},
		Example{In: "", T: reflect.TypeOf(&f64), Out: nil},

		Example{In: "true", T: reflect.TypeOf(false), Out: b},
		Example{In: "true", T: reflect.TypeOf(&b), Out: &b},
		Example{In: "", T: reflect.TypeOf(&b), Out: nil},
	}

	for _, in := range inputs {
		actual, err := parse(in.In, in.T)
		if !reflect.DeepEqual(err, in.Err) {
			t.Fatalf("expected: %v\n got: %v\n", in.Out, actual)
		}

		if !reflect.DeepEqual(in.Out, actual) {
			t.Fatalf("expected: %v (type: %T)\n got: %v (type %T)\n", in.Out, in.Out, actual, actual)
		}
	}
}

func Test_Unmarshal(t *testing.T) {
	type Example struct {
		A string `env:"foo"`
		B int64
		c string
		D string `env:"-"`
		e string `env:"-"`
		F bool
		G float64
	}

	tmpGet := get
	get = func(s string) string {
		return map[string]string{
			"foo": "hello",
			"B":   "342",
			"s":   "w",
			"F":   "true",
			"G":   "2.42",
			"D":   "word",
		}[s]
	}

	actual := Example{}
	err := Unmarshal(&actual)
	get = tmpGet
	if err != nil {
		t.Fatal(err)
	}

	expected := Example{A: "hello", B: int64(342), F: true, G: float64(2.42)}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected: %v\nGot: %v\n", expected, actual)
	}
}

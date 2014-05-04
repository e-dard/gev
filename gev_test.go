package gev

import (
	"reflect"
	"testing"
)

func Test_parse(t *testing.T) {
	type Example struct {
		In  []byte
		T   reflect.Type
		Out interface{}
		Err error
	}

	str, i64, f64, b := "foo", int64(202), float64(2.32), true
	inputs := []Example{
		// string
		Example{In: []byte("foo"), T: reflect.TypeOf(""), Out: str},
		Example{In: []byte("foo"), T: reflect.TypeOf(&str), Out: &str},
		Example{In: nil, T: reflect.TypeOf(&str), Out: nil},

		// []byte
		Example{In: []byte("foo"), T: reflect.TypeOf([]byte{}), Out: []byte("foo")},
		Example{In: []byte{}, T: reflect.TypeOf([]byte{}), Out: []byte{}},
		Example{In: nil, T: reflect.TypeOf([]byte{}), Out: nil},

		// int64
		Example{In: []byte("202"), T: reflect.TypeOf(int64(0)), Out: i64},
		Example{In: []byte("202"), T: reflect.TypeOf(&i64), Out: &i64},
		Example{In: nil, T: reflect.TypeOf(&i64), Out: nil},

		// float64
		Example{In: []byte("2.320"), T: reflect.TypeOf(float64(0.0)), Out: f64},
		Example{In: []byte("2.320"), T: reflect.TypeOf(&f64), Out: &f64},
		Example{In: nil, T: reflect.TypeOf(&f64), Out: nil},

		// bool
		Example{In: []byte("true"), T: reflect.TypeOf(false), Out: b},
		Example{In: []byte("true"), T: reflect.TypeOf(&b), Out: &b},
		Example{In: nil, T: reflect.TypeOf(&b), Out: nil},
	}

	for _, in := range inputs {
		actual, err := parse(in.In, in.T)
		if !reflect.DeepEqual(err, in.Err) {
			t.Fatalf("expected: %v\n got: %v\n", in.Out, actual)
		}

		if !reflect.DeepEqual(in.Out, actual) {
			msg := "expected: %[1]v (type: %[1]T)\n got: %[2]v (type %[2]T)\n"
			t.Fatalf(msg, in.Out, actual)
		}
	}
}

// wraps Unmarshal and uses mocked get function
func unmarshal(v interface{}) error {
	tmpGetVal, tmpVars := getVal, vars
	env := map[string]string{
		"foo": "hello",
		"B":   "342",
		"s":   "w",
		"F":   "true",
		"G":   "2.42",
		"D":   "word",
	}

	// mock funcs
	getVal = func(s string) string {
		return env[s]
	}
	vars = func() (out []string) {
		for k, _ := range env {
			out = append(out, k)
		}
		return
	}

	err := Unmarshal(v)

	// replace original functions
	getVal, vars = tmpGetVal, tmpVars
	return err
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

	actual := Example{}
	if err := unmarshal(&actual); err != nil {
		t.Fatal(err)
	}

	expected := Example{A: "hello", B: int64(342), F: true, G: float64(2.42)}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected: %v\nGot: %v\n", expected, actual)
	}
}

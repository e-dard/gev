package gev

import (
	"fmt"
	"os"
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
	iz, fz, bz := int64(0), float64(0), false
	inputs := []Example{
		// string
		Example{In: []byte("foo"), T: reflect.TypeOf(""), Out: str},
		Example{In: []byte("foo"), T: reflect.TypeOf(&str), Out: &str},
		Example{In: nil, T: reflect.TypeOf(&str), Out: nil},
		Example{In: nil, T: reflect.TypeOf(""), Out: ""},

		// []byte
		Example{In: []byte("foo"), T: reflect.TypeOf([]byte{}), Out: []byte("foo")},
		Example{In: []byte{}, T: reflect.TypeOf([]byte{}), Out: []byte{}},
		Example{In: nil, T: reflect.TypeOf([]byte{}), Out: nil},

		// int64
		Example{In: []byte("202"), T: reflect.TypeOf(int64(0)), Out: i64},
		Example{In: []byte("202"), T: reflect.TypeOf(&i64), Out: &i64},
		Example{In: nil, T: reflect.TypeOf(&i64), Out: nil},

		Example{In: nil, T: reflect.TypeOf(int64(1)), Out: int64(0),
			Err: fmt.Errorf(`cannot parse "" into type int64`)},
		Example{In: []byte("foo"), T: reflect.TypeOf(int64(1)), Out: int64(0),
			Err: fmt.Errorf(`cannot parse "foo" into type int64`)},
		Example{In: []byte("foo"), T: reflect.TypeOf(&i64), Out: &iz,
			Err: fmt.Errorf(`cannot parse "foo" into type *int64`)},

		// float64
		Example{In: []byte("2.320"), T: reflect.TypeOf(float64(0.0)), Out: f64},
		Example{In: []byte("2.320"), T: reflect.TypeOf(&f64), Out: &f64},
		Example{In: nil, T: reflect.TypeOf(&f64), Out: nil},

		Example{In: nil, T: reflect.TypeOf(float64(1)), Out: float64(0),
			Err: fmt.Errorf(`cannot parse "" into type float64`)},
		Example{In: []byte("foo"), T: reflect.TypeOf(float64(1)), Out: float64(0),
			Err: fmt.Errorf(`cannot parse "foo" into type float64`)},
		Example{In: []byte("foo"), T: reflect.TypeOf(&f64), Out: &fz,
			Err: fmt.Errorf(`cannot parse "foo" into type *float64`)},

		// bool
		Example{In: []byte("true"), T: reflect.TypeOf(false), Out: b},
		Example{In: []byte("true"), T: reflect.TypeOf(&b), Out: &b},
		Example{In: nil, T: reflect.TypeOf(&b), Out: nil},

		Example{In: nil, T: reflect.TypeOf(true), Out: false,
			Err: fmt.Errorf(`cannot parse "" into type bool`)},
		Example{In: []byte("foo"), T: reflect.TypeOf(false), Out: false,
			Err: fmt.Errorf(`cannot parse "foo" into type bool`)},
		Example{In: []byte("foo"), T: reflect.TypeOf(&b), Out: &bz,
			Err: fmt.Errorf(`cannot parse "foo" into type *bool`)},

		// unsupported types
		Example{In: []byte("foo"), T: reflect.TypeOf(int32(1)),
			Err: fmt.Errorf(`unsupported underlying type: int32`)},
		Example{In: []byte("foo"), T: reflect.TypeOf(&testing.T{}),
			Err: fmt.Errorf(`unsupported underlying type: T`)},
	}

	for _, in := range inputs {
		actual, err := parse(in.In, in.T)
		if !reflect.DeepEqual(err, in.Err) {
			t.Fatalf("expected: %v\n got: %v\n", in.Err, err)
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

func Test_getEnv(t *testing.T) {
	if err := os.Setenv("RANDOMENV", "22"); err != nil {
		panic(err)
	}

	if v := getEnv("RANDOMENV", os.Environ()); string(v) != "22" {
		t.Fatalf("Expected %q, got %q\n", "22", v)
	}

	if err := os.Setenv("NEWENV", "=22"); err != nil {
		panic(err)
	}

	if v := getEnv("NEWENV", os.Environ()); string(v) != "=22" {
		t.Fatalf("Expected %q, got %q\n", "=22", v)
	}

	if v := getEnv("ISMISSING", os.Environ()); string(v) != "" {
		t.Fatalf("Expected %q, got %q\n", "", v)
	}
}

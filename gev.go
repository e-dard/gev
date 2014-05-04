// Package gev implements functionality to unmarshal environment
// variables into struct fields.
//
package gev

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const tag = "env"

var (
	getVal = os.Getenv
	vars   = os.Environ
)

// Unmarshal inspects the process' environment for values that match
// `env` tags on v, and parses the values into the fields of v.
//
// Unmarshal can unmarshal environment variable values into the
// following types:
//
//	bool, int64, float64, string and []byte
//
// Further, it supports pointers to bool, int64, float64, string types.
// In the case of bool, int64 and float64, nil will be unmarshaled if
// the enviroment variable does not exist, or its value is the empty
// string.
//
// Unmarshal targets exported fields, and checks the gev tag value. It
// uses the following rules:
//
//	// Field will contain the value of FOO environment variable
//	Field string `env:"FOO"`
//
//	// field will be ignored by Unmarshal as it's unexported
//	field int
//
//	// Filed will be ignored by Unmarshal
//	Field bool `env:"-"`
//
func Unmarshal(v interface{}) error {
	st := reflect.TypeOf(v)
	if st.Kind() != reflect.Ptr || st.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("Unmarshal expects a pointer to a struct")
	}
	// gets the struct being pointed to
	st = st.Elem()

	// environment variables available
	names := vars()

	target := reflect.ValueOf(v).Elem()
	for i := 0; i < st.NumField(); i++ {
		fi := st.Field(i)
		// tag value
		tagv := fi.Tag.Get(tag)
		if tagv == "-" || fi.PkgPath != "" {
			// ignore fields tagged with "-" and unexported fields
			continue
		} else if tagv != "" {
			ev, err := parse(getEnv(tagv, names), fi.Type)
			if err != nil {
				return err
			}
			target.Field(i).Set(reflect.ValueOf(ev))
		} else {
			// field has no tag, but is exported field
			ev, err := parse(getEnv(fi.Name, names), fi.Type)
			if err != nil {
				return err
			}
			target.Field(i).Set(reflect.ValueOf(ev))
		}
	}
	return nil
}

// getEnv returns the value for an environment variable k. If there is
// no variable k then nil is returned, while the presence of k will
// always result in a non-nil slice being returned.
func getEnv(k string, env []string) []byte {
	for _, v := range env {
		if v == k {
			return []byte(getVal(k))
		}
	}
	return nil
}

func parse(v []byte, t reflect.Type) (out interface{}, err error) {
	if t.Kind() == reflect.Ptr {
		switch t.Elem().Kind() {
		case reflect.String:
			if v == nil {
				return
			}
			o := string(v)
			out = &o
		case reflect.Int64:
			if v == nil {
				return
			}
			o, err := strconv.ParseInt(string(v), 10, 64)
			if err != nil {
				err = fmt.Errorf("cannot parse %q into type *int64", v)
			}
			out = &o
		case reflect.Float64:
			if v == nil {
				return
			}
			o, err := strconv.ParseFloat(string(v), 64)
			if err != nil {
				err = fmt.Errorf("cannot parse %q into type *float64", v)
			}
			out = &o
		case reflect.Bool:
			if v == nil {
				return
			}
			o, err := strconv.ParseBool(string(v))
			if err != nil {
				err = fmt.Errorf("cannot parse %q into type *bool", v)
			}
			out = &o
		default:
			err = fmt.Errorf("unsupported underlying type: %T", t.Elem().Kind())
		}
		return
	}

	switch t.Kind() {
	case reflect.Slice:
		if t.Elem().Kind() != reflect.Uint8 {
			err = fmt.Errorf("cannot parse %q into type []%T", v, t.Elem().Kind())
			return
		}
		// nil will be returned if v is nil
		if v != nil {
			out = []byte(v)
		}
		return
	case reflect.String:
		out = string(v)
	case reflect.Int64:
		o, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			err = fmt.Errorf("cannot parse %q into type int64", v)
		}
		out = o
	case reflect.Float64:
		o, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			err = fmt.Errorf("cannot parse %q into type float64", v)
		}
		out = o
	case reflect.Bool:
		o, err := strconv.ParseBool(string(v))
		if err != nil {
			err = fmt.Errorf("cannot parse %q into type bool", v)
		}
		out = o
	default:
		err = fmt.Errorf("unsupported underlying type: %T", t.Kind())
	}
	return
}

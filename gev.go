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
	get = os.Getenv
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

	st = st.Elem()
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < st.NumField(); i++ {
		fi := st.Field(i)
		// tag value
		tagv := fi.Tag.Get(tag)
		if tagv == "-" || fi.PkgPath != "" {
			// ignore fields tagged with "-" and unexported fields
			continue
		} else if tagv != "" {
			v, err := parse(get(tagv), fi.Type)
			if err != nil {
				return err
			}
			val.Field(i).Set(reflect.ValueOf(v))
		} else {
			// field has no tag, but is exported field
			v, err := parse(get(fi.Name), fi.Type)
			if err != nil {
				return err
			}
			val.Field(i).Set(reflect.ValueOf(v))
		}
	}
	return nil
}

func parse(v string, t reflect.Type) (out interface{}, err error) {
	if t.Kind() == reflect.Slice {
		if t.Elem().Kind() != reflect.Uint8 {
			panic("Unsupported type")
		}
		// return slice of bytes
		out = []byte(v)
		return
	}

	if t.Kind() == reflect.Ptr {
		switch t.Elem().Kind() {
		case reflect.String:
			out = &v
		case reflect.Int64:
			if v == "" {
				break
			}
			o, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = fmt.Errorf("cannot parse %q into type *int64", v)
			}
			out = &o
		case reflect.Float64:
			if v == "" {
				break
			}
			o, err := strconv.ParseFloat(v, 64)
			if err != nil {
				err = fmt.Errorf("cannot parse %q into type *float64", v)
			}
			out = &o
		case reflect.Bool:
			if v == "" {
				break
			}
			o, err := strconv.ParseBool(v)
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
	case reflect.String:
		out = v
	case reflect.Int64:
		o, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			err = fmt.Errorf("cannot parse %q into type int64", v)
		}
		out = o
	case reflect.Float64:
		o, err := strconv.ParseFloat(v, 64)
		if err != nil {
			err = fmt.Errorf("cannot parse %q into type float64", v)
		}
		out = o
	case reflect.Bool:
		o, err := strconv.ParseBool(v)
		if err != nil {
			err = fmt.Errorf("cannot parse %q into type bool", v)
		}
		out = o
	default:
		err = fmt.Errorf("unsupported underlying type: %T", t.Kind())
	}
	return
}

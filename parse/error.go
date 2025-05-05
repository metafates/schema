package parse

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// InvalidParseError describes an invalid argument passed to [Parse].
// (The argument to [Parse] must be a non-nil pointer.)
type InvalidParseError struct {
	Type reflect.Type
}

func (e InvalidParseError) Error() string {
	if e.Type == nil {
		return "Parse(nil)"
	}

	if e.Type.Kind() != reflect.Pointer {
		return "Parse(non-pointer " + e.Type.String() + ")"
	}

	return "Validate(nil " + e.Type.String() + ")"
}

type UnconvertableTypeError struct {
	Target, Original string
}

func (e UnconvertableTypeError) Error() string {
	return fmt.Sprintf("can not convert %s to %s", e.Original, e.Target)
}

type UnknownFieldError struct {
	Name string
}

func (e UnknownFieldError) Error() string {
	return fmt.Sprintf("unknown field: %s", e.Name)
}

type ParseError struct {
	Msg   string
	Inner error

	path string
}

// Path returns the path to the value which raised this error
func (e ParseError) Path() string {
	var recursive func(path []string, err error) []string

	recursive = func(path []string, err error) []string {
		var validationErr ParseError

		if !errors.As(err, &validationErr) {
			return path
		}

		if validationErr.path != "" {
			path = append(path, validationErr.path)
		}

		return recursive(path, validationErr.Inner)
	}

	return strings.Join(recursive(nil, e), "")
}

func (e ParseError) Error() string {
	return "parse: " + e.error()
}

func (e ParseError) Unwrap() error {
	return e.Inner
}

func (e ParseError) error() string {
	segments := make([]string, 0, 3)

	if e.path != "" {
		segments = append(segments, e.path)
	}

	if e.Msg != "" {
		segments = append(segments, e.Msg)
	}

	if e.Inner != nil {
		var pe ParseError

		if errors.As(e.Inner, &pe) {
			segments = append(segments, pe.error())
		} else {
			segments = append(segments, e.Inner.Error())
		}
	}

	// path: msg: inner error
	return strings.Join(segments, ": ")
}

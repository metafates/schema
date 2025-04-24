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
		validationErr, ok := err.(ParseError)
		if !ok {
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
	segments := make([]string, 0, 3)

	if e.path != "" {
		segments = append(segments, e.path)
	}

	if e.Msg != "" {
		segments = append(segments, e.Msg)
	}

	if e.Inner != nil {
		segments = append(segments, e.Inner.Error())
	}

	// path: msg: inner error
	return strings.Join(segments, ": ")
}

func (e ParseError) Is(err error) bool {
	other, ok := err.(ParseError)
	if !ok {
		return errors.Is(e.Inner, err)
	}

	switch {
	case e.Inner == nil && other.Inner == nil:
		return e.Msg == other.Msg

	case e.Inner == nil && other.Inner != nil:
		return errors.Is(e, other.Inner)

	case e.Inner != nil && other.Inner == nil:
		return errors.Is(e.Inner, other)

	default:
		return errors.Is(e.Inner, other.Inner)
	}
}

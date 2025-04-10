package validate

import (
	"errors"
	"reflect"
	"strings"
)

// InvalidValidateError describes an invalid argument passed to [Validate].
// (The argument to [Validate] must be a non-nil pointer.)
type InvalidValidateError struct {
	Type reflect.Type
}

func (e InvalidValidateError) Error() string {
	if e.Type == nil {
		return "Validate(nil)"
	}

	if e.Type.Kind() != reflect.Pointer {
		return "Validate(non-pointer " + e.Type.String() + ")"
	}

	return "Validate(nil " + e.Type.String() + ")"
}

// ValidationError describes validation error occured at [validate.Validate]
type ValidationError struct {
	Msg   string
	Inner error

	path string
}

func (v ValidationError) WithPath(path string) ValidationError {
	v.path = path

	return v
}

func (v ValidationError) Path() string {
	var recursive func(path []string, err error) []string

	recursive = func(path []string, err error) []string {
		validationErr, ok := err.(ValidationError)
		if !ok {
			return path
		}

		if validationErr.path != "" {
			path = append(path, validationErr.path)
		}

		return recursive(path, validationErr.Inner)
	}

	return strings.Join(recursive(nil, v), "")
}

func (v ValidationError) Error() string {
	segments := make([]string, 0, 3)

	if v.path != "" {
		segments = append(segments, v.path)
	}

	if v.Msg != "" {
		segments = append(segments, v.Msg)
	}

	if v.Inner != nil {
		segments = append(segments, v.Inner.Error())
	}

	// path: msg: inner error
	return strings.Join(segments, ": ")
}

func (v ValidationError) Is(err error) bool {
	other, ok := err.(ValidationError)
	if !ok {
		return errors.Is(v.Inner, err)
	}

	switch {
	case v.Inner == nil && other.Inner == nil:
		return v.Msg == other.Msg

	case v.Inner == nil && other.Inner != nil:
		return errors.Is(v, other.Inner)

	case v.Inner != nil && other.Inner == nil:
		return errors.Is(v.Inner, other)

	default:
		return errors.Is(v.Inner, other.Inner)
	}
}

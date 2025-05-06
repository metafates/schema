package validate

import (
	"errors"
	"reflect"
	"strings"
)

// InvalidValidateError describes an invalid argument passed to [Validate].
// The argument to [Validate] must be a non-nil pointer.
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

// ValidationError describes validation error occurred at [Validate].
type ValidationError struct {
	Msg   string
	Inner error

	path string
}

// WithPath returns a copy of [ValidationError] with the given path set.
func (e ValidationError) WithPath(path string) ValidationError {
	e.path = path

	return e
}

// Path returns the path to the value which raised this error.
func (e ValidationError) Path() string {
	var recursive func(path []string, err error) []string

	recursive = func(path []string, err error) []string {
		var validationErr ValidationError

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

func (e ValidationError) Error() string {
	return "validate: " + e.error()
}

func (e ValidationError) Unwrap() error {
	return e.Inner
}

func (e ValidationError) error() string {
	segments := make([]string, 0, 3)

	if e.path != "" {
		segments = append(segments, e.path)
	}

	if e.Msg != "" {
		segments = append(segments, e.Msg)
	}

	if e.Inner != nil {
		var ve ValidationError

		if errors.As(e.Inner, &ve) {
			segments = append(segments, ve.error())
		} else {
			segments = append(segments, e.Inner.Error())
		}
	}

	// path: msg: inner error
	return strings.Join(segments, ": ")
}

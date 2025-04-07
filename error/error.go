package schemaerror

import (
	"errors"
	"fmt"
)

var ErrMissingRequiredValue = ValidationError{Msg: "missing required value"}

type ValidationError struct {
	Msg   string
	Inner error
}

func (v ValidationError) Error() string {
	if v.Inner == nil {
		return v.Msg
	}

	if v.Msg == "" {
		return v.Inner.Error()
	}

	return fmt.Sprintf("%s: %s", v.Msg, v.Inner)
}

func (v ValidationError) Is(err error) bool {
	other, ok := err.(ValidationError)
	if !ok {
		return errors.Is(v.Inner, err)
	}

	if v.Msg != other.Msg {
		return false
	}

	return errors.Is(v.Inner, other.Inner)
}

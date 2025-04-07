package schemaerror

import (
	"fmt"
)

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

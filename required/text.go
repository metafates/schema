package required

import (
	"encoding"

	"github.com/metafates/schema/validate"
)

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*Custom[any, validate.Validator[any]])(nil)

// UnmarshalText implements the [encoding.TextUnmarshaler] interface
func (c *Custom[T, V]) UnmarshalText(data []byte) error {
	return c.UnmarshalJSON(data)
}

// MarshalText implements the [encoding.TextMarshaler] interface
func (c Custom[T, V]) MarshalText() ([]byte, error) {
	return c.MarshalJSON()
}

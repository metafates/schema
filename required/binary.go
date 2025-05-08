package required

import (
	"encoding"

	"github.com/metafates/schema/validate"
)

var _ interface {
	encoding.BinaryUnmarshaler
	encoding.BinaryMarshaler
} = (*Custom[any, validate.Validator[any]])(nil)

// UnmarshalBinary implements the [encoding.BinaryUnmarshaler] interface.
func (c *Custom[T, V]) UnmarshalBinary(data []byte) error {
	return c.GobDecode(data)
}

// MarshalBinary implements the [encoding.BinaryMarshaler] interface.
func (c Custom[T, V]) MarshalBinary() ([]byte, error) {
	return c.GobEncode()
}

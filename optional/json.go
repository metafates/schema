package optional

import (
	"encoding/json"

	"github.com/metafates/schema/validate"
)

var _ interface {
	json.Unmarshaler
	json.Marshaler
} = (*Custom[any, validate.Validator[any]])(nil)

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (c *Custom[T, V]) UnmarshalJSON(data []byte) error {
	var value *T

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	// validated status will reset here
	if value == nil {
		*c = Custom[T, V]{}

		return nil
	}

	*c = Custom[T, V]{value: *value, hasValue: true}

	return nil
}

// MarshalJSON implements the [json.Marshaler] interface.
func (c Custom[T, V]) MarshalJSON() ([]byte, error) {
	if c.hasValue && !c.validated {
		panic("called MarshalJSON() on unvalidated value")
	}

	if c.hasValue {
		return json.Marshal(c.value)
	}

	return []byte("null"), nil
}

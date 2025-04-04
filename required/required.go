package required

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/metafates/schema/constraint"
	"github.com/metafates/schema/validate"
)

var _ interface {
	json.Unmarshaler
	json.Marshaler
} = (*Custom[any, validate.Any[any]])(nil)

type (
	Custom[T any, V validate.Validator[T]] struct {
		value T
	}

	T[T any] struct {
		Custom[T, validate.Any[T]]
	}

	NonEmpty[T comparable] struct {
		Custom[T, validate.NonEmpty[T]]
	}

	Positive[T constraint.Real] struct {
		Custom[T, validate.Positive[T]]
	}

	Negative[T constraint.Real] struct {
		Custom[T, validate.Negative[T]]
	}

	Alphanumeric[T ~string] struct {
		Custom[T, validate.Alphanumeric[T]]
	}
)

func (c Custom[T, V]) Value() T { return c.value }

func (c *Custom[T, V]) UnmarshalJSON(data []byte) error {
	var value *T

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	if value == nil {
		return errors.New("required value is missing")
	}

	if err := (*new(V)).Validate(*value); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	c.value = *value

	return nil
}

func (c *Custom[T, V]) UnmarshalText(data []byte) error {
	return c.UnmarshalJSON(data)
}

func (c Custom[T, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

func (c Custom[T, V]) MarshalText() ([]byte, error) {
	return c.MarshalJSON()
}

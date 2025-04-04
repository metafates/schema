package optional

import (
	"encoding/json"
	"fmt"

	"github.com/metafates/schema/constraint"
	"github.com/metafates/schema/validate"
)

type (
	Custom[T any, V validate.Validator[T]] struct {
		value    T
		hasValue bool
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

func (c Custom[T, V]) HasValue() bool   { return c.hasValue }
func (c Custom[T, V]) Value() (T, bool) { return c.value, c.hasValue }
func (c Custom[T, V]) Must() T {
	if c.hasValue {
		return c.value
	}

	panic("called must on empty optional")
}

func (c *Custom[T, V]) UnmarshalJSON(data []byte) error {
	var value *T

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	if value == nil {
		*c = Custom[T, V]{}
		return nil
	}

	if err := (*new(V)).Validate(*value); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	*c = Custom[T, V]{value: *value, hasValue: true}

	return nil
}

func (c *Custom[T, V]) UnmarshalText(data []byte) error {
	return c.UnmarshalJSON(data)
}

func (c Custom[T, V]) MarshalJSON() ([]byte, error) {
	if c.hasValue {
		return json.Marshal(c.value)
	}

	return []byte("null"), nil
}

func (c Custom[T, V]) MarshalText() ([]byte, error) {
	return c.MarshalJSON()
}

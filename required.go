package schema

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/metafates/schema/constraint"
	"github.com/metafates/schema/validate"
)

type (
	RequiredAnd[T any, V validate.Validator[T]] struct {
		value T
	}

	Required[T any] struct {
		RequiredAnd[T, validate.Nop[T]]
	}

	RequiredNonEmpty[T comparable] struct {
		RequiredAnd[T, validate.NonEmpty[T]]
	}

	RequiredPositive[T constraint.Real] struct {
		RequiredAnd[T, validate.Positive[T]]
	}
)

func (r RequiredAnd[T, V]) Value() T { return r.value }

func (r *RequiredAnd[T, V]) UnmarshalJSON(data []byte) error {
	var value *T

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	if value == nil {
		return errors.New("required value is missing")
	}

	var validator V

	if err := validator.Validate(*value); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	r.value = *value

	return nil
}

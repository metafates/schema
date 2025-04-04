package validate

import (
	"errors"

	"github.com/metafates/schema/constraint"
)

type Validator[T any] interface {
	Validate(value T) error
}

var (
	_ Validator[any] = (*Nop[any])(nil)
	_ Validator[any] = (*NonEmpty[any])(nil)
)

type (
	Nop[T any]                  struct{}
	NonEmpty[T comparable]      struct{}
	Positive[T constraint.Real] struct{}
)

func (Nop[T]) Validate(T) error {
	return nil
}

func (NonEmpty[T]) Validate(value T) error {
	var empty T

	if value == empty {
		return errors.New("empty value")
	}

	return nil
}

func (Positive[T]) Validate(value T) error {
	if value < 0 {
		return errors.New("negative value")
	}

	return nil
}

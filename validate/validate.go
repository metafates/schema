package validate

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/metafates/required/constraint"
)

type ValidateError string

func (v ValidateError) Error() string {
	return string(v)
}

type Validator[T any] interface {
	Validate(value T) error
}

var (
	_ Validator[any] = (*Any[any])(nil)
	_ Validator[any] = (*NonEmpty[any])(nil)
	_ Validator[int] = (*Positive[int])(nil)
	_ Validator[int] = (*Negative[int])(nil)
	_ Validator[any] = (*Combined[Validator[any], Validator[any], any])(nil)
)

type (
	Any[T any]                  struct{}
	NonEmpty[T comparable]      struct{}
	Positive[T constraint.Real] struct{}
	Negative[T constraint.Real] struct{}
	Alphanumeric[T ~string]     struct{}

	// Combined is a meta validator that combines other validators
	Combined[A Validator[T], B Validator[T], T any] struct{}
)

func (Any[T]) Validate(T) error {
	return nil
}

func (NonEmpty[T]) Validate(value T) error {
	var empty T

	if value == empty {
		return ValidateError("empty value")
	}

	return nil
}

func (Positive[T]) Validate(value T) error {
	if value < 0 {
		return ValidateError("negative value")
	}

	return nil
}

func (Negative[T]) Validate(value T) error {
	if value > 0 {
		return ValidateError("positive value")
	}

	return nil
}

func (Alphanumeric[T]) Validate(value T) error {
	idx := strings.IndexFunc(string(value), func(r rune) bool {
		return !unicode.IsLetter(r) || !unicode.IsNumber(r)
	})

	if value != "" && idx >= 0 {
		r := []rune(value)[idx]

		return ValidateError(fmt.Sprintf("unexpected character %q", string(r)))
	}

	return nil
}

func (Combined[A, B, T]) Validate(value T) error {
	if err := (*new(A)).Validate(value); err != nil {
		return err
	}

	if err := (*new(B)).Validate(value); err != nil {
		return err
	}

	return nil
}

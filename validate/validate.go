package validate

import (
	"errors"
	"net/mail"

	"github.com/metafates/schema/constraint"
)

type ValidateError struct{ inner error }

func (v ValidateError) Error() string {
	return v.inner.Error()
}

type Validator[T any] interface {
	Validate(value T) error
}

var (
	_ Validator[any]    = (*Any[any])(nil)
	_ Validator[any]    = (*NonEmpty[any])(nil)
	_ Validator[int]    = (*Positive[int])(nil)
	_ Validator[int]    = (*Negative[int])(nil)
	_ Validator[string] = (*Email[string])(nil)

	_ Validator[any] = (*Combined[Validator[any], Validator[any], any])(nil)
)

type (
	// Any accepts any value
	Any[T any] struct{}

	// NonEmpty accepts all non empty comparable values
	NonEmpty[T comparable] struct{}

	// Positive accepts all positive real numbers and zero
	//
	// See also [Negative]
	Positive[T constraint.Real] struct{}

	// Negative accepts all negative real numbers and zero
	//
	// See also [Positive]
	Negative[T constraint.Real] struct{}

	// Email accepts a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>"
	Email[T ~string] struct{}

	// Combined is a meta validator that combines other validators
	Combined[A Validator[T], B Validator[T], T any] struct{}
)

func (Any[T]) Validate(T) error {
	return nil
}

func (NonEmpty[T]) Validate(value T) error {
	var empty T

	if value == empty {
		return ValidateError{inner: errors.New("empty value")}
	}

	return nil
}

func (Positive[T]) Validate(value T) error {
	if value < 0 {
		return ValidateError{inner: errors.New("negative value")}
	}

	return nil
}

func (Negative[T]) Validate(value T) error {
	if value > 0 {
		return ValidateError{inner: errors.New("positive value")}
	}

	return nil
}

func (Email[T]) Validate(value T) error {
	_, err := mail.ParseAddress(string(value))
	if err != nil {
		return ValidateError{inner: err}
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

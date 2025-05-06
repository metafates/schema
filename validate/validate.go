// Package validate provides type enforced validators.
package validate

import (
	"reflect"

	"github.com/metafates/schema/internal/reflectwalk"
)

type (
	// Validator is an interface that validators must implement.
	// It's a special empty (struct{}) type that is invoked in a form of (*new(V)).Validate(...).
	// Therefore it should not depend on inner state (fields).
	Validator[T any] interface {
		Validate(value T) error
	}

	// TypeValidateable is an interface for types that can validate their types.
	// This is used by required and optional fields so that they can validate if contained values
	// satisfy the schema enforced by [Validator] backed type.
	//
	// TL;DR: do not implement nor use this method directly (codegen is exception).
	//
	// See [Validateable] interface if you want to implement custom validation.
	TypeValidateable interface {
		TypeValidate() error
	}

	// Validateable is an interface for types that can perform validation logic after
	// type validation (by [TypeValidateable]) has been called without errors.
	//
	// Primary usecase is custom cross-field validation. E.g. if X is true then Y cannot be empty.
	Validateable interface {
		Validate() error
	}
)

// Validate checks if the provided value can be validated and reports any validation errors.
//
// The validation process follows these steps:
//  1. If v implements both [TypeValidateable] and [Validateable] interfaces, its [TypeValidateable.TypeValidate]
//     method is called, followed by its [Validateable.Validate] method.
//  2. If v only implements the [TypeValidateable] interface, its [TypeValidateable.TypeValidate] method is called.
//  3. Otherwise, Validate traverses all fields in the struct pointed to by v, applying the same validation
//     logic to each field.
//
// A [ValidationError] is returned if any validation fails during any step.
//
// If v is nil or not a pointer, Validate returns an [InvalidValidateError].
func Validate(v any) error {
	switch v := v.(type) {
	case interface {
		TypeValidateable
		Validateable
	}:
		if err := v.TypeValidate(); err != nil {
			return ValidationError{Inner: err}
		}

		if err := v.Validate(); err != nil {
			return ValidationError{Inner: err}
		}

		return nil

	case TypeValidateable:
		if err := v.TypeValidate(); err != nil {
			return ValidationError{Inner: err}
		}

		return nil
	}

	// same thing [json.Unmarshal] does
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidValidateError{Type: reflect.TypeOf(v)}
	}

	var postValidate []func() error

	err := reflectwalk.WalkFields(v, func(path string, reflectValue reflect.Value) error {
		if reflectValue.CanAddr() {
			reflectValue = reflectValue.Addr()
		}

		value := reflectValue.Interface()

		if value, ok := value.(TypeValidateable); ok {
			if err := value.TypeValidate(); err != nil {
				return ValidationError{Inner: err, path: path}
			}
		}

		if value, ok := value.(Validateable); ok {
			postValidate = append(postValidate, func() error {
				if err := value.Validate(); err != nil {
					return ValidationError{Inner: err, path: path}
				}

				return nil
			})
		}

		return nil
	})
	if err != nil {
		return err
	}

	for _, f := range postValidate {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}

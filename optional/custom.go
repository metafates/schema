// Package optional provides types whose values may be either empty (null) or be present and pass validation.
//
// Optional types support the following encoding/decoding formats:
//   - json
//   - sql
//   - text
//   - binary
//   - gob
package optional

import (
	"reflect"

	"github.com/metafates/schema/parse"
	"github.com/metafates/schema/validate"
)

// Custom optional type.
// When given non-null value it errors if validation fails.
type Custom[T any, V validate.Validator[T]] struct {
	value     T
	hasValue  bool
	validated bool
}

// TypeValidate implements the [validate.TypeValidateable] interface.
// You should not call this function directly.
func (c *Custom[T, V]) TypeValidate() error {
	if !c.hasValue {
		return nil
	}

	if err := (*new(V)).Validate(c.value); err != nil {
		return validate.ValidationError{Inner: err}
	}

	// validate nested types recursively
	if err := validate.Validate(&c.value); err != nil {
		return err
	}

	c.validated = true

	return nil
}

// HasValue returns the presence of the contained value.
func (c Custom[T, V]) HasValue() bool { return c.hasValue }

// Get returns the contained value and a boolean stating its presence.
// True if value exists, false otherwise.
//
// Panics if value was not validated yet.
// See also [Custom.GetPtr].
func (c Custom[T, V]) Get() (T, bool) {
	if c.hasValue && !c.validated {
		panic("called Get() on non-empty unvalidated value")
	}

	return c.value, c.hasValue
}

// Get returns the pointer to the contained value.
// Non-nil if value exists, nil otherwise.
// Pointed value is a shallow copy.
//
// Panics if value was not validated yet.
// See also [Custom.Get].
func (c Custom[T, V]) GetPtr() *T {
	if c.hasValue && !c.validated {
		panic("called GetPtr() on non-empty unvalidated value")
	}

	var value *T

	if c.hasValue {
		valueCopy := c.value
		value = &valueCopy
	}

	return value
}

// Must returns the contained value and panics if it does not have one.
// You can check for its presence using [Custom.HasValue] or use a more safe alternative [Custom.Get].
func (c Custom[T, V]) Must() T {
	if !c.hasValue {
		panic("called must on empty optional")
	}

	value, _ := c.Get()

	return value
}

// Parse checks if given value is valid.
// If it is, a value is used to initialize this type.
// Value is converted to the target type T, if possible. If not - [parse.UnconvertableTypeError] is returned.
// It is allowed to pass convertable type wrapped in optional type.
//
// Parsed type is validated, therefore it is safe to call [Custom.Get] afterwards.
//
// Passing nil results a valid empty instance.
func (c *Custom[T, V]) Parse(value any) error {
	if value == nil {
		*c = Custom[T, V]{}

		return nil
	}

	rValue := reflect.ValueOf(value)

	if rValue.Kind() == reflect.Pointer && rValue.IsNil() {
		*c = Custom[T, V]{}

		return nil
	}

	if _, ok := value.(interface{ isOptional() }); ok {
		// NOTE: ensure this method name is in sync with [Custom.Get]
		res := rValue.MethodByName("Get").Call(nil)
		v, ok := res[0], res[1].Bool()

		if !ok {
			*c = Custom[T, V]{}

			return nil
		}

		rValue = v
	}

	v, err := convert[T](rValue)
	if err != nil {
		return parse.ParseError{Inner: err}
	}

	aux := Custom[T, V]{
		hasValue: true,
		value:    v,
	}

	if err := aux.TypeValidate(); err != nil {
		return err
	}

	*c = aux

	return nil
}

func (c *Custom[T, V]) MustParse(value any) {
	if err := c.Parse(value); err != nil {
		panic("MustParse failed")
	}
}

func (Custom[T, V]) isOptional() {}

func convert[T any](v reflect.Value) (T, error) {
	tType := reflect.TypeFor[T]()

	original := v

	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.CanConvert(tType) {
		//nolint:forcetypeassert // checked already by CanConvert
		return v.Convert(tType).Interface().(T), nil
	}

	return *new(T), parse.UnconvertableTypeError{
		Target:   tType.String(),
		Original: original.Type().String(),
	}
}

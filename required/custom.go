// Package required provides types whose values must be present and pass validation.
//
// Required types support the following encoding/decoding formats:
//   - json
//   - sql
//   - text
//   - binary
//   - gob
package required

import (
	"reflect"

	"github.com/metafates/schema/parse"
	"github.com/metafates/schema/validate"
)

var (
	ErrMissingValue  = validate.ValidationError{Msg: "missing required value"}
	ErrParseNilValue = parse.ParseError{Msg: "nil value passed for parsing"}
)

// Custom required type.
// Errors if value is missing or did not pass the validation.
type Custom[T any, V validate.Validator[T]] struct {
	value     T
	hasValue  bool
	validated bool
}

// TypeValidate implements the [validate.TypeValidateable] interface.
// You should not call this function directly.
func (c *Custom[T, V]) TypeValidate() error {
	if !c.hasValue {
		return ErrMissingValue
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

// Get returns the contained value.
// Panics if value was not validated yet.
func (c Custom[T, V]) Get() T {
	if !c.validated {
		panic("called Get() on unvalidated value")
	}

	return c.value
}

// Parse checks if given value is valid.
// If it is, a value is used to initialize this type.
// Value is converted to the target type T, if possible. If not - [parse.UnconvertableTypeError] is returned.
// It is allowed to pass convertable type wrapped in required type.
//
// Parsed type is validated, therefore it is safe to call [Custom.Get] afterwards.
func (c *Custom[T, V]) Parse(value any) error {
	if value == nil {
		return ErrParseNilValue
	}

	rValue := reflect.ValueOf(value)

	if rValue.Kind() == reflect.Pointer {
		if rValue.IsNil() {
			return ErrParseNilValue
		}

		rValue = rValue.Elem()
	}

	tType := reflect.TypeFor[T]()

	if _, ok := value.(interface{ isRequired() }); ok {
		// NOTE: ensure this method name is in sync with [Custom.Get]
		rValue = rValue.MethodByName("Get").Call(nil)[0]
	}

	if !rValue.CanConvert(tType) {
		return parse.ParseError{
			Inner: parse.UnconvertableTypeError{
				Target:   tType.String(),
				Original: rValue.Type().String(),
			},
		}
	}

	//nolint:forcetypeassert // checked already by CanConvert
	aux := Custom[T, V]{
		value:    rValue.Convert(tType).Interface().(T),
		hasValue: true,
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

func (Custom[T, V]) isRequired() {}

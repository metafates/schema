package required

import (
	"encoding/json"

	"github.com/metafates/schema/constraint"
	schemaerror "github.com/metafates/schema/error"
	"github.com/metafates/schema/validate"
)

var _ interface {
	json.Unmarshaler
	json.Marshaler
} = (*Custom[any, validate.Any[any]])(nil)

type (
	// Custom required type.
	// Erorrs if value is missing or did not pass the validation
	Custom[T any, V validate.Validator[T]] struct {
		value    T
		hasValue bool
	}

	// Any accepts any value
	Any[T any] struct {
		Custom[T, validate.Any[T]]
	}

	// NonEmpty accepts all non empty comparable values
	NonEmpty[T comparable] struct {
		Custom[T, validate.NonEmpty[T]]
	}

	// Positive accepts all positive real numbers and zero
	//
	// See also [Negative]
	Positive[T constraint.Real] struct {
		Custom[T, validate.Positive[T]]
	}

	// Negative accepts all negative real numbers and zero
	//
	// See also [Positive]
	Negative[T constraint.Real] struct {
		Custom[T, validate.Negative[T]]
	}

	// Email accepts a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>"
	Email[T constraint.Text] struct {
		Custom[T, validate.Email[T]]
	}

	// URL accepts a single url.
	// The url may be relative (a path, without a host) or absolute (starting with a scheme)
	URL[T constraint.Text] struct {
		Custom[T, validate.URL[T]]
	}

	// IP accepts an IP address.
	// The address can be in dotted decimal ("192.0.2.1"), IPv6 ("2001:db8::68"), or IPv6 with a scoped addressing zone ("fe80::1cc0:3e8c:119f:c2e1%ens18").
	IP[T constraint.Text] struct {
		Custom[T, validate.IP[T]]
	}

	// Printable accepts strings consisting of only printable runes.
	// See [unicode.IsPrint] for more information
	Printable[T constraint.Text] struct {
		Custom[T, validate.Printable[T]]
	}

	// Base64 accepts valid base64 encoded strings
	Base64[T constraint.Text] struct {
		Custom[T, validate.Base64[T]]
	}

	// ASCII accepts ascii-only strings
	ASCII[T constraint.Text] struct {
		Custom[T, validate.ASCII[T]]
	}

	// PrintableASCII combines [Printable] and [ASCII]
	PrintableASCII[T constraint.Text] struct {
		Custom[T, validate.PrintableASCII[T]]
	}

	// Latitude accepts any number in the range [-90; 90]
	//
	// See also [Longitude]
	Latitude[T constraint.Real] struct {
		Custom[T, validate.Latitude[T]]
	}

	// Longitude accepts any number in the range [-180; 180]
	//
	// See also [Latitude]
	Longitude[T constraint.Real] struct {
		Custom[T, validate.Longitude[T]]
	}

	// InPast accepts any time before current timestamp
	//
	// See also [InFuture]
	InPast[T constraint.Time] struct {
		Custom[T, validate.InPast[T]]
	}

	// InFuture accepts any time after current timestamp
	//
	// See also [InPast]
	InFuture[T constraint.Time] struct {
		Custom[T, validate.InFuture[T]]
	}
)

// Validate implementes [validate.Validateable].
// You should not call this function directly.
func (c Custom[T, V]) Validate() error {
	if !c.hasValue {
		return schemaerror.ErrMissingRequiredValue
	}

	if err := (*new(V)).Validate(c.value); err != nil {
		return schemaerror.ValidationError{Inner: err}
	}

	if err := validate.Recursively(c.value); err != nil {
		return err
	}

	return nil
}

// Value returns the contained value.
// Panics if value was not initiliased (e.g. not unmarshalled from json)
func (c Custom[T, V]) Value() T {
	if !c.hasValue {
		panic("called Value() on uninitialized required value")
	}

	return c.value
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

	*c = Custom[T, V]{value: *value, hasValue: true}

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

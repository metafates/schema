package required

import (
	"encoding/json"

	"github.com/metafates/schema/constraint"
	"github.com/metafates/schema/validate"
)

var _ interface {
	json.Unmarshaler
	json.Marshaler
} = (*Custom[any, validate.Any[any]])(nil)

var ErrMissingValue = validate.ValidationError{Msg: "missing required value"}

type (
	// Custom required type.
	// Errors if value is missing or did not pass the validation
	Custom[T any, V validate.Validator[T]] struct {
		value     T
		hasValue  bool
		validated bool
	}

	// Any accepts any value of T
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

	// Even accepts real numbers divisible by two
	Even[T constraint.Integer] struct {
		Custom[T, validate.Even[T]]
	}

	// Odd accepts real numbers not divisible by two
	Odd[T constraint.Integer] struct {
		Custom[T, validate.Odd[T]]
	}

	// Email accepts a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>"
	Email[T constraint.Text] struct {
		Custom[T, validate.Email[T]]
	}

	// URL accepts a single url.
	// The url may be relative (a path, without a host) or absolute (starting with a scheme)
	//
	// See also [HTTPURL]
	URL[T constraint.Text] struct {
		Custom[T, validate.URL[T]]
	}

	// HTTPURL accepts a single http(s) url.
	//
	// See also [URL]
	HTTPURL[T constraint.Text] struct {
		Custom[T, validate.HTTPURL[T]]
	}

	// IP accepts an IP address.
	// The address can be in dotted decimal ("192.0.2.1"), IPv6 ("2001:db8::68"), or IPv6 with a scoped addressing zone ("fe80::1cc0:3e8c:119f:c2e1%ens18").
	IP[T constraint.Text] struct {
		Custom[T, validate.IP[T]]
	}

	// MAC accepts an IEEE 802 MAC-48, EUI-48, EUI-64, or a 20-octet IP over InfiniBand link-layer address
	MAC[T constraint.Text] struct {
		Custom[T, validate.MAC[T]]
	}

	// CIDR accepts CIDR notation IP address and prefix length, like "192.0.2.0/24" or "2001:db8::/32", as defined in RFC 4632 and RFC 4291
	CIDR[T constraint.Text] struct {
		Custom[T, validate.CIDR[T]]
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

	// Unique accepts an array of unique comparable values
	Unique[S ~[]T, T comparable] struct {
		Custom[S, validate.Unique[S, T]]
	}

	// MIME accepts RFC 1521 mime type string
	MIME[T constraint.Text] struct {
		Custom[T, validate.MIME[T]]
	}
)

// Validate implementes [validate.Validateable].
// You should not call this function directly.
func (c *Custom[T, V]) Validate() error {
	if !c.hasValue {
		return ErrMissingValue
	}

	if err := (*new(V)).Validate(c.value); err != nil {
		return validate.ValidationError{Inner: err}
	}

	if err := validate.Validate(&c.value); err != nil {
		return err
	}

	c.validated = true

	return nil
}

// Value returns the contained value.
// Panics if value was not validated yet
func (c Custom[T, V]) Value() T {
	if !c.validated {
		panic("called Value() on unvalidated value")
	}

	return c.value
}

func (c *Custom[T, V]) UnmarshalJSON(data []byte) error {
	var value *T

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	// validated status will reset here
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

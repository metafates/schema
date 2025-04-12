// Optional types may be either empty or contain a value.
//
// They support the following encoding/decoding:
// - json
// - sql
// - text
package optional

import (
	"github.com/metafates/schema/constraint"
	"github.com/metafates/schema/validate"
)

type (
	// Custom optional type.
	// When given non-null value it errors if validation fails
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

	// UUID accepts a properly formatted UUID in one of the following formats:
	//   xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	//   urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	//   xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	//   {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
	UUID[T constraint.Text] struct {
		Custom[T, validate.UUID[T]]
	}

	// NonEmptyPrintable combines [NonEmpty] and [Printable]
	NonEmptyPrintable[T ~string] struct {
		Custom[T, validate.And[
			T,
			validate.NonEmpty[T],
			validate.Printable[T],
		]]
	}
)

// Validate implements the [validate.Validateable] interface.
// You should not call this function directly.
func (c *Custom[T, V]) Validate() error {
	if !c.hasValue {
		c.validated = true
		return nil
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

// HasValue returns the presence of the contained value
func (c Custom[T, V]) HasValue() bool { return c.hasValue }

// Get returns the contained value and a boolean stating its presence.
// True if value exists, false otherwise.
//
// Panics if value was not validated yet
func (c Custom[T, V]) Get() (T, bool) {
	if c.hasValue && !c.validated {
		panic("called Value() on unvalidated value")
	}

	return c.value, c.hasValue
}

// Must returns the contained value and panics if it does not have one.
// You can check for its presence using [Custom.HasValue] or use a more safe alternative [Custom.Get]
func (c Custom[T, V]) Must() T {
	if !c.hasValue {
		panic("called must on empty optional")
	}

	value, _ := c.Get()

	return value
}

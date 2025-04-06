package optional

import (
	"encoding/json"
	"fmt"

	"github.com/metafates/schema/constraint"
	"github.com/metafates/schema/validate"
)

type (
	// Custom optional type.
	// When given not-null value it errors if validation fails
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
)

// HasValue returns the presence of the contained value
func (c Custom[T, V]) HasValue() bool { return c.hasValue }

// Value returns the contained value and a boolean stating its presence.
// True if value exists, false otherwise.
func (c Custom[T, V]) Value() (T, bool) { return c.value, c.hasValue }

// Must returns the contained value and panics if it does not have one.
// You can check for its presence using [Custom.HasValue] or use a more safe alternative [Custom.Value]
func (c Custom[T, V]) Must() T {
	if c.hasValue {
		return c.value
	}

	panic("called must on empty optional")
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

	if err := (*new(V)).Validate(*value); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	*c = Custom[T, V]{value: *value, hasValue: true}

	return nil
}

func (c *Custom[T, V]) UnmarshalText(data []byte) error {
	return c.UnmarshalJSON(data)
}

func (c Custom[T, V]) MarshalJSON() ([]byte, error) {
	if c.hasValue {
		return json.Marshal(c.value)
	}

	return []byte("null"), nil
}

func (c Custom[T, V]) MarshalText() ([]byte, error) {
	return c.MarshalJSON()
}

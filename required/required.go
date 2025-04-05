package required

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/metafates/schema/constraint"
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
		value T
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
)

func (c Custom[T, V]) IsSchema() {}

// Value returns the contained value
func (c Custom[T, V]) Value() T { return c.value }

func (c *Custom[T, V]) UnmarshalJSON(data []byte) error {
	var value *T

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	if value == nil {
		return errors.New("required value is missing")
	}

	if err := (*new(V)).Validate(*value); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	c.value = *value

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

package validate

import (
	"encoding/base64"
	"errors"
	"net/mail"
	"net/netip"
	"net/url"
	"strings"
	"unicode"

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
	_ Validator[string] = (*URL[string])(nil)
	_ Validator[string] = (*IP[string])(nil)
	_ Validator[string] = (*Printable[string])(nil)
	_ Validator[string] = (*Base64[string])(nil)
	_ Validator[string] = (*ASCII[string])(nil)

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
	Email[T constraint.Text] struct{}

	// URL accepts a single url.
	// The url may be relative (a path, without a host) or absolute (starting with a scheme)
	URL[T constraint.Text] struct{}

	// IP accepts an IP address.
	// The address can be in dotted decimal ("192.0.2.1"), IPv6 ("2001:db8::68"), or IPv6 with a scoped addressing zone ("fe80::1cc0:3e8c:119f:c2e1%ens18")
	IP[T constraint.Text] struct{}

	// Printable accepts strings consisting of only printable runes.
	// See [unicode.IsPrint] for more information
	Printable[T constraint.Text] struct{}

	// Base64 accepts valid base64 encoded strings
	Base64[T constraint.Text] struct{}

	// ASCII accepts ascii-only strings
	ASCII[T constraint.Text] struct{}

	// PrintableASCII combines [Printable] and [ASCII]
	PrintableASCII[T constraint.Text] struct {
		Combined[ASCII[T], Printable[T], T]
	}

	// Combined is a meta validator that combines other validators.
	// Validators are called in the same order as type parameters.
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

func (URL[T]) Validate(value T) error {
	_, err := url.Parse(string(value))
	if err != nil {
		return ValidateError{inner: err}
	}

	return nil
}

func (IP[T]) Validate(value T) error {
	_, err := netip.ParseAddr(string(value))
	if err != nil {
		return ValidateError{inner: err}
	}

	return nil
}

func (Printable[T]) Validate(value T) error {
	contains := strings.ContainsFunc(string(value), func(r rune) bool {
		return !unicode.IsPrint(r)
	})

	if contains {
		return ValidateError{inner: errors.New("string contains unprintable character")}
	}

	return nil
}

func (Base64[T]) Validate(value T) error {
	// TODO: implement it without allocating buffer and converting to string

	_, err := base64.StdEncoding.DecodeString(string(value))
	if err != nil {
		return ValidateError{inner: err}
	}

	return nil
}

func (ASCII[T]) Validate(value T) error {
	for i := 0; i < len(value); i++ {
		if value[i] > unicode.MaxASCII {
			return ValidateError{inner: errors.New("string contains non-ascii character")}
		}
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

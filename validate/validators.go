package validate

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"mime"
	"net"
	"net/mail"
	"net/netip"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/metafates/schema/constraint"
	"github.com/metafates/schema/internal/validateutil/iso"
	"github.com/metafates/schema/internal/validateutil/uuid"
)

type (
	// Any accepts any value of T
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

	// Even accepts integers divisible by two
	Even[T constraint.Integer] struct{}

	// Odd accepts integers not divisible by two
	Odd[T constraint.Integer] struct{}

	// Email accepts a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>"
	Email[T constraint.Text] struct{}

	// URL accepts a single url.
	// The url may be relative (a path, without a host) or absolute (starting with a scheme)
	//
	// See also [HTTPURL]
	URL[T constraint.Text] struct{}

	// HTTPURL accepts a single http(s) url.
	//
	// See also [URL]
	HTTPURL[T constraint.Text] struct{}

	// IP accepts an IP address.
	// The address can be in dotted decimal ("192.0.2.1"), IPv6 ("2001:db8::68"), or IPv6 with a scoped addressing zone ("fe80::1cc0:3e8c:119f:c2e1%ens18")
	IP[T constraint.Text] struct{}

	// MAC accepts an IEEE 802 MAC-48, EUI-48, EUI-64, or a 20-octet IP over InfiniBand link-layer address
	MAC[T constraint.Text] struct{}

	// CIDR accepts CIDR notation IP address and prefix length, like "192.0.2.0/24" or "2001:db8::/32", as defined in RFC 4632 and RFC 4291
	CIDR[T constraint.Text] struct{}

	// Printable accepts strings consisting of only printable runes.
	// See [unicode.IsPrint] for more information
	Printable[T constraint.Text] struct{}

	// NonEmptyPrintable combines [NonEmpty] and [Printable]
	NonEmptyPrintable[T ~string] struct {
		And[T, NonEmpty[T], Printable[T]]
	}

	// Base64 accepts valid base64 encoded strings
	Base64[T constraint.Text] struct{}

	// ASCII accepts ascii-only strings
	ASCII[T constraint.Text] struct{}

	// PrintableASCII combines [Printable] and [ASCII]
	PrintableASCII[T constraint.Text] struct{}

	// Latitude accepts any number in the range [-90; 90]
	//
	// See also [Longitude]
	Latitude[T constraint.Real] struct{}

	// Longitude accepts any number in the range [-180; 180]
	//
	// See also [Latitude]
	Longitude[T constraint.Real] struct{}

	// InPast accepts any time before current timestamp
	//
	// See also [InFuture]
	InPast[T constraint.Time] struct{}

	// InFuture accepts any time after current timestamp
	//
	// See also [InPast]
	InFuture[T constraint.Time] struct{}

	// Unique accepts an array of unique comparable values
	Unique[S ~[]T, T comparable] struct{}

	// MIME accepts RFC 1521 mime type string
	MIME[T constraint.Text] struct{}

	// UUID accepts a properly formatted UUID in one of the following formats:
	//   xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	//   urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	//   xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	//   {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
	UUID[T constraint.Text] struct{}

	// JSON accepts valid json encoded text
	JSON[T constraint.Text] struct{}

	// CountryAlpha2 accepts case-insensitive ISO 3166 2-letter country code
	CountryAlpha2[T constraint.Text] struct{}

	// CountryAlpha2 accepts case-insensitive ISO 3166 3-letter country code
	CountryAlpha3[T constraint.Text] struct{}

	// CountryAlpha2 accepts either [CountryAlpha2] or [CountryAlpha3]
	Country[T constraint.Text] struct {
		Or[T, CountryAlpha2[T], CountryAlpha3[T]]
	}

	// CurrencyAlpha accepts case-insensitive ISO 4217 alphabetic currency code
	CurrencyAlpha[T constraint.Text] struct{}

	// And is a meta validator that combines other validators with AND operator.
	// Validators are called in the same order as specified by type parameters.
	//
	// See also [Or]
	And[T any, A Validator[T], B Validator[T]] struct{}

	// And is a meta validator that combines other validators with OR operator.
	// Validators are called in the same order as type parameters.
	//
	// See also [And]
	Or[T any, A Validator[T], B Validator[T]] struct{}
)

func (Any[T]) Validate(T) error {
	return nil
}

func (NonEmpty[T]) Validate(value T) error {
	var empty T

	if value == empty {
		return errors.New("empty value")
	}

	return nil
}

func (Positive[T]) Validate(value T) error {
	if value < 0 {
		return errors.New("negative value")
	}

	return nil
}

func (Negative[T]) Validate(value T) error {
	if value > 0 {
		return errors.New("positive value")
	}

	return nil
}

func (Even[T]) Validate(value T) error {
	if value%2 != 0 {
		return fmt.Errorf("value must be even")
	}

	return nil
}

func (Odd[T]) Validate(value T) error {
	if value%2 == 0 {
		return fmt.Errorf("value must be odd")
	}

	return nil
}

func (Email[T]) Validate(value T) error {
	_, err := mail.ParseAddress(string(value))
	if err != nil {
		return err
	}

	return nil
}

func (URL[T]) Validate(value T) error {
	_, err := url.Parse(string(value))
	if err != nil {
		return err
	}

	return nil
}

func (HTTPURL[T]) Validate(value T) error {
	u, err := url.Parse(string(value))
	if err != nil {
		return err
	}

	if u.Host == "" {
		return errors.New("empty host")
	}

	switch u.Scheme {
	case "http", "https":
		return nil

	default:
		return errors.New("non-http(s) scheme")
	}
}

func (IP[T]) Validate(value T) error {
	_, err := netip.ParseAddr(string(value))
	if err != nil {
		return err
	}

	return nil
}

func (MAC[T]) Validate(value T) error {
	_, err := net.ParseMAC(string(value))
	if err != nil {
		return err
	}

	return nil
}

func (CIDR[T]) Validate(value T) error {
	_, _, err := net.ParseCIDR(string(value))
	if err != nil {
		return err
	}

	return nil
}

func (Printable[T]) Validate(value T) error {
	contains := strings.ContainsFunc(string(value), func(r rune) bool {
		return !unicode.IsPrint(r)
	})

	if contains {
		return errors.New("string contains unprintable character")
	}

	return nil
}

func (Base64[T]) Validate(value T) error {
	// TODO: implement it without allocating buffer and converting to string

	_, err := base64.StdEncoding.DecodeString(string(value))
	if err != nil {
		return err
	}

	return nil
}

func (ASCII[T]) Validate(value T) error {
	for i := 0; i < len(value); i++ {
		if value[i] > unicode.MaxASCII {
			return errors.New("string contains non-ascii character")
		}
	}

	return nil
}

func (PrintableASCII[T]) Validate(value T) error {
	for i := 0; i < len(value); i++ {
		if value[i] > unicode.MaxASCII {
			return errors.New("string contains non-ascii character")
		}

		if !unicode.IsPrint(rune(value[i])) {
			return errors.New("string contains unprintable character")
		}
	}

	return nil
}

func (Latitude[T]) Validate(value T) error {
	abs := math.Abs(float64(value))

	if abs > 90 {
		return errors.New("invalid latitude")
	}

	return nil
}

func (Longitude[T]) Validate(value T) error {
	abs := math.Abs(float64(value))

	if abs > 180 {
		return errors.New("invalid longitude")
	}

	return nil
}

func (InPast[T]) Validate(value T) error {
	if value.Compare(time.Now()) > 0 {
		return errors.New("time is not in the past")
	}

	return nil
}

func (InFuture[T]) Validate(value T) error {
	if value.Compare(time.Now()) < 0 {
		return errors.New("time is not in the future")
	}

	return nil
}

func (Unique[S, T]) Validate(value S) error {
	visited := make(map[T]struct{})

	for _, v := range value {
		if _, ok := visited[v]; ok {
			return errors.New("duplicate value found")
		}

		visited[v] = struct{}{}
	}

	return nil
}

func (MIME[T]) Validate(value T) error {
	_, _, err := mime.ParseMediaType(string(value))
	if err != nil {
		return err
	}

	return nil
}

func (UUID[T]) Validate(value T) error {
	// converting to bytes is cheaper than vice versa
	if err := uuid.Validate([]byte(value)); err != nil {
		return err
	}

	return nil
}

func (JSON[T]) Validate(value T) error {
	if !json.Valid([]byte(value)) {
		return errors.New("invalid json")
	}

	return nil
}

func (CountryAlpha2[T]) Validate(value T) error {
	v := strings.ToLower(string(value))

	if _, ok := iso.CountryAlpha2[v]; !ok {
		return errors.New("unknown 2-letter country code")
	}

	return nil
}

func (CountryAlpha3[T]) Validate(value T) error {
	v := strings.ToLower(string(value))

	if _, ok := iso.CountryAlpha3[v]; !ok {
		return errors.New("unknown 3-letter country code")
	}

	return nil
}

func (CurrencyAlpha[T]) Validate(value T) error {
	v := strings.ToLower(string(value))

	if _, ok := iso.CurrencyAlpha[v]; !ok {
		return errors.New("unknown currency alphabetic code")
	}

	return nil
}

func (And[T, A, B]) Validate(value T) error {
	if err := (*new(A)).Validate(value); err != nil {
		return err
	}

	if err := (*new(B)).Validate(value); err != nil {
		return err
	}

	return nil
}

func (Or[T, A, B]) Validate(value T) error {
	errA := (*new(A)).Validate(value)
	if errA == nil {
		return nil
	}

	errB := (*new(B)).Validate(value)
	if errB == nil {
		return nil
	}

	return errors.Join(errA, errB)
}

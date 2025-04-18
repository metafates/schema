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

	"github.com/metafates/schema/constraint"
	"github.com/metafates/schema/internal/iso"
	"github.com/metafates/schema/internal/uuid"
	"github.com/metafates/schema/validate/charset"
)

type (
	// Any accepts any value of T.
	Any[T any] struct{}

	// Zero accepts all zero values.
	//
	// The zero value is:
	// - 0 for numeric types,
	// - false for the boolean type, and
	// - "" (the empty string) for strings.
	//
	// See [NonZero]
	Zero[T comparable] struct{}

	// NonZero accepts all non-zero values.
	//
	// The zero value is:
	// - 0 for numeric types,
	// - false for the boolean type, and
	// - "" (the empty string) for strings.
	//
	// See [Zero]
	NonZero[T comparable] struct{}

	// Positive accepts all positive real numbers excluding zero.
	//
	// See [Positive0] for zero inlcuding variant.
	Positive[T constraint.Real] struct{}

	// Negative accepts all negative real numbers excluding zero.
	//
	// See [Negative0] for zero including variant.
	Negative[T constraint.Real] struct{}

	// Positive0 accepts all positive real numbers including zero.
	//
	// See [Positive] for zero excluding variant.
	Positive0[T constraint.Real] struct {
		Or[T, Positive[T], Zero[T]]
	}

	// Negative0 accepts all negative real numbers including zero.
	//
	// See [Negative] for zero excluding variant.
	Negative0[T constraint.Real] struct {
		Or[T, Negative[T], Zero[T]]
	}

	// Even accepts integers divisible by two.
	Even[T constraint.Integer] struct{}

	// Odd accepts integers not divisible by two.
	Odd[T constraint.Integer] struct{}

	// Email accepts a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>".
	Email[T constraint.Text] struct{}

	// URL accepts a single url.
	// The url may be relative (a path, without a host) or absolute (starting with a scheme).
	//
	// See also [HTTPURL].
	URL[T constraint.Text] struct{}

	// HTTPURL accepts a single http(s) url.
	//
	// See also [URL].
	HTTPURL[T constraint.Text] struct{}

	// IP accepts an IP address.
	// The address can be in dotted decimal ("192.0.2.1"), IPv6 ("2001:db8::68"), or IPv6 with a scoped addressing zone ("fe80::1cc0:3e8c:119f:c2e1%ens18").
	IP[T constraint.Text] struct{}

	// IP accepts an IP V4 address (e.g. "192.0.2.1").
	IPV4[T constraint.Text] struct{}

	// IP accepts an IP V6 address, including IPv4-mapped IPv6 addresses.
	// The address can be regular IPv6 ("2001:db8::68"), or IPv6 with a scoped addressing zone ("fe80::1cc0:3e8c:119f:c2e1%ens18").
	IPV6[T constraint.Text] struct{}

	// MAC accepts an IEEE 802 MAC-48, EUI-48, EUI-64, or a 20-octet IP over InfiniBand link-layer address.
	MAC[T constraint.Text] struct{}

	// CIDR accepts CIDR notation IP address and prefix length, like "192.0.2.0/24" or "2001:db8::/32", as defined in RFC 4632 and RFC 4291.
	CIDR[T constraint.Text] struct{}

	// Base64 accepts valid base64 encoded strings.
	Base64[T constraint.Text] struct{}

	// Charset0 accepts (possibly empty) text which contains only runes acceptable by filter.
	//
	// See [Charset] for a non-empty variant.
	Charset0[T constraint.Text, F charset.Filter] struct{}

	// Charset accepts non-empty text which contains only runes acceptable by filter.
	Charset[T constraint.Text, F charset.Filter] struct{}

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

	// Unique accepts a slice-like of unique values
	//
	// See [UniqueSlice] for a slice shortcut
	Unique[S ~[]T, T comparable] struct{}

	// Unique accepts a slice of unique values
	//
	// See [Unique] for a more generic version
	UniqueSlice[T comparable] struct {
		Unique[[]T, T]
	}

	// NonEmpty accepts a non-empty slice-like (len > 0)
	//
	// See [NonEmptySlice] for a slice shortcut
	NonEmpty[S ~[]T, T any] struct{}

	// NonEmpty accepts a non-empty slice (len > 0)
	//
	// See [NonEmpty] for a more generic version
	NonEmptySlice[T any] struct {
		NonEmpty[[]T, T]
	}

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
	CountryAlpha[T constraint.Text] struct {
		Or[T, CountryAlpha2[T], CountryAlpha3[T]]
	}

	// CurrencyAlpha accepts case-insensitive ISO 4217 alphabetic currency code
	CurrencyAlpha[T constraint.Text] struct{}

	// And is a meta validator that combines other validators with AND operator.
	// Validators are called in the same order as specified by type parameters.
	//
	// See also [Or], [Not]
	And[T any, A Validator[T], B Validator[T]] struct{}

	// And is a meta validator that combines other validators with OR operator.
	// Validators are called in the same order as type parameters.
	//
	// See also [And], [Not]
	Or[T any, A Validator[T], B Validator[T]] struct{}

	// Not is a meta validator that inverts given validator.
	//
	// See also [And], [Or]
	Not[T any, V Validator[T]] struct{}
)

func (Any[T]) Validate(T) error {
	return nil
}

func (Zero[T]) Validate(value T) error {
	var empty T

	if value != empty {
		return errors.New("non-zero value")
	}

	return nil
}

func (NonZero[T]) Validate(value T) error {
	var empty T

	if value == empty {
		return errors.New("zero value")
	}

	return nil
}

func (Positive[T]) Validate(value T) error {
	if value < 0 {
		return errors.New("negative value")
	}

	if value == 0 {
		return errors.New("zero value")
	}

	return nil
}

func (Negative[T]) Validate(value T) error {
	if value > 0 {
		return errors.New("positive value")
	}

	if value == 0 {
		return errors.New("zero value")
	}

	return nil
}

func (Even[T]) Validate(value T) error {
	if value%2 != 0 {
		return fmt.Errorf("odd value")
	}

	return nil
}

func (Odd[T]) Validate(value T) error {
	if value%2 == 0 {
		return fmt.Errorf("even value")
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

func (IPV4[T]) Validate(value T) error {
	a, err := netip.ParseAddr(string(value))
	if err != nil {
		return err
	}

	if !a.Is4() {
		return errors.New("ipv6 address")
	}

	return nil
}

func (IPV6[T]) Validate(value T) error {
	a, err := netip.ParseAddr(string(value))
	if err != nil {
		return err
	}

	if !a.Is6() {
		return errors.New("ipv6 address")
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

func (Base64[T]) Validate(value T) error {
	// TODO: implement it without allocating buffer and converting to string

	_, err := base64.StdEncoding.DecodeString(string(value))
	if err != nil {
		return err
	}

	return nil
}

func (Charset0[T, F]) Validate(value T) error {
	var f F

	for _, r := range string(value) {
		if err := f.Filter(r); err != nil {
			return err
		}
	}

	return nil
}

func (Charset[T, F]) Validate(value T) error {
	if len(value) == 0 {
		return errors.New("empty text")
	}

	return (*new(Charset0[T, F])).Validate(value)
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

func (NonEmpty[S, T]) Validate(value S) error {
	if len(value) == 0 {
		return errors.New("empty slice")
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
	if err := uuid.Validate(string(value)); err != nil {
		return err
	}

	return nil
}

func (JSON[T]) Validate(value T) error {
	if !json.Valid([]byte(string(value))) {
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

func (Not[T, V]) Validate(value T) error {
	var v V

	if err := v.Validate(value); err != nil {
		return nil
	}

	return errors.New(fmt.Sprint(v))
}

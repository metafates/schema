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
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/metafates/schema/constraint"
	"github.com/metafates/schema/internal/reflectwalk"
)

type (
	// Validator is an interface that validators must implement
	Validator[T any] interface {
		Validate(value T) error
	}

	// Validateable is an interface that states "this type can validate itself".
	// It is invoked for all values in the struct/slice/map/... by [Recursively]
	Validateable interface {
		Validate() error
	}
)

// ensure that all validators implement [Validator] interface
// if something is missing - please contribute!
var (
	_ Validator[any]       = (*Any[any])(nil)
	_ Validator[any]       = (*NonEmpty[any])(nil)
	_ Validator[int]       = (*Positive[int])(nil)
	_ Validator[int]       = (*Negative[int])(nil)
	_ Validator[string]    = (*Email[string])(nil)
	_ Validator[string]    = (*URL[string])(nil)
	_ Validator[string]    = (*IP[string])(nil)
	_ Validator[string]    = (*MAC[string])(nil)
	_ Validator[string]    = (*CIDR[string])(nil)
	_ Validator[string]    = (*Printable[string])(nil)
	_ Validator[string]    = (*PrintableASCII[string])(nil)
	_ Validator[string]    = (*Base64[string])(nil)
	_ Validator[string]    = (*ASCII[string])(nil)
	_ Validator[int]       = (*Latitude[int])(nil)
	_ Validator[int]       = (*Longitude[int])(nil)
	_ Validator[time.Time] = (*InPast[time.Time])(nil)
	_ Validator[time.Time] = (*InFuture[time.Time])(nil)
	_ Validator[[]string]  = (*Unique[[]string, string])(nil)
	_ Validator[string]    = (*MIME[string])(nil)

	_ Validator[any] = (*And[any, Validator[any], Validator[any]])(nil)
	_ Validator[any] = (*Or[any, Validator[any], Validator[any]])(nil)
)

// OnUnmarshal is a type that validates its inner value as part of unmarshalling
type OnUnmarshal[T any] struct{ Inner T }

func (w *OnUnmarshal[T]) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &w.Inner); err != nil {
		return err
	}

	if err := Validate(&w.Inner); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

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

// Validate traverses all fields in the given value pointed to by v, calling [Validateable.Validate] for each field
// and returns [ValidationError] if validation fails.
//
// If v is nil or not a pointer, Validate returns an [InvalidValidateError].
func Validate(v any) error {
	// same thing [json.Unmarshal] does
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidValidateError{Type: reflect.TypeOf(v)}
	}

	// we can skip fields traversal in case v implements [Validateable].
	if validatable, ok := v.(Validateable); ok {
		if err := validatable.Validate(); err != nil {
			return ValidationError{Inner: err}
		}

		return nil
	}

	return reflectwalk.WalkFields(v, func(path string, value reflect.Value) error {
		if value.CanAddr() {
			value = value.Addr()
		}

		r, ok := value.Interface().(Validateable)
		if !ok {
			return nil
		}

		if err := r.Validate(); err != nil {
			return ValidationError{
				path:  path,
				Inner: err,
			}
		}

		return nil
	})
}

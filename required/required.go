// Package required provides types whose values must be present and pass validation.
//
// Required types support the following encoding/decoding formats:
//   - json
//   - sql
//   - text
package required

import (
	"reflect"

	"github.com/metafates/schema/constraint"
	"github.com/metafates/schema/parse"
	"github.com/metafates/schema/validate"
	"github.com/metafates/schema/validate/charset"
)

var (
	ErrMissingValue  = validate.ValidationError{Msg: "missing required value"}
	ErrParseNilValue = parse.ParseError{Msg: "nil value passed for parsing"}
)

type (
	// Custom required type.
	// Errors if value is missing or did not pass the validation.
	Custom[T any, V validate.Validator[T]] struct {
		value     T
		hasValue  bool
		validated bool
	}

	// Any accepts any value of T.
	Any[T any] = Custom[T, validate.Any[T]]

	// NonZero accepts all non-zero values.
	//
	// The zero value is:
	// 	- 0 for numeric types,
	// 	- false for the boolean type, and
	// 	- "" (the empty string) for strings.
	NonZero[T comparable] = Custom[T, validate.NonZero[T]]

	// Positive accepts all positive real numbers and zero.
	//
	// See also [Negative].
	Positive[T constraint.Real] = Custom[T, validate.Positive[T]]

	// Negative accepts all negative real numbers and zero.
	//
	// See also [Positive].
	Negative[T constraint.Real] = Custom[T, validate.Negative[T]]

	// Positive0 accepts all positive real numbers including zero.
	//
	// See [Positive] for zero excluding variant.
	Positive0[T constraint.Real] = Custom[T, validate.Positive0[T]]

	// Negative0 accepts all negative real numbers including zero.
	//
	// See [Negative] for zero excluding variant.
	Negative0[T constraint.Real] = Custom[T, validate.Negative0[T]]

	// Even accepts real numbers divisible by two.
	Even[T constraint.Integer] = Custom[T, validate.Even[T]]

	// Odd accepts real numbers not divisible by two.
	Odd[T constraint.Integer] = Custom[T, validate.Odd[T]]

	// Email accepts a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>".
	Email[T constraint.Text] = Custom[T, validate.Email[T]]

	// URL accepts a single url.
	// The url may be relative (a path, without a host) or absolute (starting with a scheme).
	//
	// See also [HTTPURL].
	URL[T constraint.Text] = Custom[T, validate.URL[T]]

	// HTTPURL accepts a single http(s) url.
	//
	// See also [URL].
	HTTPURL[T constraint.Text] = Custom[T, validate.HTTPURL[T]]

	// IP accepts an IP address.
	// The address can be in dotted decimal ("192.0.2.1"),
	// IPv6 ("2001:db8::68"), or IPv6 with a scoped addressing zone ("fe80::1cc0:3e8c:119f:c2e1%ens18").
	IP[T constraint.Text] = Custom[T, validate.IP[T]]

	// IPV4 accepts an IP V4 address (e.g. "192.0.2.1").
	IPV4[T constraint.Text] = Custom[T, validate.IPV4[T]]

	// IPV6 accepts an IP V6 address, including IPv4-mapped IPv6 addresses.
	// The address can be regular IPv6 ("2001:db8::68"), or IPv6 with
	// a scoped addressing zone ("fe80::1cc0:3e8c:119f:c2e1%ens18").
	IPV6[T constraint.Text] = Custom[T, validate.IPV6[T]]

	// MAC accepts an IEEE 802 MAC-48, EUI-48, EUI-64, or a 20-octet IP over InfiniBand link-layer address.
	MAC[T constraint.Text] = Custom[T, validate.MAC[T]]

	// CIDR accepts CIDR notation IP address and prefix length,
	// like "192.0.2.0/24" or "2001:db8::/32", as defined in RFC 4632 and RFC 4291.
	CIDR[T constraint.Text] = Custom[T, validate.CIDR[T]]

	// Base64 accepts valid base64 encoded strings.
	Base64[T constraint.Text] = Custom[T, validate.Base64[T]]

	// Charset0 accepts (possibly empty) text which contains only runes acceptable by filter.
	//
	// See [Charset] for a non-empty variant.
	Charset0[T constraint.Text, F charset.Filter] = Custom[T, validate.Charset0[T, F]]

	// Charset accepts non-empty text which contains only runes acceptable by filter.
	Charset[T constraint.Text, F charset.Filter] = Custom[T, validate.Charset[T, F]]

	// Latitude accepts any number in the range [-90; 90].
	//
	// See also [Longitude].
	Latitude[T constraint.Real] = Custom[T, validate.Latitude[T]]

	// Longitude accepts any number in the range [-180; 180].
	//
	// See also [Latitude].
	Longitude[T constraint.Real] = Custom[T, validate.Longitude[T]]

	// InPast accepts any time before current timestamp.
	//
	// See also [InFuture].
	InPast[T constraint.Time] = Custom[T, validate.InPast[T]]

	// InFuture accepts any time after current timestamp.
	//
	// See also [InPast].
	InFuture[T constraint.Time] = Custom[T, validate.InFuture[T]]

	// Unique accepts a slice-like of unique values.
	//
	// See [UniqueSlice] for a slice shortcut.
	Unique[S ~[]T, T comparable] = Custom[S, validate.Unique[S, T]]

	// Unique accepts a slice of unique values.
	//
	// See [Unique] for a more generic version.
	UniqueSlice[T comparable] = Custom[[]T, validate.UniqueSlice[T]]

	// NonEmpty accepts a non-empty slice-like (len > 0).
	//
	// See [NonEmptySlice] for a slice shortcut.
	NonEmpty[S ~[]T, T any] = Custom[S, validate.NonEmpty[S, T]]

	// NonEmpty accepts a non-empty slice (len > 0).
	//
	// See [NonEmpty] for a more generic version.
	NonEmptySlice[T any] = Custom[[]T, validate.NonEmptySlice[T]]

	// MIME accepts RFC 1521 mime type string.
	MIME[T constraint.Text] = Custom[T, validate.MIME[T]]

	// UUID accepts a properly formatted UUID in one of the following formats:
	//   xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	//   urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	//   xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	//   {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
	UUID[T constraint.Text] = Custom[T, validate.UUID[T]]

	// JSON accepts valid json encoded text.
	JSON[T constraint.Text] = Custom[T, validate.JSON[T]]

	// CountryAlpha2 accepts case-insensitive ISO 3166 2-letter country code.
	CountryAlpha2[T constraint.Text] = Custom[T, validate.CountryAlpha2[T]]

	// CountryAlpha3 accepts case-insensitive ISO 3166 3-letter country code.
	CountryAlpha3[T constraint.Text] = Custom[T, validate.CountryAlpha3[T]]

	// CountryAlpha accepts either [CountryAlpha2] or [CountryAlpha3].
	CountryAlpha[T constraint.Text] = Custom[T, validate.CountryAlpha[T]]

	// CurrencyAlpha accepts case-insensitive ISO 4217 alphabetic currency code.
	CurrencyAlpha[T constraint.Text] = Custom[T, validate.CurrencyAlpha[T]]

	// LangAlpha2 accepts case-insesitive ISO 639 2-letter language code.
	LangAlpha2[T constraint.Text] = Custom[T, validate.LangAlpha2[T]]

	// LangAlpha2 accepts case-insesitive ISO 639 3-letter language code.
	LangAlpha3[T constraint.Text] = Custom[T, validate.LangAlpha3[T]]

	// LangAlpha accepts either [LangAlpha2] or [LangAlpha3].
	LangAlpha[T constraint.Text] = Custom[T, validate.LangAlpha[T]]
)

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

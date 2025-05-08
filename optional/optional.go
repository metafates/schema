// Package optional provides types whose values may be either empty (null) or be present and pass validation.
//
// Optional types support the following encoding/decoding formats:
//   - json
//   - sql
//   - text
//   - binary
//   - gob
package optional

import (
	"reflect"

	"github.com/metafates/schema/constraint"
	"github.com/metafates/schema/parse"
	"github.com/metafates/schema/validate"
	"github.com/metafates/schema/validate/charset"
)

type (
	// Custom optional type.
	// When given non-null value it errors if validation fails.
	Custom[T any, V validate.Validator[T]] struct {
		value     T
		hasValue  bool
		validated bool
	}

	// Any accepts any value of T.
	Any[T any] struct {
		Custom[T, validate.Any[T]]
	}

	// NonZero accepts all non-zero values.
	//
	// The zero value is:
	// - 0 for numeric types,
	// - false for the boolean type, and
	// - "" (the empty string) for strings.
	NonZero[T comparable] struct {
		Custom[T, validate.NonZero[T]]
	}

	// Positive accepts all positive real numbers and zero.
	//
	// See also [Negative].
	Positive[T constraint.Real] struct {
		Custom[T, validate.Positive[T]]
	}

	// Negative accepts all negative real numbers and zero.
	//
	// See also [Positive].
	Negative[T constraint.Real] struct {
		Custom[T, validate.Negative[T]]
	}

	// Positive0 accepts all positive real numbers including zero.
	//
	// See [Positive] for zero excluding variant.
	Positive0[T constraint.Real] struct {
		Custom[T, validate.Positive0[T]]
	}

	// Negative0 accepts all negative real numbers including zero.
	//
	// See [Negative] for zero excluding variant.
	Negative0[T constraint.Real] struct {
		Custom[T, validate.Negative0[T]]
	}

	// Even accepts real numbers divisible by two.
	Even[T constraint.Integer] struct {
		Custom[T, validate.Even[T]]
	}

	// Odd accepts real numbers not divisible by two.
	Odd[T constraint.Integer] struct {
		Custom[T, validate.Odd[T]]
	}

	// Email accepts a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>".
	Email[T constraint.Text] struct {
		Custom[T, validate.Email[T]]
	}

	// URL accepts a single url.
	// The url may be relative (a path, without a host) or absolute (starting with a scheme).
	//
	// See also [HTTPURL].
	URL[T constraint.Text] struct {
		Custom[T, validate.URL[T]]
	}

	// HTTPURL accepts a single http(s) url.
	//
	// See also [URL].
	HTTPURL[T constraint.Text] struct {
		Custom[T, validate.HTTPURL[T]]
	}

	// IP accepts an IP address.
	// The address can be in dotted decimal ("192.0.2.1"),
	// IPv6 ("2001:db8::68"), or IPv6 with a scoped addressing zone ("fe80::1cc0:3e8c:119f:c2e1%ens18").
	IP[T constraint.Text] struct {
		Custom[T, validate.IP[T]]
	}

	// IPV4 accepts an IP V4 address (e.g. "192.0.2.1").
	IPV4[T constraint.Text] struct {
		Custom[T, validate.IPV4[T]]
	}

	// IPV6 accepts an IP V6 address, including IPv4-mapped IPv6 addresses.
	// The address can be regular IPv6 ("2001:db8::68"), or IPv6 with
	// a scoped addressing zone ("fe80::1cc0:3e8c:119f:c2e1%ens18").
	IPV6[T constraint.Text] struct {
		Custom[T, validate.IPV6[T]]
	}

	// MAC accepts an IEEE 802 MAC-48, EUI-48, EUI-64, or a 20-octet IP over InfiniBand link-layer address.
	MAC[T constraint.Text] struct {
		Custom[T, validate.MAC[T]]
	}

	// CIDR accepts CIDR notation IP address and prefix length,
	// like "192.0.2.0/24" or "2001:db8::/32", as defined in RFC 4632 and RFC 4291.
	CIDR[T constraint.Text] struct {
		Custom[T, validate.CIDR[T]]
	}

	// Base64 accepts valid base64 encoded strings.
	Base64[T constraint.Text] struct {
		Custom[T, validate.Base64[T]]
	}

	// Charset accepts text which contains only runes acceptable by filter.
	//
	// NOTE: empty strings will also pass. Use [NonZeroCharset] if you need non-empty strings.
	Charset[T constraint.Text, F charset.Filter] struct {
		Custom[T, validate.Charset0[T, F]]
	}

	// NonZeroCharset combines [NonZero] and [Charset].
	NonZeroCharset[T constraint.Text, F charset.Filter] struct {
		Custom[T, validate.Charset[T, F]]
	}

	// Latitude accepts any number in the range [-90; 90].
	//
	// See also [Longitude].
	Latitude[T constraint.Real] struct {
		Custom[T, validate.Latitude[T]]
	}

	// Longitude accepts any number in the range [-180; 180].
	//
	// See also [Latitude].
	Longitude[T constraint.Real] struct {
		Custom[T, validate.Longitude[T]]
	}

	// InPast accepts any time before current timestamp.
	//
	// See also [InFuture].
	InPast[T constraint.Time] struct {
		Custom[T, validate.InPast[T]]
	}

	// InFuture accepts any time after current timestamp.
	//
	// See also [InPast].
	InFuture[T constraint.Time] struct {
		Custom[T, validate.InFuture[T]]
	}

	// Unique accepts a slice-like of unique values.
	//
	// See [UniqueSlice] for a slice shortcut.
	Unique[S ~[]T, T comparable] struct {
		Custom[S, validate.Unique[S, T]]
	}

	// Unique accepts a slice of unique values.
	//
	// See [Unique] for a more generic version.
	UniqueSlice[T comparable] struct {
		Custom[[]T, validate.UniqueSlice[T]]
	}

	// NonEmpty accepts a non-empty slice-like (len > 0).
	//
	// See [NonEmptySlice] for a slice shortcut.
	NonEmpty[S ~[]T, T any] struct {
		Custom[S, validate.NonEmpty[S, T]]
	}

	// NonEmpty accepts a non-empty slice (len > 0).
	//
	// See [NonEmpty] for a more generic version.
	NonEmptySlice[T any] struct {
		Custom[[]T, validate.NonEmptySlice[T]]
	}

	// MIME accepts RFC 1521 mime type string.
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

	// JSON accepts valid json encoded text.
	JSON[T constraint.Text] struct {
		Custom[T, validate.JSON[T]]
	}

	// CountryAlpha2 accepts case-insensitive ISO 3166 2-letter country code.
	CountryAlpha2[T constraint.Text] struct {
		Custom[T, validate.CountryAlpha2[T]]
	}

	// CountryAlpha3 accepts case-insensitive ISO 3166 3-letter country code.
	CountryAlpha3[T constraint.Text] struct {
		Custom[T, validate.CountryAlpha3[T]]
	}

	// CountryAlpha accepts either [CountryAlpha2] or [CountryAlpha3].
	CountryAlpha[T constraint.Text] struct {
		Custom[T, validate.CountryAlpha[T]]
	}

	// CurrencyAlpha accepts case-insensitive ISO 4217 alphabetic currency code.
	CurrencyAlpha[T constraint.Text] struct {
		Custom[T, validate.CurrencyAlpha[T]]
	}

	// LangAlpha2 accepts case-insesitive ISO 639 2-letter language code.
	LangAlpha2[T constraint.Text] struct {
		Custom[T, validate.LangAlpha2[T]]
	}

	// LangAlpha2 accepts case-insesitive ISO 639 3-letter language code.
	LangAlpha3[T constraint.Text] struct {
		Custom[T, validate.LangAlpha3[T]]
	}

	// LangAlpha accepts either [LangAlpha2] or [LangAlpha3].
	LangAlpha[T constraint.Text] struct {
		Custom[T, validate.LangAlpha[T]]
	}
)

// TypeValidate implements the [validate.TypeValidateable] interface.
// You should not call this function directly.
func (c *Custom[T, V]) TypeValidate() error {
	if !c.hasValue {
		return nil
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

// HasValue returns the presence of the contained value.
func (c Custom[T, V]) HasValue() bool { return c.hasValue }

// Get returns the contained value and a boolean stating its presence.
// True if value exists, false otherwise.
//
// Panics if value was not validated yet.
// See also [Custom.GetPtr].
func (c Custom[T, V]) Get() (T, bool) {
	if c.hasValue && !c.validated {
		panic("called Get() on non-empty unvalidated value")
	}

	return c.value, c.hasValue
}

// Get returns the pointer to the contained value.
// Non-nil if value exists, nil otherwise.
// Pointed value is a shallow copy.
//
// Panics if value was not validated yet.
// See also [Custom.Get].
func (c Custom[T, V]) GetPtr() *T {
	if c.hasValue && !c.validated {
		panic("called GetPtr() on non-empty unvalidated value")
	}

	var value *T

	if c.hasValue {
		valueCopy := c.value
		value = &valueCopy
	}

	return value
}

// Must returns the contained value and panics if it does not have one.
// You can check for its presence using [Custom.HasValue] or use a more safe alternative [Custom.Get].
func (c Custom[T, V]) Must() T {
	if !c.hasValue {
		panic("called must on empty optional")
	}

	value, _ := c.Get()

	return value
}

// Parse checks if given value is valid.
// If it is, a value is used to initialize this type.
// Value is converted to the target type T, if possible. If not - [parse.UnconvertableTypeError] is returned.
// It is allowed to pass convertable type wrapped in optional type.
//
// Parsed type is validated, therefore it is safe to call [Custom.Get] afterwards.
//
// Passing nil results a valid empty instance.
func (c *Custom[T, V]) Parse(value any) error {
	if value == nil {
		*c = Custom[T, V]{}

		return nil
	}

	rValue := reflect.ValueOf(value)

	if rValue.Kind() == reflect.Pointer && rValue.IsNil() {
		*c = Custom[T, V]{}

		return nil
	}

	if _, ok := value.(interface{ isOptional() }); ok {
		// NOTE: ensure this method name is in sync with [Custom.Get]
		res := rValue.MethodByName("Get").Call(nil)
		v, ok := res[0], res[1].Bool()

		if !ok {
			*c = Custom[T, V]{}

			return nil
		}

		rValue = v
	}

	v, err := convert[T](rValue)
	if err != nil {
		return parse.ParseError{Inner: err}
	}

	aux := Custom[T, V]{
		hasValue: true,
		value:    v,
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

func (Custom[T, V]) isOptional() {}

func convert[T any](v reflect.Value) (T, error) {
	tType := reflect.TypeFor[T]()

	original := v

	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.CanConvert(tType) {
		//nolint:forcetypeassert // checked already by CanConvert
		return v.Convert(tType).Interface().(T), nil
	}

	return *new(T), parse.UnconvertableTypeError{
		Target:   tType.String(),
		Original: original.Type().String(),
	}
}

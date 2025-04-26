package optional_test

import (
	"fmt"
	"time"

	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/validate"
	"github.com/metafates/schema/validate/charset"
)

func ExampleAny() {
	var value optional.Any[string]

	// Parse a value
	value.MustParse("any value is acceptable")
	val, exists := value.Get()
	fmt.Printf("Value: %s, Exists: %t\n", val, exists)

	// Empty string is also valid
	value.MustParse("")
	val, exists = value.Get()
	fmt.Printf("Empty string value: %q, Exists: %t\n", val, exists)

	// Parse nil for an empty optional
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// GetPtr returns nil for empty optionals
	value.MustParse(nil)
	ptr := value.GetPtr()
	fmt.Printf("GetPtr returns nil: %t\n", ptr == nil)

	// Output:
	// Value: any value is acceptable, Exists: true
	// Empty string value: "", Exists: true
	// After nil: Exists: false
	// GetPtr returns nil: true
}

func ExampleNonZero() {
	var strValue optional.NonZero[string]

	// Valid non-zero value
	strValue.MustParse("non-empty string")
	val, exists := strValue.Get()
	fmt.Printf("Value: %s, Exists: %t\n", val, exists)

	// Empty optionals are valid
	strValue.MustParse(nil)
	fmt.Printf("HasValue after nil: %t\n", strValue.HasValue())

	// Empty strings will fail validation
	var invalidStr optional.NonZero[string]
	err := invalidStr.Parse("")
	fmt.Printf("Empty string error: %t\n", err != nil)

	// Zero number will fail validation
	var numValue optional.NonZero[int]
	err = numValue.Parse(0)
	fmt.Printf("Zero number error: %t\n", err != nil)

	// Must() will panic on empty optionals
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Must() panicked on empty optional as expected")
		}
	}()
	strValue.MustParse(nil)
	strValue.Must() // This will panic

	// Output:
	// Value: non-empty string, Exists: true
	// HasValue after nil: false
	// Empty string error: true
	// Zero number error: true
	// Must() panicked on empty optional as expected
}

func ExamplePositive() {
	var intValue optional.Positive[int]

	// Valid positive value
	intValue.MustParse(42)
	val, exists := intValue.Get()
	fmt.Printf("Value: %d, Exists: %t\n", val, exists)

	// Empty optional is valid
	intValue.MustParse(nil)
	val, exists = intValue.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Zero and negative values will fail validation
	var invalidValue optional.Positive[int]
	err := invalidValue.Parse(0)
	fmt.Printf("Zero error: %t\n", err != nil)

	err = invalidValue.Parse(-5)
	fmt.Printf("Negative error: %t\n", err != nil)

	// Output:
	// Value: 42, Exists: true
	// After nil: Exists: false
	// Zero error: true
	// Negative error: true
}

func ExampleNegative() {
	var intValue optional.Negative[int]

	// Valid negative value
	intValue.MustParse(-42)
	val, exists := intValue.Get()
	fmt.Printf("Value: %d, Exists: %t\n", val, exists)

	// Empty optional is valid
	intValue.MustParse(nil)
	val, exists = intValue.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Zero and positive values will fail validation
	var invalidValue optional.Negative[int]
	err := invalidValue.Parse(0)
	fmt.Printf("Zero error: %t\n", err != nil)

	err = invalidValue.Parse(5)
	fmt.Printf("Positive error: %t\n", err != nil)

	// Output:
	// Value: -42, Exists: true
	// After nil: Exists: false
	// Zero error: true
	// Positive error: true
}

func ExamplePositive0() {
	var intValue optional.Positive0[int]

	// Valid positive value
	intValue.MustParse(42)
	val, exists := intValue.Get()
	fmt.Printf("Value: %d, Exists: %t\n", val, exists)

	// Zero is valid for Positive0
	intValue.MustParse(0)
	val, exists = intValue.Get()
	fmt.Printf("Zero value: %d, Exists: %t\n", val, exists)

	// Empty optional is valid
	intValue.MustParse(nil)
	val, exists = intValue.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Negative values will fail validation
	var invalidValue optional.Positive0[int]
	err := invalidValue.Parse(-5)
	fmt.Printf("Negative error: %t\n", err != nil)

	// Output:
	// Value: 42, Exists: true
	// Zero value: 0, Exists: true
	// After nil: Exists: false
	// Negative error: true
}

func ExampleNegative0() {
	var intValue optional.Negative0[int]

	// Valid negative value
	intValue.MustParse(-42)
	val, exists := intValue.Get()
	fmt.Printf("Value: %d, Exists: %t\n", val, exists)

	// Zero is valid for Negative0
	intValue.MustParse(0)
	val, exists = intValue.Get()
	fmt.Printf("Zero value: %d, Exists: %t\n", val, exists)

	// Empty optional is valid
	intValue.MustParse(nil)
	val, exists = intValue.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Positive values will fail validation
	var invalidValue optional.Negative0[int]
	err := invalidValue.Parse(5)
	fmt.Printf("Positive error: %t\n", err != nil)

	// Output:
	// Value: -42, Exists: true
	// Zero value: 0, Exists: true
	// After nil: Exists: false
	// Positive error: true
}

func ExampleEven() {
	var value optional.Even[int]

	// Valid even value
	value.MustParse(42)
	val, exists := value.Get()
	fmt.Printf("Value: %d, Exists: %t\n", val, exists)

	// Zero is even
	value.MustParse(0)
	val, exists = value.Get()
	fmt.Printf("Zero value: %d, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Odd numbers will fail validation
	var invalidValue optional.Even[int]
	err := invalidValue.Parse(3)
	fmt.Printf("Odd number error: %t\n", err != nil)

	// Output:
	// Value: 42, Exists: true
	// Zero value: 0, Exists: true
	// After nil: Exists: false
	// Odd number error: true
}

func ExampleOdd() {
	var value optional.Odd[int]

	// Valid odd value
	value.MustParse(43)
	val, exists := value.Get()
	fmt.Printf("Value: %d, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Even numbers will fail validation
	var invalidValue optional.Odd[int]
	err := invalidValue.Parse(2)
	fmt.Printf("Even number error: %t\n", err != nil)

	err = invalidValue.Parse(0)
	fmt.Printf("Zero error: %t\n", err != nil)

	// Output:
	// Value: 43, Exists: true
	// After nil: Exists: false
	// Even number error: true
	// Zero error: true
}

func ExampleEmail() {
	var value optional.Email[string]

	// Valid email
	value.MustParse("user@example.com")
	val, exists := value.Get()
	fmt.Printf("Value: %s, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Invalid emails will fail validation
	var invalidValue optional.Email[string]
	err := invalidValue.Parse("not-an-email")
	fmt.Printf("Invalid email error: %t\n", err != nil)

	// Using GetPtr()
	value.MustParse("user@example.com")
	ptr := value.GetPtr()
	fmt.Printf("GetPtr value: %s, IsNil: %t\n", *ptr, ptr == nil)

	// Output:
	// Value: user@example.com, Exists: true
	// After nil: Exists: false
	// Invalid email error: true
	// GetPtr value: user@example.com, IsNil: false
}

func ExampleURL() {
	var value optional.URL[string]

	// Valid URL
	value.MustParse("https://example.com")
	val, exists := value.Get()
	fmt.Printf("Value: %s, Exists: %t\n", val, exists)

	// Relative URLs are valid too
	value.MustParse("/relative/path")
	val, exists = value.Get()
	fmt.Printf("Relative URL: %s, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Invalid URLs will fail validation
	var invalidValue optional.URL[string]
	err := invalidValue.Parse("htt ps://com example")
	fmt.Printf("Invalid URL error: %t\n", err != nil)

	// Output:
	// Value: https://example.com, Exists: true
	// Relative URL: /relative/path, Exists: true
	// After nil: Exists: false
	// Invalid URL error: true
}

func ExampleHTTPURL() {
	var value optional.HTTPURL[string]

	// Valid HTTP URL
	value.MustParse("https://example.com")
	val, exists := value.Get()
	fmt.Printf("Value: %s, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Non-HTTP URLs will fail validation
	var invalidValue optional.HTTPURL[string]
	err := invalidValue.Parse("ftp://example.com")
	fmt.Printf("Non-HTTP URL error: %t\n", err != nil)

	// Relative URLs will fail validation
	err = invalidValue.Parse("/relative/path")
	fmt.Printf("Relative URL error: %t\n", err != nil)

	// Output:
	// Value: https://example.com, Exists: true
	// After nil: Exists: false
	// Non-HTTP URL error: true
	// Relative URL error: true
}

func ExampleIP() {
	var value optional.IP[string]

	// Valid IPv4
	value.MustParse("192.168.1.1")
	val, exists := value.Get()
	fmt.Printf("IPv4: %s, Exists: %t\n", val, exists)

	// Valid IPv6
	value.MustParse("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	val, exists = value.Get()
	fmt.Printf("Has IPv6: %t\n", exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Invalid IPs will fail validation
	var invalidValue optional.IP[string]
	err := invalidValue.Parse("999.999.999.999")
	fmt.Printf("Invalid IP error: %t\n", err != nil)

	// Output:
	// IPv4: 192.168.1.1, Exists: true
	// Has IPv6: true
	// After nil: Exists: false
	// Invalid IP error: true
}

func ExampleIPV4() {
	var value optional.IPV4[string]

	// Valid IPv4
	value.MustParse("192.168.1.1")
	val, exists := value.Get()
	fmt.Printf("Value: %s, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// IPv6 addresses will fail validation
	var invalidValue optional.IPV4[string]
	err := invalidValue.Parse("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	fmt.Printf("IPv6 address error: %t\n", err != nil)

	// Output:
	// Value: 192.168.1.1, Exists: true
	// After nil: Exists: false
	// IPv6 address error: true
}

func ExampleIPV6() {
	var value optional.IPV6[string]

	// Valid IPv6
	value.MustParse("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	_, exists := value.Get()
	fmt.Printf("Has value: %t\n", exists)

	// Empty optional is valid
	value.MustParse(nil)
	_, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// IPv4 addresses will fail validation
	var invalidValue optional.IPV6[string]
	err := invalidValue.Parse("192.168.1.1")
	fmt.Printf("IPv4 address error: %t\n", err != nil)

	// Output:
	// Has value: true
	// After nil: Exists: false
	// IPv4 address error: true
}

func ExampleMAC() {
	var value optional.MAC[string]

	// Valid MAC address
	value.MustParse("00:1B:44:11:3A:B7")
	val, exists := value.Get()
	fmt.Printf("Value: %s, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Invalid MAC addresses will fail validation
	var invalidValue optional.MAC[string]
	err := invalidValue.Parse("not a mac address")
	fmt.Printf("Invalid MAC error: %t\n", err != nil)

	// Output:
	// Value: 00:1B:44:11:3A:B7, Exists: true
	// After nil: Exists: false
	// Invalid MAC error: true
}

func ExampleCIDR() {
	var value optional.CIDR[string]

	// Valid CIDR
	value.MustParse("192.168.1.0/24")
	val, exists := value.Get()
	fmt.Printf("Value: %s, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Invalid CIDRs will fail validation
	var invalidValue optional.CIDR[string]
	err := invalidValue.Parse("192.168.1.1") // missing prefix
	fmt.Printf("Missing prefix error: %t\n", err != nil)

	// Output:
	// Value: 192.168.1.0/24, Exists: true
	// After nil: Exists: false
	// Missing prefix error: true
}

func ExampleBase64() {
	var value optional.Base64[string]

	// Valid base64
	value.MustParse("SGVsbG8gV29ybGQ=") // "Hello World"
	_, exists := value.Get()
	fmt.Printf("Value exists: %t\n", exists)

	// Empty optional is valid
	value.MustParse(nil)
	_, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Invalid base64 will fail validation
	var invalidValue optional.Base64[string]
	err := invalidValue.Parse("not!valid!base64!")

	fmt.Printf("Invalid base64 error: %t\n", err != nil)

	// Output:
	// Value exists: true
	// After nil: Exists: false
	// Invalid base64 error: true
}

func ExampleCharset() {
	// Charset allows empty strings
	var value optional.Charset[string, charset.Letter]

	value.MustParse("abcDEF")
	val, exists := value.Get()
	fmt.Printf("Value: %s, Exists: %t\n", val, exists)

	// Empty strings are allowed
	value.MustParse("")
	val, exists = value.Get()
	fmt.Printf("Empty string: %q, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Non-alphabetic characters will fail validation
	var invalidValue optional.Charset[string, charset.Letter]

	err := invalidValue.Parse("abc123")
	fmt.Printf("Non-alphabetic error: %t\n", err != nil)

	// Output:
	// Value: abcDEF, Exists: true
	// Empty string: "", Exists: true
	// After nil: Exists: false
	// Non-alphabetic error: true
}

func ExampleNonZeroCharset() {
	// NonZeroCharset requires non-empty strings
	var value optional.NonZeroCharset[string, charset.Or[charset.Letter, charset.Number]]

	value.MustParse("abc123DEF")
	val, exists := value.Get()
	fmt.Printf("Value: %s, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Empty strings will fail validation
	var invalidValue optional.NonZeroCharset[string, charset.Or[charset.Letter, charset.Number]]

	err := invalidValue.Parse("")
	fmt.Printf("Empty string error: %t\n", err != nil)

	// Non-alphanumeric characters will fail validation
	err = invalidValue.Parse("abc123!")
	fmt.Printf("Non-alphanumeric error: %t\n", err != nil)

	// Output:
	// Value: abc123DEF, Exists: true
	// After nil: Exists: false
	// Empty string error: true
	// Non-alphanumeric error: true
}

func ExampleLatitude() {
	var value optional.Latitude[float64]

	// Valid latitude
	value.MustParse(37.7749) // San Francisco
	val, exists := value.Get()
	fmt.Printf("Value: %.4f, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Out of range latitudes will fail validation
	var invalidValue optional.Latitude[float64]
	err := invalidValue.Parse(91)
	fmt.Printf("Too high error: %t\n", err != nil)

	// Output:
	// Value: 37.7749, Exists: true
	// After nil: Exists: false
	// Too high error: true
}

func ExampleLongitude() {
	var value optional.Longitude[float64]

	// Valid longitude
	value.MustParse(-122.4194) // San Francisco
	val, exists := value.Get()
	fmt.Printf("Value: %.4f, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Out of range longitudes will fail validation
	var invalidValue optional.Longitude[float64]
	err := invalidValue.Parse(181)
	fmt.Printf("Too high error: %t\n", err != nil)

	// Output:
	// Value: -122.4194, Exists: true
	// After nil: Exists: false
	// Too high error: true
}

func ExampleInPast() {
	var value optional.InPast[time.Time]

	// Valid past time
	pastTime := time.Now().Add(-24 * time.Hour) // 1 day ago
	value.MustParse(pastTime)
	val, exists := value.Get()
	fmt.Printf("Is in past: %t, Exists: %t\n", val.Before(time.Now()), exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Future times will fail validation
	var invalidValue optional.InPast[time.Time]
	futureTime := time.Now().Add(24 * time.Hour) // 1 day in future
	err := invalidValue.Parse(futureTime)
	fmt.Printf("Future time error: %t\n", err != nil)

	// Output:
	// Is in past: true, Exists: true
	// After nil: Exists: false
	// Future time error: true
}

func ExampleInFuture() {
	var value optional.InFuture[time.Time]

	// Valid future time
	futureTime := time.Now().Add(24 * time.Hour) // 1 day in future
	value.MustParse(futureTime)
	val, exists := value.Get()
	fmt.Printf("Is in future: %t, Exists: %t\n", val.After(time.Now()), exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Past times will fail validation
	var invalidValue optional.InFuture[time.Time]
	pastTime := time.Now().Add(-24 * time.Hour) // 1 day ago
	err := invalidValue.Parse(pastTime)
	fmt.Printf("Past time error: %t\n", err != nil)

	// Output:
	// Is in future: true, Exists: true
	// After nil: Exists: false
	// Past time error: true
}

func ExampleUnique() {
	// Unique with a slice type
	var value optional.Unique[[]int, int]
	value.MustParse([]int{1, 2, 3, 4, 5})
	val, exists := value.Get()
	fmt.Printf("Value: %v, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Duplicate values will fail validation
	var invalidValue optional.Unique[[]int, int]
	err := invalidValue.Parse([]int{1, 2, 3, 1, 4})
	fmt.Printf("Duplicate value error: %t\n", err != nil)

	// Output:
	// Value: [1 2 3 4 5], Exists: true
	// After nil: Exists: false
	// Duplicate value error: true
}

func ExampleUniqueSlice() {
	// UniqueSlice is a simplified version of Unique
	var value optional.UniqueSlice[string]
	value.MustParse([]string{"apple", "banana", "cherry"})
	val, exists := value.Get()
	fmt.Printf("Value: %v, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Duplicate values will fail validation
	var invalidValue optional.UniqueSlice[string]
	err := invalidValue.Parse([]string{"apple", "banana", "apple"})
	fmt.Printf("Duplicate value error: %t\n", err != nil)

	// Output:
	// Value: [apple banana cherry], Exists: true
	// After nil: Exists: false
	// Duplicate value error: true
}

func ExampleNonEmpty() {
	// NonEmpty with a slice type
	var value optional.NonEmpty[[]string, string]
	value.MustParse([]string{"hello", "world"})
	val, exists := value.Get()
	fmt.Printf("Value: %v, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Empty slices will fail validation
	var invalidValue optional.NonEmpty[[]string, string]
	err := invalidValue.Parse([]string{})
	fmt.Printf("Empty slice error: %t\n", err != nil)

	// Output:
	// Value: [hello world], Exists: true
	// After nil: Exists: false
	// Empty slice error: true
}

func ExampleNonEmptySlice() {
	// NonEmptySlice is a simplified version of NonEmpty
	var value optional.NonEmptySlice[int]
	value.MustParse([]int{1, 2, 3})
	val, exists := value.Get()
	fmt.Printf("Value: %v, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Empty slices will fail validation
	var invalidValue optional.NonEmptySlice[int]
	err := invalidValue.Parse([]int{})
	fmt.Printf("Empty slice error: %t\n", err != nil)

	// Output:
	// Value: [1 2 3], Exists: true
	// After nil: Exists: false
	// Empty slice error: true
}

func ExampleMIME() {
	var value optional.MIME[string]

	// Valid MIME type
	value.MustParse("text/html")
	val, exists := value.Get()
	fmt.Printf("Value: %s, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Invalid MIME types will fail validation
	var invalidValue optional.MIME[string]
	err := invalidValue.Parse("not a mime type")
	fmt.Printf("Invalid MIME type error: %t\n", err != nil)

	// Output:
	// Value: text/html, Exists: true
	// After nil: Exists: false
	// Invalid MIME type error: true
}

func ExampleUUID() {
	var value optional.UUID[string]

	// Valid UUID
	value.MustParse("550e8400-e29b-41d4-a716-446655440000")
	_, exists := value.Get()
	fmt.Printf("Value exists: %t\n", exists)

	// Empty optional is valid
	value.MustParse(nil)
	_, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Invalid UUIDs will fail validation
	var invalidValue optional.UUID[string]
	err := invalidValue.Parse("not-a-uuid")
	fmt.Printf("Invalid UUID error: %t\n", err != nil)

	// Output:
	// Value exists: true
	// After nil: Exists: false
	// Invalid UUID error: true
}

func ExampleJSON() {
	var value optional.JSON[string]

	// Valid JSON
	value.MustParse(`{"name": "John", "age": 30}`)
	_, exists := value.Get()
	fmt.Printf("Has value: %t\n", exists)

	// Empty optional is valid
	value.MustParse(nil)
	_, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Invalid JSON will fail validation
	var invalidValue optional.JSON[string]
	err := invalidValue.Parse(`{"name": "John", "age": }`) // missing value
	fmt.Printf("Invalid JSON error: %t\n", err != nil)

	// Output:
	// Has value: true
	// After nil: Exists: false
	// Invalid JSON error: true
}

func ExampleCountryAlpha2() {
	var value optional.CountryAlpha2[string]

	// Valid country code
	value.MustParse("US")
	val, exists := value.Get()
	fmt.Printf("Value: %s, Exists: %t\n", val, exists)

	// Case-insensitive
	value.MustParse("gb")
	val, exists = value.Get()
	fmt.Printf("Case-insensitive: %s, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Invalid country codes will fail validation
	var invalidValue optional.CountryAlpha2[string]
	err := invalidValue.Parse("USA") // too long
	fmt.Printf("Three-letter code error: %t\n", err != nil)

	// Output:
	// Value: US, Exists: true
	// Case-insensitive: gb, Exists: true
	// After nil: Exists: false
	// Three-letter code error: true
}

func ExampleCountryAlpha3() {
	var value optional.CountryAlpha3[string]

	// Valid country code
	value.MustParse("USA")
	val, exists := value.Get()
	fmt.Printf("Value: %s, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Invalid country codes will fail validation
	var invalidValue optional.CountryAlpha3[string]
	err := invalidValue.Parse("US") // too short
	fmt.Printf("Two-letter code error: %t\n", err != nil)

	// Output:
	// Value: USA, Exists: true
	// After nil: Exists: false
	// Two-letter code error: true
}

func ExampleCountryAlpha() {
	var value optional.CountryAlpha[string]

	// Accepts both 2-letter codes
	value.MustParse("US")
	val, exists := value.Get()
	fmt.Printf("2-letter: %s, Exists: %t\n", val, exists)

	// And 3-letter codes
	value.MustParse("GBR")
	val, exists = value.Get()
	fmt.Printf("3-letter: %s, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Invalid country codes will fail validation
	var invalidValue optional.CountryAlpha[string]
	err := invalidValue.Parse("USAX") // too long
	fmt.Printf("Too long code error: %t\n", err != nil)

	// Output:
	// 2-letter: US, Exists: true
	// 3-letter: GBR, Exists: true
	// After nil: Exists: false
	// Too long code error: true
}

func ExampleCurrencyAlpha() {
	var value optional.CurrencyAlpha[string]

	// Valid currency code
	value.MustParse("USD")
	val, exists := value.Get()
	fmt.Printf("Value: %s, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Invalid currency codes will fail validation
	var invalidValue optional.CurrencyAlpha[string]
	err := invalidValue.Parse("US") // too short
	fmt.Printf("Too short code error: %t\n", err != nil)

	// Output:
	// Value: USD, Exists: true
	// After nil: Exists: false
	// Too short code error: true
}

func ExampleLangAlpha2() {
	var value optional.LangAlpha2[string]

	// Valid language code
	value.MustParse("en")
	val, exists := value.Get()
	fmt.Printf("Value: %s, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Invalid language codes will fail validation
	var invalidValue optional.LangAlpha2[string]
	err := invalidValue.Parse("eng") // too long
	fmt.Printf("Too long code error: %t\n", err != nil)

	// Output:
	// Value: en, Exists: true
	// After nil: Exists: false
	// Too long code error: true
}

func ExampleLangAlpha3() {
	var value optional.LangAlpha3[string]

	// Valid language code
	value.MustParse("eng")
	val, exists := value.Get()
	fmt.Printf("Value: %s, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Invalid language codes will fail validation
	var invalidValue optional.LangAlpha3[string]
	err := invalidValue.Parse("en") // too short
	fmt.Printf("Too short code error: %t\n", err != nil)

	// Output:
	// Value: eng, Exists: true
	// After nil: Exists: false
	// Too short code error: true
}

func ExampleLangAlpha() {
	var value optional.LangAlpha[string]

	// Accepts both 2-letter codes
	value.MustParse("en")
	val, exists := value.Get()
	fmt.Printf("2-letter: %s, Exists: %t\n", val, exists)

	// And 3-letter codes
	value.MustParse("ger")
	val, exists = value.Get()
	fmt.Printf("3-letter: %s, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Invalid language codes will fail validation
	var invalidValue optional.LangAlpha[string]
	err := invalidValue.Parse("engx") // too long
	fmt.Printf("Too long code error: %t\n", err != nil)

	// Output:
	// 2-letter: en, Exists: true
	// 3-letter: ger, Exists: true
	// After nil: Exists: false
	// Too long code error: true
}

func ExampleCustom() {
	// Custom allows creating an optional type with any validator
	var value optional.Custom[int, validate.Positive[int]]
	value.MustParse(42)
	val, exists := value.Get()
	fmt.Printf("Value: %d, Exists: %t\n", val, exists)

	// Empty optional is valid
	value.MustParse(nil)
	val, exists = value.Get()
	fmt.Printf("After nil: Exists: %t\n", exists)

	// Invalid values will fail validation
	var invalidValue optional.Custom[int, validate.Positive[int]]
	err := invalidValue.Parse(-5)
	fmt.Printf("Invalid value error: %t\n", err != nil)

	// Output:
	// Value: 42, Exists: true
	// After nil: Exists: false
	// Invalid value error: true
}

package required_test

import (
	"fmt"
	"time"

	"github.com/metafates/schema/required"
	"github.com/metafates/schema/validate"
	"github.com/metafates/schema/validate/charset"
)

func ExampleAny() {
	var value required.Any[string]
	value.MustParse("any value is acceptable")
	fmt.Println(value.Get())

	// Even an empty string is valid for Any
	value.MustParse("")
	fmt.Println("Empty string accepted:", value.Get() == "")

	// Output:
	// any value is acceptable
	// Empty string accepted: true
}

func ExampleNonZero() {
	var strValue required.NonZero[string]
	strValue.MustParse("non-empty string")
	fmt.Println(strValue.Get())

	var numValue required.NonZero[int]
	numValue.MustParse(42)
	fmt.Println(numValue.Get())

	// Empty strings and zero values will fail validation
	var invalidStr required.NonZero[string]
	err := invalidStr.Parse("")
	fmt.Println("Empty string error:", err != nil)

	var invalidNum required.NonZero[int]
	err = invalidNum.Parse(0)
	fmt.Println("Zero number error:", err != nil)

	// Output:
	// non-empty string
	// 42
	// Empty string error: true
	// Zero number error: true
}

func ExamplePositive() {
	var intValue required.Positive[int]
	intValue.MustParse(42)
	fmt.Println(intValue.Get())

	var floatValue required.Positive[float64]
	floatValue.MustParse(3.14)
	fmt.Println(floatValue.Get())

	// Zero and negative values will fail
	var invalidValue required.Positive[int]
	err := invalidValue.Parse(0)
	fmt.Println("Zero error:", err != nil)

	err = invalidValue.Parse(-5)
	fmt.Println("Negative error:", err != nil)

	// Output:
	// 42
	// 3.14
	// Zero error: true
	// Negative error: true
}

func ExampleNegative() {
	var intValue required.Negative[int]
	intValue.MustParse(-42)
	fmt.Println(intValue.Get())

	var floatValue required.Negative[float64]
	floatValue.MustParse(-3.14)
	fmt.Println(floatValue.Get())

	// Zero and positive values will fail
	var invalidValue required.Negative[int]
	err := invalidValue.Parse(0)
	fmt.Println("Zero error:", err != nil)

	err = invalidValue.Parse(5)
	fmt.Println("Positive error:", err != nil)

	// Output:
	// -42
	// -3.14
	// Zero error: true
	// Positive error: true
}

func ExamplePositive0() {
	var intValue required.Positive0[int]
	intValue.MustParse(42)
	fmt.Println(intValue.Get())

	// Zero is valid for Positive0
	intValue.MustParse(0)
	fmt.Println("Zero value:", intValue.Get())

	// Negative values will fail
	var invalidValue required.Positive0[int]
	err := invalidValue.Parse(-5)
	fmt.Println("Negative error:", err != nil)

	// Output:
	// 42
	// Zero value: 0
	// Negative error: true
}

func ExampleNegative0() {
	var intValue required.Negative0[int]
	intValue.MustParse(-42)
	fmt.Println(intValue.Get())

	// Zero is valid for Negative0
	intValue.MustParse(0)
	fmt.Println("Zero value:", intValue.Get())

	// Positive values will fail
	var invalidValue required.Negative0[int]
	err := invalidValue.Parse(5)
	fmt.Println("Positive error:", err != nil)

	// Output:
	// -42
	// Zero value: 0
	// Positive error: true
}

func ExampleEven() {
	var value required.Even[int]
	value.MustParse(42)
	fmt.Println(value.Get())

	value.MustParse(0)
	fmt.Println("Zero is even:", value.Get())

	value.MustParse(-8)
	fmt.Println("Negative even:", value.Get())

	// Odd numbers will fail
	var invalidValue required.Even[int]
	err := invalidValue.Parse(3)
	fmt.Println("Odd number error:", err != nil)

	// Output:
	// 42
	// Zero is even: 0
	// Negative even: -8
	// Odd number error: true
}

func ExampleOdd() {
	var value required.Odd[int]
	value.MustParse(43)
	fmt.Println(value.Get())

	value.MustParse(-9)
	fmt.Println("Negative odd:", value.Get())

	// Even numbers will fail
	var invalidValue required.Odd[int]
	err := invalidValue.Parse(2)
	fmt.Println("Even number error:", err != nil)

	err = invalidValue.Parse(0)
	fmt.Println("Zero error:", err != nil)

	// Output:
	// 43
	// Negative odd: -9
	// Even number error: true
	// Zero error: true
}

func ExampleEmail() {
	var value required.Email[string]
	value.MustParse("user@example.com")
	fmt.Println(value.Get())

	value.MustParse("John Doe <john.doe@example.com>")
	fmt.Println(value.Get())

	// Invalid emails will fail
	var invalidValue required.Email[string]
	err := invalidValue.Parse("not-an-email")
	fmt.Println("Invalid email error:", err != nil)

	err = invalidValue.Parse("@missing-user.com")
	fmt.Println("Missing username error:", err != nil)

	// Output:
	// user@example.com
	// John Doe <john.doe@example.com>
	// Invalid email error: true
	// Missing username error: true
}

func ExampleURL() {
	var value required.URL[string]
	value.MustParse("https://example.com")
	fmt.Println(value.Get())

	value.MustParse("http://localhost:8080/path?query=value")
	fmt.Println(value.Get())

	// Relative URLs are valid too
	value.MustParse("/relative/path")
	fmt.Println("Relative URL:", value.Get())

	// Invalid URLs will fail
	var invalidValue required.URL[string]
	err := invalidValue.Parse("htt ps://com example")
	fmt.Println("Invalid URL error:", err != nil)

	// Output:
	// https://example.com
	// http://localhost:8080/path?query=value
	// Relative URL: /relative/path
	// Invalid URL error: true
}

func ExampleHTTPURL() {
	var value required.HTTPURL[string]
	value.MustParse("https://example.com")
	fmt.Println(value.Get())

	value.MustParse("http://localhost:8080/path?query=value")
	fmt.Println(value.Get())

	// Non-HTTP URLs will fail
	var invalidValue required.HTTPURL[string]
	err := invalidValue.Parse("ftp://example.com")
	fmt.Println("Non-HTTP URL error:", err != nil)

	// Relative URLs will fail
	err = invalidValue.Parse("/relative/path")
	fmt.Println("Relative URL error:", err != nil)

	// Output:
	// https://example.com
	// http://localhost:8080/path?query=value
	// Non-HTTP URL error: true
	// Relative URL error: true
}

func ExampleIP() {
	var value required.IP[string]
	value.MustParse("192.168.1.1")
	fmt.Println(value.Get())

	value.MustParse("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	fmt.Println(value.Get())

	// Invalid IPs will fail
	var invalidValue required.IP[string]
	err := invalidValue.Parse("999.999.999.999")
	fmt.Println("Invalid IP error:", err != nil)

	err = invalidValue.Parse("not an ip")
	fmt.Println("Non-IP string error:", err != nil)

	// Output:
	// 192.168.1.1
	// 2001:0db8:85a3:0000:0000:8a2e:0370:7334
	// Invalid IP error: true
	// Non-IP string error: true
}

func ExampleIPV4() {
	var value required.IPV4[string]
	value.MustParse("192.168.1.1")
	fmt.Println(value.Get())

	value.MustParse("10.0.0.1")
	fmt.Println(value.Get())

	// IPv6 addresses will fail
	var invalidValue required.IPV4[string]
	err := invalidValue.Parse("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	fmt.Println("IPv6 address error:", err != nil)

	// Invalid IPs will fail
	err = invalidValue.Parse("999.999.999.999")
	fmt.Println("Invalid IP error:", err != nil)

	// Output:
	// 192.168.1.1
	// 10.0.0.1
	// IPv6 address error: true
	// Invalid IP error: true
}

func ExampleIPV6() {
	var value required.IPV6[string]
	value.MustParse("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	fmt.Println(value.Get())

	value.MustParse("::1") // localhost in IPv6
	fmt.Println(value.Get())

	// IPv4 addresses will fail
	var invalidValue required.IPV6[string]
	err := invalidValue.Parse("192.168.1.1")
	fmt.Println("IPv4 address error:", err != nil)

	// Invalid IPs will fail
	err = invalidValue.Parse("not an ip")
	fmt.Println("Invalid IP error:", err != nil)

	// Output:
	// 2001:0db8:85a3:0000:0000:8a2e:0370:7334
	// ::1
	// IPv4 address error: true
	// Invalid IP error: true
}

func ExampleMAC() {
	var value required.MAC[string]
	value.MustParse("00:1B:44:11:3A:B7")
	fmt.Println(value.Get())

	value.MustParse("00-1B-44-11-3A-B7") // hyphen format
	fmt.Println(value.Get())

	// Invalid MAC addresses will fail
	var invalidValue required.MAC[string]
	err := invalidValue.Parse("not a mac address")
	fmt.Println("Invalid MAC error:", err != nil)

	err = invalidValue.Parse("00:1B:44:11:3A") // too short
	fmt.Println("Too short MAC error:", err != nil)

	// Output:
	// 00:1B:44:11:3A:B7
	// 00-1B-44-11-3A-B7
	// Invalid MAC error: true
	// Too short MAC error: true
}

func ExampleCIDR() {
	var value required.CIDR[string]
	value.MustParse("192.168.1.0/24")
	fmt.Println(value.Get())

	value.MustParse("2001:db8::/32")
	fmt.Println(value.Get())

	// Invalid CIDRs will fail
	var invalidValue required.CIDR[string]
	err := invalidValue.Parse("192.168.1.1") // missing prefix
	fmt.Println("Missing prefix error:", err != nil)

	err = invalidValue.Parse("192.168.1.0/33") // invalid prefix length
	fmt.Println("Invalid prefix error:", err != nil)

	// Output:
	// 192.168.1.0/24
	// 2001:db8::/32
	// Missing prefix error: true
	// Invalid prefix error: true
}

func ExampleBase64() {
	var value required.Base64[string]
	value.MustParse("SGVsbG8gV29ybGQ=") // "Hello World"
	fmt.Println(value.Get())

	value.MustParse("dGVzdA==") // "test"
	fmt.Println(value.Get())

	// Invalid base64 will fail
	var invalidValue required.Base64[string]
	err := invalidValue.Parse("not!valid!base64!")
	fmt.Println("Invalid base64 error:", err != nil)

	// Output:
	// SGVsbG8gV29ybGQ=
	// dGVzdA==
	// Invalid base64 error: true
}

func ExampleCharset0() {
	// Charset0 allows empty strings
	var value required.Charset0[string, charset.Letter]
	value.MustParse("abcDEF")
	fmt.Println(value.Get())

	value.MustParse("")
	fmt.Println("Empty string accepted:", value.Get() == "")

	// Non-alphabetic characters will fail
	var invalidValue required.Charset0[string, charset.Letter]
	err := invalidValue.Parse("abc123")
	fmt.Println("Non-alphabetic error:", err != nil)

	// Output:
	// abcDEF
	// Empty string accepted: true
	// Non-alphabetic error: true
}

func ExampleCharset() {
	// Charset requires non-empty strings
	var value required.Charset[string, charset.Or[charset.Letter, charset.Number]]
	value.MustParse("abc123DEF")
	fmt.Println(value.Get())

	// Empty strings will fail
	var invalidValue required.Charset[string, charset.Or[charset.Letter, charset.Number]]
	err := invalidValue.Parse("")
	fmt.Println("Empty string error:", err != nil)

	// Non-alphanumeric characters will fail
	err = invalidValue.Parse("abc123!")
	fmt.Println("Non-alphanumeric error:", err != nil)

	// Output:
	// abc123DEF
	// Empty string error: true
	// Non-alphanumeric error: true
}

func ExampleLatitude() {
	var value required.Latitude[float64]
	value.MustParse(37.7749) // San Francisco
	fmt.Println(value.Get())

	value.MustParse(-33.8688) // Sydney (negative latitude)
	fmt.Println(value.Get())

	value.MustParse(90) // North Pole (max latitude)
	fmt.Println("North Pole:", value.Get())

	value.MustParse(-90) // South Pole (min latitude)
	fmt.Println("South Pole:", value.Get())

	// Out of range latitudes will fail
	var invalidValue required.Latitude[float64]
	err := invalidValue.Parse(91)
	fmt.Println("Too high error:", err != nil)

	err = invalidValue.Parse(-91)
	fmt.Println("Too low error:", err != nil)

	// Output:
	// 37.7749
	// -33.8688
	// North Pole: 90
	// South Pole: -90
	// Too high error: true
	// Too low error: true
}

func ExampleLongitude() {
	var value required.Longitude[float64]
	value.MustParse(-122.4194) // San Francisco
	fmt.Println(value.Get())

	value.MustParse(151.2093) // Sydney
	fmt.Println(value.Get())

	value.MustParse(180) // International Date Line (max longitude)
	fmt.Println("Date Line East:", value.Get())

	value.MustParse(-180) // International Date Line (min longitude)
	fmt.Println("Date Line West:", value.Get())

	// Out of range longitudes will fail
	var invalidValue required.Longitude[float64]
	err := invalidValue.Parse(181)
	fmt.Println("Too high error:", err != nil)

	err = invalidValue.Parse(-181)
	fmt.Println("Too low error:", err != nil)

	// Output:
	// -122.4194
	// 151.2093
	// Date Line East: 180
	// Date Line West: -180
	// Too high error: true
	// Too low error: true
}

func ExampleInPast() {
	var value required.InPast[time.Time]
	pastTime := time.Now().Add(-24 * time.Hour) // 1 day ago
	value.MustParse(pastTime)
	fmt.Println("Is in past:", value.Get().Before(time.Now()))

	// Future times will fail
	var invalidValue required.InPast[time.Time]
	futureTime := time.Now().Add(24 * time.Hour) // 1 day in future
	err := invalidValue.Parse(futureTime)
	fmt.Println("Future time error:", err != nil)

	// Output:
	// Is in past: true
	// Future time error: true
}

func ExampleInFuture() {
	var value required.InFuture[time.Time]
	futureTime := time.Now().Add(24 * time.Hour) // 1 day in future
	value.MustParse(futureTime)
	fmt.Println("Is in future:", value.Get().After(time.Now()))

	// Past times will fail
	var invalidValue required.InFuture[time.Time]
	pastTime := time.Now().Add(-24 * time.Hour) // 1 day ago
	err := invalidValue.Parse(pastTime)
	fmt.Println("Past time error:", err != nil)

	// Output:
	// Is in future: true
	// Past time error: true
}

func ExampleUnique() {
	// Unique with a slice type
	var value required.Unique[[]int, int]
	value.MustParse([]int{1, 2, 3, 4, 5})
	fmt.Println(value.Get())

	// Duplicate values will fail
	var invalidValue required.Unique[[]int, int]
	err := invalidValue.Parse([]int{1, 2, 3, 1, 4})
	fmt.Println("Duplicate value error:", err != nil)

	// Output:
	// [1 2 3 4 5]
	// Duplicate value error: true
}

func ExampleUniqueSlice() {
	// UniqueSlice is a simplified version of Unique
	var value required.UniqueSlice[string]
	value.MustParse([]string{"apple", "banana", "cherry"})
	fmt.Println(value.Get())

	// Duplicate values will fail
	var invalidValue required.UniqueSlice[string]
	err := invalidValue.Parse([]string{"apple", "banana", "apple"})
	fmt.Println("Duplicate value error:", err != nil)

	// Output:
	// [apple banana cherry]
	// Duplicate value error: true
}

func ExampleNonEmpty() {
	// NonEmpty with a slice type
	var value required.NonEmpty[[]string, string]
	value.MustParse([]string{"hello", "world"})
	fmt.Println(value.Get())

	// Empty slices will fail
	var invalidValue required.NonEmpty[[]string, string]
	err := invalidValue.Parse([]string{})
	fmt.Println("Empty slice error:", err != nil)

	// Output:
	// [hello world]
	// Empty slice error: true
}

func ExampleNonEmptySlice() {
	// NonEmptySlice is a simplified version of NonEmpty
	var value required.NonEmptySlice[int]
	value.MustParse([]int{1, 2, 3})
	fmt.Println(value.Get())

	// Empty slices will fail
	var invalidValue required.NonEmptySlice[int]
	err := invalidValue.Parse([]int{})
	fmt.Println("Empty slice error:", err != nil)

	// Output:
	// [1 2 3]
	// Empty slice error: true
}

func ExampleMIME() {
	var value required.MIME[string]
	value.MustParse("text/html")
	fmt.Println(value.Get())

	value.MustParse("application/json")
	fmt.Println(value.Get())

	value.MustParse("image/png")
	fmt.Println(value.Get())

	// Invalid MIME types will fail
	var invalidValue required.MIME[string]
	err := invalidValue.Parse("not a mime type")
	fmt.Println("Invalid MIME type error:", err != nil)

	// Output:
	// text/html
	// application/json
	// image/png
	// Invalid MIME type error: true
}

func ExampleUUID() {
	var value required.UUID[string]
	value.MustParse("550e8400-e29b-41d4-a716-446655440000")
	fmt.Println(value.Get())

	// Different UUID formats are accepted
	value.MustParse("urn:uuid:550e8400-e29b-41d4-a716-446655440000")
	fmt.Println("URN format accepted")

	value.MustParse("550e8400e29b41d4a716446655440000")
	fmt.Println("No hyphens format accepted")

	value.MustParse("{550e8400-e29b-41d4-a716-446655440000}")
	fmt.Println("Braces format accepted")

	// Invalid UUIDs will fail
	var invalidValue required.UUID[string]
	err := invalidValue.Parse("not-a-uuid")
	fmt.Println("Invalid UUID error:", err != nil)

	err = invalidValue.Parse("550e8400-e29b-41d4-a716-44665544000") // too short
	fmt.Println("Too short UUID error:", err != nil)

	// Output:
	// 550e8400-e29b-41d4-a716-446655440000
	// URN format accepted
	// No hyphens format accepted
	// Braces format accepted
	// Invalid UUID error: true
	// Too short UUID error: true
}

func ExampleJSON() {
	var value required.JSON[string]
	value.MustParse(`{"name": "John", "age": 30}`)
	fmt.Println(value.Get())

	value.MustParse(`[1, 2, 3, 4, 5]`)
	fmt.Println(value.Get())

	value.MustParse(`"simple string"`)
	fmt.Println(value.Get())

	// Invalid JSON will fail
	var invalidValue required.JSON[string]
	err := invalidValue.Parse(`{"name": "John", "age": }`) // missing value
	fmt.Println("Invalid JSON error:", err != nil)

	err = invalidValue.Parse(`not json`)
	fmt.Println("Not JSON error:", err != nil)

	// Output:
	// {"name": "John", "age": 30}
	// [1, 2, 3, 4, 5]
	// "simple string"
	// Invalid JSON error: true
	// Not JSON error: true
}

func ExampleCountryAlpha2() {
	var value required.CountryAlpha2[string]
	value.MustParse("US")
	fmt.Println(value.Get())

	// Case-insensitive
	value.MustParse("gb")
	fmt.Println(value.Get())

	// Invalid country codes will fail
	var invalidValue required.CountryAlpha2[string]
	err := invalidValue.Parse("USA") // too long
	fmt.Println("Three-letter code error:", err != nil)

	err = invalidValue.Parse("XX") // non-existent
	fmt.Println("Non-existent code error:", err != nil)

	// Output:
	// US
	// gb
	// Three-letter code error: true
	// Non-existent code error: true
}

func ExampleCountryAlpha3() {
	var value required.CountryAlpha3[string]
	value.MustParse("USA")
	fmt.Println(value.Get())

	// Case-insensitive
	value.MustParse("gbr")
	fmt.Println(value.Get())

	// Invalid country codes will fail
	var invalidValue required.CountryAlpha3[string]
	err := invalidValue.Parse("US") // too short
	fmt.Println("Two-letter code error:", err != nil)

	err = invalidValue.Parse("XXX") // non-existent
	fmt.Println("Non-existent code error:", err != nil)

	// Output:
	// USA
	// gbr
	// Two-letter code error: true
	// Non-existent code error: true
}

func ExampleCountryAlpha() {
	var value required.CountryAlpha[string]
	// Accepts both 2-letter codes
	value.MustParse("US")
	fmt.Println(value.Get())

	// And 3-letter codes
	value.MustParse("GBR")
	fmt.Println(value.Get())

	// Case-insensitive
	value.MustParse("jp")
	fmt.Println(value.Get())

	// Invalid country codes will fail
	var invalidValue required.CountryAlpha[string]
	err := invalidValue.Parse("USAX") // too long
	fmt.Println("Too long code error:", err != nil)

	err = invalidValue.Parse("XX") // non-existent
	fmt.Println("Non-existent code error:", err != nil)

	// Output:
	// US
	// GBR
	// jp
	// Too long code error: true
	// Non-existent code error: true
}

func ExampleCurrencyAlpha() {
	var value required.CurrencyAlpha[string]
	value.MustParse("USD")
	fmt.Println(value.Get())

	// Case-insensitive
	value.MustParse("eur")
	fmt.Println(value.Get())

	value.MustParse("GBP")
	fmt.Println(value.Get())

	// Invalid currency codes will fail
	var invalidValue required.CurrencyAlpha[string]
	err := invalidValue.Parse("US") // too short
	fmt.Println("Too short code error:", err != nil)

	err = invalidValue.Parse("ABC") // non-existent
	fmt.Println("Non-existent code error:", err != nil)

	// Output:
	// USD
	// eur
	// GBP
	// Too short code error: true
	// Non-existent code error: true
}

func ExampleLangAlpha2() {
	var value required.LangAlpha2[string]
	value.MustParse("en")
	fmt.Println(value.Get())

	// Case-insensitive
	value.MustParse("DE")
	fmt.Println(value.Get())

	value.MustParse("fr")
	fmt.Println(value.Get())

	// Invalid language codes will fail
	var invalidValue required.LangAlpha2[string]
	err := invalidValue.Parse("eng") // too long
	fmt.Println("Too long code error:", err != nil)

	err = invalidValue.Parse("xx") // non-existent
	fmt.Println("Non-existent code error:", err != nil)

	// Output:
	// en
	// DE
	// fr
	// Too long code error: true
	// Non-existent code error: true
}

func ExampleLangAlpha3() {
	var value required.LangAlpha3[string]
	value.MustParse("eng")
	fmt.Println(value.Get())

	// Case-insensitive
	value.MustParse("GER")
	fmt.Println(value.Get())

	value.MustParse("fre")
	fmt.Println(value.Get())

	// Invalid language codes will fail
	var invalidValue required.LangAlpha3[string]
	err := invalidValue.Parse("en") // too short
	fmt.Println("Too short code error:", err != nil)

	err = invalidValue.Parse("xxx") // non-existent
	fmt.Println("Non-existent code error:", err != nil)

	// Output:
	// eng
	// GER
	// fre
	// Too short code error: true
	// Non-existent code error: true
}

func ExampleLangAlpha() {
	var value required.LangAlpha[string]
	// Accepts both 2-letter codes
	value.MustParse("en")
	fmt.Println(value.Get())

	// And 3-letter codes
	value.MustParse("GER")
	fmt.Println(value.Get())

	// Case-insensitive
	value.MustParse("fRE")
	fmt.Println(value.Get())

	// Invalid language codes will fail
	var invalidValue required.LangAlpha[string]
	err := invalidValue.Parse("engx") // too long
	fmt.Println("Too long code error:", err != nil)

	err = invalidValue.Parse("xx") // non-existent
	fmt.Println("Non-existent code error:", err != nil)

	// Output:
	// en
	// GER
	// fRE
	// Too long code error: true
	// Non-existent code error: true
}

func ExampleCustom() {
	// Custom allows creating a required type with any validator
	var value required.Custom[int, validate.Positive[int]]
	value.MustParse(42)
	fmt.Println(value.Get())

	// Invalid values will fail according to the validator
	var invalidValue required.Custom[int, validate.Positive[int]]
	err := invalidValue.Parse(-5)
	fmt.Println("Invalid value error:", err != nil)

	// Output:
	// 42
	// Invalid value error: true
}

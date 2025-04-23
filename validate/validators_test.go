package validate_test

import (
	"reflect"
	"testing"
	"time"

	schemajson "github.com/metafates/schema/encoding/json"
	"github.com/metafates/schema/internal/testutil"
	"github.com/metafates/schema/required"
	. "github.com/metafates/schema/validate"
	"github.com/metafates/schema/validate/charset"
)

type TestCase[T any] struct {
	Name    string
	Input   T
	WantErr bool
}

type Suite[T any, V Validator[T]] []TestCase[T]

func (s Suite[T, V]) GetName() string {
	return reflect.TypeFor[V]().Name()
}

func (s Suite[T, V]) Test(t *testing.T) {
	var v V

	for _, tc := range s {
		t.Run(tc.Name, func(t *testing.T) {
			err := v.Validate(tc.Input)

			if tc.WantErr {
				testutil.Error(t, err)
			} else {
				testutil.NoError(t, err)
			}
		})
	}
}

func (s Suite[T, V]) Benchmark(b *testing.B) {
	var v V

	for _, tc := range s {
		if tc.WantErr {
			b.Run(tc.Name, func(b *testing.B) {
				for b.Loop() {
					testutil.Error(b, v.Validate(tc.Input))
				}
			})
		} else {
			b.Run(tc.Name, func(b *testing.B) {
				for b.Loop() {
					testutil.NoError(b, v.Validate(tc.Input))
				}
			})
		}
	}
}

type Testable interface {
	GetName() string

	Test(t *testing.T)
	Benchmark(b *testing.B)
}

var suites = []Testable{
	Suite[string, NonZero[string]]{
		{
			Name:  "non zero string",
			Input: "foo bar",
		},
		{
			Name:    "zero string",
			Input:   "",
			WantErr: true,
		},
		{
			Name:  "zero-width space",
			Input: "\u200B",
		},
	},
	Suite[int, Positive[int]]{
		{Name: "positive", Input: 42},
		{Name: "zero", Input: 0, WantErr: true},
		{Name: "negative", Input: -14, WantErr: true},
	},
	Suite[int, Negative[int]]{
		{Name: "positive", Input: 42, WantErr: true},
		{Name: "zero", Input: 0, WantErr: true},
		{Name: "negative", Input: -14},
	},
	Suite[string, Email[string]]{
		{
			Name:  "valid email",
			Input: `"john doe"@example.com (a comment)`,
		},
		{
			Name:    "empty string",
			Input:   "",
			WantErr: true,
		},
		{
			Name:    "invalid email",
			Input:   `a"b(c)d,e:f;g<h>i[j\k]l@example.com`,
			WantErr: true,
		},
	},
	Suite[string, URL[string]]{
		{
			Name:  "valid absolute url",
			Input: "https://example.com",
		},
		{
			Name:  "valid relative url",
			Input: "/example/com",
		},
		{
			Name:    "invalid url",
			Input:   "htt ps://com example",
			WantErr: true,
		},
	},
	Suite[string, HTTPURL[string]]{
		{
			Name:  "valid https url",
			Input: "https://example.com",
		},
		{
			Name:    "invalid url",
			Input:   "http s://e xample.com",
			WantErr: true,
		},
		{
			Name:    "relative url",
			Input:   "/example/com",
			WantErr: true,
		},
		{
			Name:    "non-http schema",
			Input:   "rpc://example.com",
			WantErr: true,
		},
	},
	Suite[string, IP[string]]{
		{
			Name:  "valid ipv4",
			Input: "127.0.0.1",
		},
		{
			Name:  "valid ipv6",
			Input: "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		},
		{
			Name:    "invalid ip",
			Input:   "127.0.0.A",
			WantErr: true,
		},
	},
	Suite[string, MAC[string]]{
		{
			Name:  "valid mac address",
			Input: "00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01",
		},
		{
			Name:    "invalid mac address",
			Input:   "lorem ipsum",
			WantErr: true,
		},
		{
			Name:    "empty string",
			Input:   "",
			WantErr: true,
		},
	},
	Suite[string, CIDR[string]]{
		{
			Name:  "valid cidr",
			Input: "192.0.2.0/24",
		},
		{
			Name:    "invalid cidr",
			Input:   "192.0.2.0@24",
			WantErr: true,
		},
		{
			Name:    "empty string",
			Input:   "",
			WantErr: true,
		},
	},
	Suite[string, Base64[string]]{
		{
			Name:  "valid base64 string",
			Input: "bG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQK",
		},
		{
			Name:    "invalid base64 string",
			Input:   "foo=bar",
			WantErr: true,
		},
	},
	Suite[string, Charset0[string, charset.ASCII]]{
		{
			Name:  "ascii only",
			Input: "The quick brown fox jumps over the lazy dog",
		},
		{
			Name:  "ascii with unprintable ascii",
			Input: "hello\x1B",
		},
		{
			Name:    "non-ascii",
			Input:   "Eĥoŝanĝoj ĉiuĵaŭde.", // esperanto btw
			WantErr: true,
		},
		{
			Name:  "empty",
			Input: "",
		},
	},
	Suite[string, Charset[string, charset.ASCII]]{
		{
			Name:  "ascii only",
			Input: "The quick brown fox jumps over the lazy dog",
		},
		{
			Name:  "ascii with unprintable ascii",
			Input: "hello\x1B",
		},
		{
			Name:    "non-ascii",
			Input:   "Eĥoŝanĝoj ĉiuĵaŭde.", // esperanto btw
			WantErr: true,
		},
		{
			Name:    "empty",
			Input:   "",
			WantErr: true,
		},
	},
	Suite[float64, Latitude[float64]]{
		{
			Name:  "valid negative latitude",
			Input: -24.20002,
		},
		{
			Name:  "valid positive latitude",
			Input: 90,
		},
		{
			Name:    "invalid positive latitude",
			Input:   90.0001,
			WantErr: true,
		},
		{
			Name:    "invalid negative latitude",
			Input:   -100,
			WantErr: true,
		},
	},
	Suite[float64, Longitude[float64]]{
		{
			Name:  "valid negative longitude",
			Input: -24.20002,
		},
		{
			Name:  "valid positive longitude",
			Input: 180,
		},
		{
			Name:    "invalid positive longitude",
			Input:   180.0001,
			WantErr: true,
		},
		{
			Name:    "invalid negative longitude",
			Input:   -200,
			WantErr: true,
		},
	},
	Suite[time.Time, InPast[time.Time]]{
		{
			Name:  "time is in the past",
			Input: time.Now().Add(-time.Hour),
		},
		{
			Name:    "time is in the future",
			Input:   time.Now().Add(time.Hour),
			WantErr: true,
		},
	},
	Suite[time.Time, InFuture[time.Time]]{
		{
			Name:    "time is in the past",
			Input:   time.Now().Add(-time.Hour),
			WantErr: true,
		},
		{
			Name:  "time is in the future",
			Input: time.Now().Add(time.Hour),
		},
	},
	Suite[[]string, Unique[[]string, string]]{
		{
			Name:  "unique strings",
			Input: []string{"foo", "bar", "foobar"},
		},
		{
			Name:    "duplicate strings",
			Input:   []string{"foo", "bar", "foo"},
			WantErr: true,
		},
		{
			Name:  "empty slice",
			Input: []string{},
		},
	},
	Suite[[]int, NonEmpty[[]int, int]]{
		{
			Name:  "non-empty slice",
			Input: []int{1, 2, 3},
		},
		{
			Name:    "empty slice",
			Input:   []int{},
			WantErr: true,
		},
		{
			Name:    "nil slice",
			Input:   nil,
			WantErr: true,
		},
	},
	Suite[string, MIME[string]]{
		{
			Name:  "simple valid mime type",
			Input: "application/json",
		},
		{
			Name:  "complex valid mime type",
			Input: `multipart/mixed; boundary="boundary-example"`,
		},
		{
			Name:    "invalid mime type",
			Input:   "blah blah",
			WantErr: true,
		},
		{
			Name:    "empty string",
			Input:   "",
			WantErr: true,
		},
	},
	Suite[string, UUID[string]]{
		{
			Name:  "standard valid uuid",
			Input: "550e8400-e29b-41d4-a716-446655440000",
		},
		{
			Name:  "urn valid uuid",
			Input: "urn:uuid:9b9773f5-ceb6-4e20-9bf6-7f83d6964877",
		},
		{
			Name:  "no-hyphens valid uuid",
			Input: "f47ac10b58cc4372a5670e02b2c3d479",
		},
		{
			Name:  "curly-braces valid uuid",
			Input: "{3d673a77-5f73-4608-a364-2a7c5c271d0c}",
		},
		{
			Name:    "empty string",
			Input:   "",
			WantErr: true,
		},
		{
			Name:    "invalid uuid",
			Input:   "XXXXXXXXXXXXXXXX hi",
			WantErr: true,
		},
	},
	Suite[int, Even[int]]{
		{Name: "even", Input: 2},
		{Name: "odd", Input: 3, WantErr: true},
	},
	Suite[int, Odd[int]]{
		{Name: "even", Input: 2, WantErr: true},
		{Name: "odd", Input: 3},
	},
	Suite[string, JSON[string]]{
		{
			Name:  "valid json",
			Input: `{"key":"value","array":[1, 2, 3]}`,
		},
		{
			Name:    "invalid json",
			Input:   `{"key" - "value" wait is going on??,"array":[1, 2, 3]}`,
			WantErr: true,
		},
	},
	Suite[string, CountryAlpha2[string]]{
		{
			Name:  "valid country code",
			Input: "us",
		},
		{
			Name:  "valid uppercase country code",
			Input: "US",
		},
		{
			Name:    "invalid country code",
			Input:   "..",
			WantErr: true,
		},
		{
			Name:    "empty country code",
			Input:   "",
			WantErr: true,
		},
		{
			Name:    "long country code",
			Input:   "usa",
			WantErr: true,
		},
	},
	Suite[string, CountryAlpha3[string]]{
		{
			Name:  "valid country code",
			Input: "usa",
		},
		{
			Name:  "valid uppercase country code",
			Input: "USA",
		},
		{
			Name:    "invalid country code",
			Input:   "...",
			WantErr: true,
		},
		{
			Name:    "empty country code",
			Input:   "",
			WantErr: true,
		},
		{
			Name:    "short country code",
			Input:   "us",
			WantErr: true,
		},
	},
	Suite[string, CurrencyAlpha[string]]{
		{
			Name:  "valid currency code",
			Input: "usd",
		},
		{
			Name:  "valid currency country code",
			Input: "USD",
		},
		{
			Name:    "invalid currency code",
			Input:   "...",
			WantErr: true,
		},
		{
			Name:    "empty currency code",
			Input:   "",
			WantErr: true,
		},
	},
	Suite[string, LangAlpha2[string]]{
		{
			Name:  "valid lang code",
			Input: "en",
		},
		{
			Name:  "valid uppercase lang code",
			Input: "EN",
		},
		{
			Name:    "invalid lang code",
			Input:   "..",
			WantErr: true,
		},
		{
			Name:    "empty lang code",
			Input:   "",
			WantErr: true,
		},
		{
			Name:    "long lang code",
			Input:   "eng",
			WantErr: true,
		},
	},
	Suite[string, LangAlpha3[string]]{
		{
			Name:  "valid lang code",
			Input: "eng",
		},
		{
			Name:  "valid uppercase lang code",
			Input: "ENG",
		},
		{
			Name:    "invalid lang code",
			Input:   "...",
			WantErr: true,
		},
		{
			Name:    "empty lang code",
			Input:   "",
			WantErr: true,
		},
		{
			Name:    "short lang code",
			Input:   "en",
			WantErr: true,
		},
	},
	Suite[int, And[int, NonZero[int], Positive[int]]]{
		{Name: "positive non zero", Input: 2},
		{Name: "zero", Input: 0, WantErr: true},
		{Name: "negative", Input: -2, WantErr: true},
	},
	Suite[int, Or[int, Even[int], Positive[int]]]{
		{Name: "positive even", Input: 2},
		{Name: "positive odd", Input: 3},
		{Name: "zero", Input: 0},
		{Name: "negative even", Input: -2},
		{Name: "negative odd", Input: -3, WantErr: true},
	},
	Suite[int, Not[int, Positive[int]]]{
		{Name: "positive", Input: 42, WantErr: true},
		{Name: "zero", Input: 0},
		{Name: "negative", Input: -14},
	},
}

func BenchmarkValidator(b *testing.B) {
	if testing.Short() {
		b.Skip("skipped benchmarking validators due to -short flag")
	}

	for _, s := range suites {
		b.Run(s.GetName(), s.Benchmark)
	}
}

func TestValidator(t *testing.T) {
	for _, s := range suites {
		t.Run(s.GetName(), s.Test)
	}
}

func TestValidate(t *testing.T) {
	type User struct {
		Name required.NonZero[string] `json:"name"`
		Age  int
	}

	t.Run("struct", func(t *testing.T) {
		var user User

		data := []byte(`{"name":"foo", "Age": 99}`)

		testutil.NoError(t, schemajson.Unmarshal(data, &user))

		testutil.Equal(t, "foo", user.Name.Get())
		testutil.Equal(t, 99, user.Age)
	})

	t.Run("slice", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			var users []User

			data := []byte(`[{"name": "foo"}, {"name": "bar", "Age": 99}]`)

			testutil.NoError(t, schemajson.Unmarshal(data, &users))

			for i, name := range []string{"foo", "bar"} {
				testutil.Equal(t, name, users[i].Name.Get())
			}
		})

		t.Run("error", func(t *testing.T) {
			var users []User

			data := []byte(`[{"name": "foo"}, {"bar": "other"}]`)

			testutil.Error(t, schemajson.Unmarshal(data, &users))
		})
	})
}

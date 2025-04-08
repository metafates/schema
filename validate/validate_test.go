package validate

import (
	"reflect"
	"testing"
	"time"
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

			switch {
			case !tc.WantErr && err != nil:
				t.Errorf("unexpected error: %v", err)
			case tc.WantErr && err == nil:
				t.Error("error did not occur")
			}
		})
	}
}

type Testable interface {
	GetName() string
	Test(t *testing.T)
}

func TestValidator(t *testing.T) {
	for _, tc := range []Testable{
		Suite[string, NonEmpty[string]]{
			{
				Name:  "non empty string",
				Input: "foo bar",
			},
			{
				Name:    "empty string",
				Input:   "",
				WantErr: true,
			},
			{
				Name:  "zero-width space",
				Input: "\u200B",
			},
		},
		Suite[int, Positive[int]]{
			{
				Name:  "positive",
				Input: 42,
			},
			{
				Name:  "zero",
				Input: 0,
			},
			{
				Name:    "negative",
				Input:   -14,
				WantErr: true,
			},
		},
		Suite[int, Negative[int]]{
			{
				Name:    "positive",
				Input:   42,
				WantErr: true,
			},
			{
				Name:  "zero",
				Input: 0,
			},
			{
				Name:  "negative",
				Input: -14,
			},
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
		Suite[string, Printable[string]]{
			{
				Name:  "printable string",
				Input: "lorem ipsum❤️",
			},
			{
				Name:    "mixed with unprintable",
				Input:   "lorem ipsum\x00",
				WantErr: true,
			},
			{
				Name:    "all unprintable",
				Input:   "\x1B\x08",
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
		Suite[string, ASCII[string]]{
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
		},
		Suite[string, PrintableASCII[string]]{
			{
				Name:  "ascii only",
				Input: "The quick brown fox jumps over the lazy dog",
			},
			{
				Name:    "ascii with unprintable ascii",
				Input:   "hello\x1B",
				WantErr: true,
			},
			{
				Name:    "non-ascii",
				Input:   "Eĥoŝanĝoj ĉiuĵaŭde.",
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
	} {
		t.Run(tc.GetName(), tc.Test)
	}
}

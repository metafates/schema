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

	"github.com/metafates/schema/internal/iso"
	"github.com/metafates/schema/internal/uuid"
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
		return errors.New("odd value")
	}

	return nil
}

func (Odd[T]) Validate(value T) error {
	if value%2 == 0 {
		return errors.New("even value")
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

	return Charset0[T, F]{}.Validate(value)
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

func (LangAlpha2[T]) Validate(value T) error {
	v := strings.ToLower(string(value))

	if _, ok := iso.LanguageAlpha2[v]; !ok {
		return errors.New("unknown 2-letter language code")
	}

	return nil
}

func (LangAlpha3[T]) Validate(value T) error {
	v := strings.ToLower(string(value))

	if _, ok := iso.LanguageAlpha3[v]; !ok {
		return errors.New("unknown 3-letter language code")
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

	//nolint:nilerr
	if err := v.Validate(value); err != nil {
		return nil
	}

	return errors.New(fmt.Sprint(v))
}

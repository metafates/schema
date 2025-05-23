package uuid

import (
	"errors"
	"fmt"
	"strings"
)

// xvalues returns the value of a byte as a hexadecimal digit or 255.
var xvalues = [256]byte{
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
}

// xtob states whether hex characters x1 and x2 could be converted into a byte.
func xtob(x1, x2 byte) bool {
	b1 := xvalues[x1]
	b2 := xvalues[x2]

	return b1 != 255 && b2 != 255
}

// Validate if given string is a valid UUID.
//
// https://github.com/google/uuid/blob/0f11ee6918f41a04c201eceeadf612a377bc7fbc/uuid.go#L195
//
//nolint:cyclop
func Validate(s string) error {
	const standardLen = 36

	switch len(s) {
	// Standard UUID format
	case standardLen:

	// UUID with "urn:uuid:" prefix
	case standardLen + 9:
		if !strings.EqualFold(s[:9], "urn:uuid:") {
			return fmt.Errorf("invalid urn prefix: %q", s[:9])
		}

		s = s[9:]

	// UUID enclosed in braces
	case standardLen + 2:
		if s[0] != '{' || s[len(s)-1] != '}' {
			return errors.New("invalid bracketed UUID format")
		}

		s = s[1 : len(s)-1]

	// UUID without hyphens
	case standardLen - 4:
		for i := 0; i < len(s); i += 2 {
			if !xtob(s[i], s[i+1]) {
				return errors.New("invalid UUID format")
			}
		}

	default:
		return fmt.Errorf("invalid UUID length: %d", len(s))
	}

	// Check for standard UUID format
	if len(s) == standardLen {
		if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
			return errors.New("invalid UUID format")
		}

		for _, x := range []int{0, 2, 4, 6, 9, 11, 14, 16, 19, 21, 24, 26, 28, 30, 32, 34} {
			if !xtob(s[x], s[x+1]) {
				return errors.New("invalid UUID format")
			}
		}
	}

	return nil
}

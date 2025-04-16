package charset

import (
	"testing"

	"github.com/metafates/schema/internal/testutil"
)

func TestFilter(t *testing.T) {
	for _, tc := range []struct {
		name           string
		filter         Filter
		valid, invalid []rune
	}{
		{
			name:    "ascii",
			filter:  ASCII{},
			valid:   []rune{'A'},
			invalid: []rune{'Ж'},
		},
		{
			name:    "graphic",
			filter:  Graphic{},
			valid:   []rune{'Ж'},
			invalid: []rune{0},
		},
		{
			name:    "print",
			filter:  Print{},
			valid:   []rune{'A'},
			invalid: []rune{0},
		},
		{
			name:    "control",
			filter:  Control{},
			valid:   []rune{0},
			invalid: []rune{'A'},
		},
		{
			name:    "letter",
			filter:  Letter{},
			valid:   []rune{'A'},
			invalid: []rune{'?'},
		},
		{
			name:    "mark",
			filter:  Mark{},
			valid:   []rune{0x300}, // Combining Grave Accent
			invalid: []rune{'?'},
		},
		{
			name:    "punct",
			filter:  Punct{},
			valid:   []rune{';'},
			invalid: []rune{'A'},
		},
		{
			name:    "space",
			filter:  Space{},
			valid:   []rune{' '},
			invalid: []rune{'_'},
		},
		{
			name:    "symbol",
			filter:  Symbol{},
			valid:   []rune{'✨'},
			invalid: []rune{0},
		},
		{
			name:    "and",
			filter:  And[Space, Print]{},
			valid:   []rune{' '},
			invalid: []rune{0xA0, '8'},
		},
		{
			name:    "or",
			filter:  Or[Letter, Number]{},
			valid:   []rune{'A', '1'},
			invalid: []rune{' ', '?'},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("valid", func(t *testing.T) {
				for _, r := range tc.valid {
					testutil.NoError(t, tc.filter.Filter(r))
				}
			})

			t.Run("invalid", func(t *testing.T) {
				for _, r := range tc.invalid {
					testutil.Error(t, tc.filter.Filter(r))
				}
			})
		})
	}
}

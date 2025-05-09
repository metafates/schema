// Package charset provides various charset filters to be used in combination with charset validator
package charset

import (
	"errors"
	"fmt"
	"unicode"
)

// Filter represents charset filter.
type Filter interface {
	Filter(r rune) error
}

type (
	// Any accepts any rune.
	Any struct{}

	// ASCII accepts ASCII runes.
	ASCII struct{}

	// Graphic wraps [unicode.IsGraphic].
	Graphic struct{}

	// Print wraps [unicode.IsPrint].
	Print struct{}

	// Control wraps [unicode.IsControl].
	Control struct{}

	// Letter wraps [unicode.IsLetter].
	Letter struct{}

	// Mark wraps [unicode.IsMark].
	Mark struct{}

	// Number wraps [unicode.IsNumber].
	Number struct{}

	// Punct wraps [unicode.IsPunct].
	Punct struct{}

	// Space wraps [unicode.IsSpace].
	Space struct{}

	// Symbol wraps [unicode.IsSymbol].
	Symbol struct{}

	// And is a meta filter that combines multiple filters using AND operator.
	And[A, B Filter] struct{}

	// Or is a meta filter that combines multiple filters using OR operator.
	Or[A, B Filter] struct{}

	// Not is a meta filter that inverts given filter.
	Not[F Filter] struct{}
)

// Common aliases.
type (
	// ASCIINumber intersects [ASCII] and [Number].
	ASCIINumber = And[ASCII, Number]

	// ASCIIPrint intersects [ASCII] and [Print].
	ASCIIPrint = And[ASCII, Print]

	// ASCIILetter intersects [ASCII] and [Letter].
	ASCIILetter = And[ASCII, Letter]

	// ASCIIPunct intersects [ASCII] and [Punct].
	ASCIIPunct = And[ASCII, Punct]
)

func (Any) Filter(rune) error       { return nil }
func (ASCII) Filter(r rune) error   { return assert(r <= unicode.MaxASCII, "non-ascii character") }
func (Graphic) Filter(r rune) error { return assert(unicode.IsGraphic(r), "non-graphic character") }
func (Print) Filter(r rune) error   { return assert(unicode.IsPrint(r), "non-printable character") }
func (Control) Filter(r rune) error { return assert(unicode.IsControl(r), "non-control character") }
func (Letter) Filter(r rune) error  { return assert(unicode.IsLetter(r), "non-letter character") }
func (Mark) Filter(r rune) error    { return assert(unicode.IsMark(r), "non-mark character") }
func (Number) Filter(r rune) error  { return assert(unicode.IsNumber(r), "non-number character") }

func (Punct) Filter(
	r rune,
) error {
	return assert(unicode.IsPunct(r), "non-punctuation character")
}
func (Space) Filter(r rune) error  { return assert(unicode.IsSpace(r), "non-space character") }
func (Symbol) Filter(r rune) error { return assert(unicode.IsSymbol(r), "non-symbol character") }

func (And[A, B]) Filter(r rune) error {
	if err := (*new(A)).Filter(r); err != nil {
		return err
	}

	if err := (*new(B)).Filter(r); err != nil {
		return err
	}

	return nil
}

func (Or[A, B]) Filter(r rune) error {
	errA := (*new(A)).Filter(r)
	if errA == nil {
		return nil
	}

	errB := (*new(B)).Filter(r)
	if errB == nil {
		return nil
	}

	return errors.Join(errA, errB)
}

func (Not[F]) Filter(r rune) error {
	var f F

	if err := f.Filter(r); err != nil {
		//nolint:nilerr
		return nil
	}

	return errors.New(fmt.Sprint(f))
}

func assert(condition bool, msg string) error {
	if !condition {
		return errors.New(msg)
	}

	return nil
}

package schemajson

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/metafates/schema/validate"
)

// Decoder wraps [json.Decoder] with validation step after decoding
//
// See [json.Decoder] documentation
type Decoder struct {
	*json.Decoder
}

// NewDecoder returns a new decoder that reads from r.
//
// See [json.NewDecoder] documentation
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{json.NewDecoder(r)}
}

// Decode wraps [json.Decoder.Decode] and calls [validate.Validate] afterwards.
//
// See also [Unmarshal]
func (dec *Decoder) Decode(v any) error {
	if err := dec.Decoder.Decode(v); err != nil {
		return err
	}

	if err := validate.Validate(v); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

// Unmarshal wraps [json.Unmarshal] and calls [validate.Validate] afterwards.
//
// See also [Decoder.Decode]
func Unmarshal(data []byte, v any) error {
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	if err := validate.Validate(v); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

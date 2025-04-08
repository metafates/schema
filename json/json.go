package schemajson

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/metafates/schema/validate"
)

type Decoder struct {
	*json.Decoder
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{json.NewDecoder(r)}
}

func (dec *Decoder) Decode(v any) error {
	if err := dec.Decoder.Decode(v); err != nil {
		return err
	}

	if err := validate.Recursively(v); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func Unmarshal(data []byte, v any) error {
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	if err := validate.Recursively(v); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

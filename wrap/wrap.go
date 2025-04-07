package wrap

import (
	"encoding/json"
	"fmt"

	"github.com/metafates/schema/validate"
)

// Wrapped is a type that adds support for validating the fields that implement [Validator] interface.
type Wrapped[T any] struct{ Inner T }

func (w *Wrapped[T]) UnmarshalJSON(data []byte) error {
	var inner T

	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}

	if err := validate.Recursively(inner); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	w.Inner = inner

	return nil
}

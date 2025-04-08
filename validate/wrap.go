package validate

import (
	"encoding/json"
	"fmt"
)

// Wrap is a type that validates its inner value as part of unmarshalling
type Wrap[T any] struct{ Inner T }

func (w *Wrap[T]) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &w.Inner); err != nil {
		return err
	}

	if err := Recursively(w.Inner); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

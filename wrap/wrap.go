package wrap

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/metafates/schema/internal/reflectwalk"
)

type Validater interface {
	Validate() error
}

// Wrapped is a type that adds support for validating the fields that implement [Validater] interface.
type Wrapped[T any] struct{ Inner T }

func (w *Wrapped[T]) UnmarshalJSON(data []byte) error {
	var inner T

	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}

	err := reflectwalk.WalkFields(inner, func(path string, value reflect.Value) error {
		r, ok := value.Interface().(Validater)
		if !ok {
			return nil
		}

		if err := r.Validate(); err != nil {
			return fmt.Errorf("%s: validate: %w", path, err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	w.Inner = inner

	return nil
}

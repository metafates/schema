package wrap

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/metafates/schema/internal/reflectwalk"
)

// Wrap is a type that adds support for checking the presence of the required fields, if any
type Wrap[T any] struct{ Inner T }

func (w *Wrap[T]) UnmarshalJSON(data []byte) error {
	var inner T

	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}

	err := reflectwalk.WalkFields(inner, func(path string, value reflect.Value) error {
		r, ok := value.Interface().(interface{ IsValid() bool })
		if !ok {
			return nil
		}

		if !r.IsValid() {
			return fmt.Errorf("%s: missing required value", path)
		}

		return nil
	})
	if err != nil {
		return err
	}

	w.Inner = inner

	return nil
}

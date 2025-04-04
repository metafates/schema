package requiredjson

import (
	"encoding/json"

	"github.com/metafates/required/internal/wrap"
)

func Unmarshal[T any](data []byte, v *T) error {
	var wrapped wrap.Wrap[T]

	if err := json.Unmarshal(data, &wrapped); err != nil {
		return err
	}

	*v = wrapped.Inner

	return nil
}

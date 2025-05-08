package required

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/metafates/schema/validate"
)

var _ interface {
	gob.GobDecoder
	gob.GobEncoder
} = (*Custom[any, validate.Validator[any]])(nil)

// GobDecode implements the [gob.GobDecoder] interface.
func (c *Custom[T, V]) GobDecode(data []byte) error {
	if len(data) == 0 {
		return errors.New("GobDecode: no data")
	}

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	var value T

	if err := dec.Decode(&value); err != nil {
		return err
	}

	*c = Custom[T, V]{value: value, hasValue: true}

	return nil
}

// GobEncode implements the [gob.GobEncoder] interface.
func (c Custom[T, V]) GobEncode() ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(c.Get()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

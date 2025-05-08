package optional

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

	if data[0] == 0 {
		*c = Custom[T, V]{}

		return nil
	}

	buf := bytes.NewBuffer(data[1:])
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
	if !c.hasValue {
		return []byte{0}, nil
	}

	var buf bytes.Buffer

	buf.WriteByte(1)

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(c.Must()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

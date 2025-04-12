package optional

import "encoding/json"

func (c *Custom[T, V]) UnmarshalJSON(data []byte) error {
	var value *T

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	// validated status will reset here
	if value == nil {
		*c = Custom[T, V]{}

		return nil
	}

	*c = Custom[T, V]{value: *value, hasValue: true}

	return nil
}

func (c Custom[T, V]) MarshalJSON() ([]byte, error) {
	if c.hasValue {
		return json.Marshal(c.value)
	}

	return []byte("null"), nil
}

package optional

func (c *Custom[T, V]) UnmarshalText(data []byte) error {
	return c.UnmarshalJSON(data)
}

func (c Custom[T, V]) MarshalText() ([]byte, error) {
	return c.MarshalJSON()
}

package optional

import (
	"encoding/json"
	"testing"
)

func TestOptional(t *testing.T) {
	t.Run("missing value", func(t *testing.T) {
		var foo Any[string]

		if err := json.Unmarshal([]byte(`null`), &foo); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		requireEqual(t, false, foo.hasValue)
		requireEqual(t, "", foo.value)
	})

	t.Run("invalid value", func(t *testing.T) {
		var foo Positive[int]

		if err := json.Unmarshal([]byte(`-24`), &foo); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		requireEqual(t, true, foo.hasValue)
		requireEqual(t, -24, foo.value) // won't be validated during unmarshalling
		{
			value, ok := foo.Value()
			requireEqual(t, foo.value, value)
			requireEqual(t, foo.hasValue, ok)
		}
	})

	t.Run("valid value", func(t *testing.T) {
		var foo Positive[int]

		if err := json.Unmarshal([]byte(`24`), &foo); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		requireEqual(t, true, foo.hasValue)
		requireEqual(t, 24, foo.value)
		{
			value, ok := foo.Value()
			requireEqual(t, foo.value, value)
			requireEqual(t, foo.hasValue, ok)
		}
	})
}

func requireEqual[T comparable](t *testing.T, want, actual T) {
	if want != actual {
		t.Fatalf("not equal: %v and %v", want, actual)
	}
}

package required

import (
	"encoding/json"
	"testing"
)

func TestRequired(t *testing.T) {
	t.Run("missing value", func(t *testing.T) {
		var foo Any[string]

		if err := json.Unmarshal([]byte(`null`), &foo); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		requireEqual(t, false, foo.hasValue)
		requireEqual(t, "", foo.value)

		requireError(t, foo.Validate())
	})

	t.Run("invalid value", func(t *testing.T) {
		var foo Positive[int]

		if err := json.Unmarshal([]byte(`-24`), &foo); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		requireEqual(t, true, foo.hasValue)
		requireEqual(t, -24, foo.value) // won't be validated during unmarshalling
		requireEqual(t, foo.value, foo.Value())

		requireError(t, foo.Validate())
	})

	t.Run("valid value", func(t *testing.T) {
		var foo Positive[int]

		if err := json.Unmarshal([]byte(`24`), &foo); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		requireEqual(t, true, foo.hasValue)
		requireEqual(t, 24, foo.value)
		requireEqual(t, foo.value, foo.Value())

		requireNoError(t, foo.Validate())
	})
}

func requireEqual[T comparable](t *testing.T, want, actual T) {
	if want != actual {
		t.Fatalf("not equal: %v and %v", want, actual)
	}
}

func requireNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireError(t *testing.T, err error) {
	if err == nil {
		t.Fatalf("error is nil")
	}
}

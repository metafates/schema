package validate_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/metafates/schema/required"
	"github.com/metafates/schema/validate"
)

func TestWrap(t *testing.T) {
	type User struct {
		Name required.NonEmpty[string] `json:"name"`
		Age  int
	}

	t.Run("struct", func(t *testing.T) {
		var wrapped validate.Wrap[User]

		data := []byte(`{"name":"foo", "Age": 99}`)

		if err := json.Unmarshal(data, &wrapped); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}

		inner := wrapped.Inner

		if inner.Name.Value() != "foo" {
			t.Errorf("name value: want %q, got %q", "foo", inner.Name.Value())
		}

		if inner.Age != 99 {
			t.Errorf("age: want 99 got %d", inner.Age)
		}
	})

	t.Run("slice", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			var wrapped validate.Wrap[[]User]

			data := []byte(`[{"name": "foo"}, {"name": "bar", "Age": 99}]`)

			if err := json.Unmarshal(data, &wrapped); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			inner := wrapped.Inner

			for i, name := range []string{"foo", "bar"} {
				if inner[i].Name.Value() != name {
					t.Errorf("[%d].name: want %q, got %q", i, name, inner[i].Name.Value())
				}
			}
		})

		t.Run("error", func(t *testing.T) {
			var wrapped validate.Wrap[[]User]

			data := []byte(`[{"name": "foo"}, {"bar": "other"}]`)

			if err := json.Unmarshal(data, &wrapped); err == nil {
				t.Fatal("unmarshal: no error")
			}
		})
	})
}

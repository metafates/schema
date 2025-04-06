package wrap

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/metafates/schema/required"
)

//go:embed testdata.json
var testdata []byte

type Data struct {
	ID         string   `json:"_id"`
	Index      int      `json:"index"`
	GUID       string   `json:"guid"`
	IsActive   bool     `json:"isActive"`
	Balance    string   `json:"balance"`
	Picture    string   `json:"picture"`
	Age        int      `json:"age"`
	EyeColor   string   `json:"eyeColor"`
	Name       string   `json:"name"`
	Gender     string   `json:"gender"`
	Company    string   `json:"company"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Address    string   `json:"address"`
	About      string   `json:"about"`
	Registered string   `json:"registered"`
	Latitude   float64  `json:"latitude"`
	Longitude  float64  `json:"longitude"`
	Tags       []string `json:"tags"`
	Friends    []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"friends"`
	Greeting      string `json:"greeting"`
	FavoriteFruit string `json:"favoriteFruit"`
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	b.Run("wrapped", func(b *testing.B) {
		for b.Loop() {
			var w Wrap[Data]

			if err := json.Unmarshal(testdata, &w); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("regular", func(b *testing.B) {
		for b.Loop() {
			var w Data

			if err := json.Unmarshal(testdata, &w); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestWrap(t *testing.T) {
	type User struct {
		Name required.NonEmpty[string] `json:"name"`
		Age  int
	}

	t.Run("struct", func(t *testing.T) {
		var wrapped Wrap[User]

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
			var wrapped Wrap[[]User]

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
			var wrapped Wrap[[]User]

			data := []byte(`[{"name": "foo"}, {"bar": "other"}]`)

			if err := json.Unmarshal(data, &wrapped); err == nil {
				t.Fatal("unmarshal: no error")
			}
		})
	})
}

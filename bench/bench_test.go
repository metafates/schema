package bench

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/metafates/schema/validate"
)

//go:embed testdata.json
var testdata []byte

type Data []struct {
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
	b.Run("unmarshal and validation with wrap", func(b *testing.B) {
		for b.Loop() {
			var w validate.OnUnmarshal[Data]

			if err := json.Unmarshal(testdata, &w); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("separate unmarshal and validation manually", func(b *testing.B) {
		for b.Loop() {
			var w Data

			if err := json.Unmarshal(testdata, &w); err != nil {
				b.Fatal(err)
			}

			if err := validate.Validate(&w); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("unmarshal without validation", func(b *testing.B) {
		for b.Loop() {
			var w Data

			if err := json.Unmarshal(testdata, &w); err != nil {
				b.Fatal(err)
			}
		}
	})
}

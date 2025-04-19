package bench

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/metafates/schema/required"
	"github.com/metafates/schema/validate"
)

//go:embed testdata.json
var testdata []byte

//go:generate schemagen -type DataWithGen

type DataWithGen []struct {
	ID         required.NonZero[string]    `json:"_id"`
	Index      int                         `json:"index"`
	GUID       string                      `json:"guid"`
	IsActive   bool                        `json:"isActive"`
	Balance    string                      `json:"balance"`
	Picture    string                      `json:"picture"`
	Age        required.Positive[int]      `json:"age"`
	EyeColor   string                      `json:"eyeColor"`
	Name       string                      `json:"name"`
	Gender     string                      `json:"gender"`
	Company    string                      `json:"company"`
	Email      string                      `json:"email"`
	Phone      string                      `json:"phone"`
	Address    string                      `json:"address"`
	About      string                      `json:"about"`
	Registered string                      `json:"registered"`
	Latitude   required.Latitude[float64]  `json:"latitude"`
	Longitude  required.Longitude[float64] `json:"longitude"`
	Tags       []string                    `json:"tags"`
	Friends    []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"friends"`
	Greeting      string `json:"greeting"`
	FavoriteFruit string `json:"favoriteFruit"`
}

type Data []struct {
	ID         required.NonZero[string]    `json:"_id"`
	Index      int                         `json:"index"`
	GUID       string                      `json:"guid"`
	IsActive   bool                        `json:"isActive"`
	Balance    string                      `json:"balance"`
	Picture    string                      `json:"picture"`
	Age        required.Positive[int]      `json:"age"`
	EyeColor   string                      `json:"eyeColor"`
	Name       string                      `json:"name"`
	Gender     string                      `json:"gender"`
	Company    string                      `json:"company"`
	Email      string                      `json:"email"`
	Phone      string                      `json:"phone"`
	Address    string                      `json:"address"`
	About      string                      `json:"about"`
	Registered string                      `json:"registered"`
	Latitude   required.Latitude[float64]  `json:"latitude"`
	Longitude  required.Longitude[float64] `json:"longitude"`
	Tags       []string                    `json:"tags"`
	Friends    []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"friends"`
	Greeting      string `json:"greeting"`
	FavoriteFruit string `json:"favoriteFruit"`
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	b.Run("reflection", func(b *testing.B) {
		b.Run("with validation", func(b *testing.B) {
			for b.Loop() {
				var data Data

				if err := json.Unmarshal(testdata, &data); err != nil {
					b.Fatal(err)
				}

				if err := validate.Validate(&data); err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run("without validation", func(b *testing.B) {
			for b.Loop() {
				var data Data

				if err := json.Unmarshal(testdata, &data); err != nil {
					b.Fatal(err)
				}
			}
		})
	})

	b.Run("codegen", func(b *testing.B) {
		b.Run("with validation", func(b *testing.B) {
			for b.Loop() {
				var data DataWithGen

				if err := json.Unmarshal(testdata, &data); err != nil {
					b.Fatal(err)
				}

				if err := validate.Validate(&data); err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run("without validation", func(b *testing.B) {
			for b.Loop() {
				var data DataWithGen

				if err := json.Unmarshal(testdata, &data); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

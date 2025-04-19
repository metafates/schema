package bench

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/metafates/schema/internal/bench"
	"github.com/metafates/schema/validate"
)

//go:generate schemagen -type A
type A struct{}

//go:embed testdata.json
var testdata []byte

func BenchmarkUnmarshalJSON(b *testing.B) {
	b.Run("reflection", func(b *testing.B) {
		b.Run("with validation", func(b *testing.B) {
			for b.Loop() {
				var w bench.Data

				if err := json.Unmarshal(testdata, &w); err != nil {
					b.Fatal(err)
				}

				if err := validate.Validate(&w); err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run("without validation", func(b *testing.B) {
			for b.Loop() {
				var w bench.Data

				if err := json.Unmarshal(testdata, &w); err != nil {
					b.Fatal(err)
				}
			}
		})
	})

	b.Run("codegen", func(b *testing.B) {
		b.Run("with validation", func(b *testing.B) {
			for b.Loop() {
				var w bench.DataWithGen

				if err := json.Unmarshal(testdata, &w); err != nil {
					b.Fatal(err)
				}

				if err := validate.Validate(&w); err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run("without validation", func(b *testing.B) {
			for b.Loop() {
				var w bench.DataWithGen

				if err := json.Unmarshal(testdata, &w); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

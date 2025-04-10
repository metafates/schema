package bench

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/metafates/schema/internal/bench"
	"github.com/metafates/schema/validate"
)

//go:embed testdata.json
var testdata []byte

func BenchmarkUnmarshalJSON(b *testing.B) {
	b.Run("without codegen", func(b *testing.B) {
		b.Run("unmarshal and validation with wrap", func(b *testing.B) {
			for b.Loop() {
				var w validate.OnUnmarshal[bench.Data]

				if err := json.Unmarshal(testdata, &w); err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run("separate unmarshal and validation manually", func(b *testing.B) {
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

		b.Run("unmarshal without validation", func(b *testing.B) {
			for b.Loop() {
				var w bench.Data

				if err := json.Unmarshal(testdata, &w); err != nil {
					b.Fatal(err)
				}
			}
		})
	})

	b.Run("with codegen", func(b *testing.B) {
		b.Run("unmarshal and validation with wrap", func(b *testing.B) {
			for b.Loop() {
				var w validate.OnUnmarshal[bench.DataWithGen]

				if err := json.Unmarshal(testdata, &w); err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run("separate unmarshal and validation manually", func(b *testing.B) {
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

		b.Run("unmarshal without validation", func(b *testing.B) {
			for b.Loop() {
				var w bench.DataWithGen

				if err := json.Unmarshal(testdata, &w); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

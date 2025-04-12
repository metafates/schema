package optional

import (
	"encoding/json"
	"testing"

	"github.com/metafates/schema/internal/testutil"
)

func TestOptional(t *testing.T) {
	t.Run("missing value", func(t *testing.T) {
		var foo Any[string]

		testutil.NoError(t, json.Unmarshal([]byte(`null`), &foo))

		testutil.Equal(t, false, foo.hasValue)
		testutil.Equal(t, foo.hasValue, foo.HasValue())
		testutil.Equal(t, "", foo.value)

		testutil.NoError(t, foo.Validate())
		testutil.Equal(t, true, foo.validated)

		testutil.Panic(t, func() { foo.Must() })
		testutil.NoPanic(t, func() { foo.Get() })
		testutil.NoPanic(t, func() { foo.MarshalJSON() })
		testutil.NoPanic(t, func() { foo.MarshalText() })
		testutil.NoPanic(t, func() { foo.Value() })
	})

	t.Run("invalid value", func(t *testing.T) {
		var foo Positive[int]

		testutil.NoError(t, json.Unmarshal([]byte(`-24`), &foo))

		testutil.Equal(t, true, foo.hasValue)
		testutil.Equal(t, foo.hasValue, foo.HasValue())
		testutil.Equal(t, -24, foo.value)

		testutil.Error(t, foo.Validate())
		testutil.Equal(t, false, foo.validated)

		testutil.Panic(t, func() { foo.Must() })
		testutil.Panic(t, func() { foo.Get() })
		testutil.Panic(t, func() { foo.MarshalJSON() })
		testutil.Panic(t, func() { foo.MarshalText() })
		testutil.Panic(t, func() { foo.Value() })
	})

	t.Run("nested invalid value", func(t *testing.T) {
		type Foo struct {
			Field Positive[int]
		}

		var foo Any[Foo]

		testutil.NoError(t, json.Unmarshal([]byte(`{"field":-1}`), &foo))
		testutil.Error(t, foo.Validate())
	})

	t.Run("valid value", func(t *testing.T) {
		var foo Positive[int]

		testutil.NoError(t, json.Unmarshal([]byte(`24`), &foo))

		testutil.Equal(t, true, foo.hasValue)
		testutil.Equal(t, foo.hasValue, foo.HasValue())
		testutil.Equal(t, 24, foo.value)

		testutil.NoError(t, foo.Validate())
		testutil.Equal(t, true, foo.validated)

		testutil.NoPanic(t, func() { foo.Must() })
		testutil.NoPanic(t, func() { foo.Get() })
		testutil.NoPanic(t, func() { foo.MarshalJSON() })
		testutil.NoPanic(t, func() { foo.MarshalText() })
		testutil.NoPanic(t, func() { foo.Value() })

		t.Run("reuse as invalid", func(t *testing.T) {
			testutil.NoError(t, json.Unmarshal([]byte(`24`), &foo))

			testutil.NoError(t, json.Unmarshal([]byte(`-24`), &foo))

			testutil.Equal(t, true, foo.hasValue)
			testutil.Equal(t, -24, foo.value)

			testutil.Error(t, foo.Validate())
			testutil.Equal(t, false, foo.validated)

			testutil.Panic(t, func() { foo.Must() })
			testutil.Panic(t, func() { foo.Get() })
			testutil.Panic(t, func() { foo.MarshalJSON() })
			testutil.Panic(t, func() { foo.MarshalText() })
			testutil.Panic(t, func() { foo.Value() })
		})
	})
}
